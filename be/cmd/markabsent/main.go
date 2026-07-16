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

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	dateStr := flag.String("date", "", "Target date YYYY-MM-DD (default: yesterday in local timezone)")
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

	var targetDate time.Time
	if *dateStr != "" {
		targetDate, err = time.ParseInLocation("2006-01-02", *dateStr, loc)
		if err != nil {
			logger.Fatalf("Invalid date %q, expected YYYY-MM-DD", *dateStr)
		}
	} else {
		now := time.Now().In(loc)
		targetDate = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc).AddDate(0, 0, -1)
	}

	dateLabel := targetDate.Format("2006-01-02")
	logger.Printf("[markabsent] target date: %s (dry-run: %v)", dateLabel, *dryRun)

	db := initDB(logger)

	// Use start/end of day range — consistent with how the attendance repo queries dates.
	startOfDay := time.Date(targetDate.Year(), targetDate.Month(), targetDate.Day(), 0, 0, 0, 0, loc)
	endOfDay := startOfDay.Add(24 * time.Hour)

	// Find the default office (needed for the absent record's required office_id).
	var defaultOffice models.OfficeLocation
	if err := db.Where("is_active = true").First(&defaultOffice).Error; err != nil {
		logger.Fatalf("[markabsent] no active office location found — cannot create absent records: %v", err)
	}

	// All active shift assignments that cover the target date.
	var assignments []models.UserShiftAssignment
	if err := db.
		Where("is_active = true AND start_date <= ? AND (end_date IS NULL OR end_date >= ?)",
			startOfDay, startOfDay).
		Find(&assignments).Error; err != nil {
		logger.Fatalf("[markabsent] failed to fetch shift assignments: %v", err)
	}

	logger.Printf("[markabsent] %d active shift assignment(s) on %s", len(assignments), dateLabel)

	created, skipped, failed := 0, 0, 0

	for _, assignment := range assignments {
		// Check if an attendance record already exists for this user+shift+date.
		var count int64
		db.Model(&models.Attendance{}).
			Where("user_id = ? AND shift_id = ? AND date >= ? AND date < ?",
				assignment.UserID, assignment.ShiftID, startOfDay, endOfDay).
			Count(&count)

		if count > 0 {
			skipped++
			continue
		}

		if *dryRun {
			logger.Printf("  [dry-run] absent  user_id=%-4d  shift_id=%-3d  date=%s",
				assignment.UserID, assignment.ShiftID, dateLabel)
			created++
			continue
		}

		record := &models.Attendance{
			UserID:   assignment.UserID,
			ShiftID:  assignment.ShiftID,
			OfficeID: defaultOffice.ID,
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

	logger.Printf("[markabsent] done — created: %d, skipped: %d, failed: %d", created, skipped, failed)
	if failed > 0 {
		os.Exit(1)
	}
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
