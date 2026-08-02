package main

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/reshap0318/hadirYuk/internal/database"
	"github.com/reshap0318/hadirYuk/internal/helpers"
	"github.com/reshap0318/hadirYuk/internal/models"
	"gorm.io/gorm"
)

// ponytail: lookbackDays re-evaluates a few days back on every run so a shift
// that was still in progress on a prior run (skipped, see shiftEndDateTime)
// gets caught automatically — no need to time the cron schedule around shift
// end times. Bump this if a shift's flexi/duration ever exceeds ~2 days.
const lookbackDays = 3

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	dateStr := flag.String("date", "", "Target date YYYY-MM-DD (default: lookback window ending yesterday)")
	dryRun := flag.Bool("dry-run", false, "Print records that would be created without writing to DB")
	flag.Parse()

	tz := helpers.GetEnv("TZ", "Asia/Jakarta")
	loc, err := time.LoadLocation(tz)
	if err != nil {
		log.Fatalf("Invalid timezone %q: %v", tz, err)
	}

	logger, err := helpers.NewLoggerWithSuffix("storage/logs", "markabsent")
	if err != nil {
		log.Fatalf("[markabsent] failed to initialize logger: %v", err)
	}
	defer logger.Close()

	now := time.Now().In(loc)
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)

	var targetDates []time.Time
	if *dateStr != "" {
		d, err := time.ParseInLocation("2006-01-02", *dateStr, loc)
		if err != nil {
			logger.Fatalf("Invalid date %q, expected YYYY-MM-DD", *dateStr)
		}
		targetDates = []time.Time{d}
	} else {
		// Process a rolling window ending yesterday, oldest first — a shift
		// skipped on a previous run (still in progress) gets re-evaluated here.
		for i := lookbackDays; i >= 1; i-- {
			targetDates = append(targetDates, today.AddDate(0, 0, -i))
		}
	}

	db := initDB(logger)

	// Find the default office (needed for the absent record's required office_id).
	var defaultOffice models.OfficeLocation
	if err := db.Where("is_active = true").First(&defaultOffice).Error; err != nil {
		logger.Fatalf("[markabsent] no active office location found — cannot create absent records: %v", err)
	}

	totalCreated, totalSkipped, totalFailed := 0, 0, 0
	for _, targetDate := range targetDates {
		created, skipped, failed := processDate(db, logger, loc, now, targetDate, defaultOffice.ID, *dryRun)
		totalCreated += created
		totalSkipped += skipped
		totalFailed += failed
	}

	logger.Printf("[markabsent] done — created: %d, skipped: %d, failed: %d", totalCreated, totalSkipped, totalFailed)
	if totalFailed > 0 {
		os.Exit(1)
	}
}

// processDate marks absent every active shift assignment on targetDate that has
// no attendance record and whose shift has already ended (cross-midnight aware).
// Assignments whose shift is still running are left alone — a later run
// (today's or a subsequent day's, via lookbackDays) will re-check them.
func processDate(db *gorm.DB, logger *helpers.Logger, loc *time.Location, now, targetDate time.Time, officeID uint, dryRun bool) (created, skipped, failed int) {
	dateLabel := targetDate.Format("2006-01-02")
	startOfDay := targetDate
	endOfDay := startOfDay.Add(24 * time.Hour)

	var assignments []models.UserShiftAssignment
	if err := db.
		Preload("Shift").
		Where("is_active = true AND start_date <= ? AND (end_date IS NULL OR end_date >= ?)",
			startOfDay, startOfDay).
		Find(&assignments).Error; err != nil {
		logger.Printf("[markabsent] ERROR fetching assignments for %s: %v", dateLabel, err)
		failed++
		return
	}

	logger.Printf("[markabsent] %s: %d active shift assignment(s)", dateLabel, len(assignments))

	for _, assignment := range assignments {
		// Already has an attendance record (present/late/absent) — nothing to do.
		var count int64
		db.Model(&models.Attendance{}).
			Where("user_id = ? AND shift_id = ? AND date >= ? AND date < ?",
				assignment.UserID, assignment.ShiftID, startOfDay, endOfDay).
			Count(&count)
		if count > 0 {
			skipped++
			continue
		}

		// Shift still running (e.g. cross-midnight "Malam" 22:00-06:00) — its
		// check-in window hasn't closed yet, so it's too early to call it absent.
		if shiftEnd, ok := shiftEndDateTime(assignment.Shift, startOfDay, loc); ok && now.Before(shiftEnd) {
			skipped++
			continue
		}

		if dryRun {
			logger.Printf("  [dry-run] absent  user_id=%-4d  shift_id=%-3d  date=%s",
				assignment.UserID, assignment.ShiftID, dateLabel)
			created++
			continue
		}

		record := &models.Attendance{
			UserID:   assignment.UserID,
			ShiftID:  assignment.ShiftID,
			OfficeID: officeID,
			Date:     startOfDay,
			Status:   "absent",
		}

		if err := db.Create(record).Error; err != nil {
			logger.Printf("[markabsent] ERROR  user_id=%d  shift_id=%d  date=%s  err=%v",
				assignment.UserID, assignment.ShiftID, dateLabel, err)
			failed++
			continue
		}

		logger.Printf("  marked absent  user_id=%-4d  shift_id=%-3d  date=%s",
			assignment.UserID, assignment.ShiftID, dateLabel)
		created++
	}

	return
}

// shiftEndDateTime resolves a shift's actual end datetime on the given date,
// rolling over to the next day for cross-midnight shifts (end < start).
func shiftEndDateTime(shift models.Shift, date time.Time, loc *time.Location) (time.Time, bool) {
	start, errStart := time.Parse("15:04", shift.StartTime)
	end, errEnd := time.Parse("15:04", shift.EndTime)
	if errStart != nil || errEnd != nil {
		return time.Time{}, false
	}

	endDay := date
	if end.Hour()*60+end.Minute() < start.Hour()*60+start.Minute() {
		endDay = date.AddDate(0, 0, 1)
	}

	return time.Date(endDay.Year(), endDay.Month(), endDay.Day(), end.Hour(), end.Minute(), 0, 0, loc), true
}

func initDB(logger *helpers.Logger) *gorm.DB {
	cfg := database.MySQLConfig{
		Host:     helpers.GetEnv("DB_HOST", "127.0.0.1"),
		Port:     helpers.GetEnv("DB_PORT", "3306"),
		User:     helpers.GetEnv("DB_USERNAME", "root"),
		Password: helpers.GetEnv("DB_PASSWORD", ""),
		DBName:   helpers.GetEnv("DB_DATABASE", "hadir_yuk"),
	}
	db, err := database.NewMySQL(cfg)
	if err != nil {
		logger.Fatalf("[markabsent] failed to connect to database: %v", err)
	}
	return db
}
