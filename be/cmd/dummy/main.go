// dummy generates a randomized attendance record for TODAY for the
// portfolio demo employees (see cmd/migration/seeders/dummy_data.sql — the
// email list below must stay in sync with it). Meant to be run once a day by
// a scheduler (see storage/logs/dummy.log for output), so the demo
// dashboard always has fresh "today" data without anyone touching the DB by
// hand.
//
// DEMO/DEV ONLY — never wire this into the production entrypoint/crontab
// (etc/crontab), it fabricates attendance rows.
package main

import (
	"flag"
	"log"
	"math"
	"math/rand"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/reshap0318/hadirYuk/internal/database"
	"github.com/reshap0318/hadirYuk/internal/helpers"
	"github.com/reshap0318/hadirYuk/internal/models"
	"gorm.io/gorm"
)

// Keep this list in sync with cmd/migration/seeders/dummy_data.sql section 4.
var targetEmails = []string{
	"budi.santoso@app.com",
	"siti.nurhaliza@app.com",
	"andi.wijaya@app.com",
	"dewi.lestari@app.com",
	"rizky.pratama@app.com",
	"putri.ayu@app.com",
	"maya.kusuma@app.com",
}

const (
	presentPct = 70 // 0-69
	latePct    = 90 // 70-89, remainder (90-99) is absent
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	dryRun := flag.Bool("dry-run", false, "Print what would be created without writing to DB")
	flag.Parse()

	tz := helpers.GetEnv("TZ", "Asia/Jakarta")
	loc, err := time.LoadLocation(tz)
	if err != nil {
		log.Fatalf("Invalid timezone %q: %v", tz, err)
	}

	logger, err := helpers.NewLoggerWithSuffix("storage/logs", "dummy")
	if err != nil {
		log.Fatalf("[dummy] failed to initialize logger: %v", err)
	}
	defer logger.Close()

	now := time.Now().In(loc)
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)

	if today.Weekday() == time.Saturday || today.Weekday() == time.Sunday {
		logger.Printf("[dummy] %s is a weekend, nothing to do", today.Format("2006-01-02"))
		return
	}

	db := initDB(logger)

	var shift models.Shift
	if err := db.Where("name = ?", "Pagi").First(&shift).Error; err != nil {
		logger.Fatalf("[dummy] shift 'Pagi' not found — run `go run cmd/migration/main.go dummy` first: %v", err)
	}

	var office models.OfficeLocation
	if err := db.Where("name = ?", "Kantor Pusat").First(&office).Error; err != nil {
		logger.Fatalf("[dummy] office 'Kantor Pusat' not found — run `go run cmd/migration/main.go dummy` first: %v", err)
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	created, skipped := 0, 0
	for _, email := range targetEmails {
		var user models.User
		if err := db.Where("email = ?", email).First(&user).Error; err != nil {
			logger.Printf("[dummy] user %s not found, skipping — run the dummy seeder first", email)
			skipped++
			continue
		}

		var assignment models.UserShiftAssignment
		err := db.Where("user_id = ? AND shift_id = ? AND start_date <= ? AND (end_date IS NULL OR end_date >= ?) AND is_active = true",
			user.ID, shift.ID, today, today).First(&assignment).Error
		if err != nil {
			logger.Printf("[dummy] %s has no active shift assignment for %s, skipping", email, today.Format("2006-01-02"))
			skipped++
			continue
		}

		var count int64
		db.Model(&models.Attendance{}).
			Where("user_id = ? AND shift_id = ? AND date = ?", user.ID, shift.ID, today).
			Count(&count)
		if count > 0 {
			logger.Printf("[dummy] %s already has an attendance record for %s, skipping", email, today.Format("2006-01-02"))
			skipped++
			continue
		}

		record, roll := buildRecord(rng, user.ID, shift.ID, office, today, loc)
		if record == nil {
			logger.Printf("  %-28s -> absent (roll=%d), no record created", email, roll)
			continue
		}

		if *dryRun {
			logger.Printf("  [dry-run] %-28s -> %-7s check-in %s (roll=%d)", email, record.Status, record.TimeIn.Format("15:04"), roll)
			created++
			continue
		}

		if err := db.Create(record).Error; err != nil {
			logger.Printf("[dummy] ERROR creating record for %s: %v", email, err)
			continue
		}
		logger.Printf("  %-28s -> %-7s check-in %s (roll=%d)", email, record.Status, record.TimeIn.Format("15:04"), roll)
		created++
	}

	logger.Printf("[dummy] done — created: %d, skipped: %d", created, skipped)
	if *dryRun {
		os.Exit(0)
	}
}

// buildRecord randomly rolls today's attendance outcome for one employee.
// Returns (nil, roll) for the absent case — caller does not insert a row,
// mirroring how a real absence just means no check-in ever happened.
func buildRecord(rng *rand.Rand, userID, shiftID uint, office models.OfficeLocation, today time.Time, loc *time.Location) (*models.Attendance, int) {
	roll := rng.Intn(100)

	shiftStart := time.Date(today.Year(), today.Month(), today.Day(), 8, 0, 0, 0, loc)

	var status string
	var timeIn time.Time
	switch {
	case roll < presentPct:
		status = "present"
		// -10..+15 minutes around shift start, within the 15min flexi window.
		timeIn = shiftStart.Add(time.Duration(rng.Intn(26)-10) * time.Minute)
	case roll < latePct:
		status = "late"
		// 16..55 minutes late, past the flexi window.
		timeIn = shiftStart.Add(time.Duration(16+rng.Intn(40)) * time.Minute)
	default:
		return nil, roll
	}

	latIn, lngIn, distance := jitterLocation(rng, office.Latitude, office.Longitude)

	return &models.Attendance{
		UserID:         userID,
		ShiftID:        shiftID,
		OfficeID:       office.ID,
		Date:           today,
		Status:         status,
		TimeIn:         &timeIn,
		LatIn:          &latIn,
		LngIn:          &lngIn,
		DistanceMeters: &distance,
	}, roll
}

// jitterLocation returns a point randomly offset up to ~80m from the office
// (comfortably inside a typical geofence radius) plus its haversine distance.
func jitterLocation(rng *rand.Rand, officeLat, officeLng float64) (lat, lng, distanceMeters float64) {
	const maxOffsetMeters = 80.0
	const metersPerDegLat = 111320.0

	offsetMeters := rng.Float64() * maxOffsetMeters
	bearing := rng.Float64() * 2 * math.Pi

	dLat := (offsetMeters * math.Cos(bearing)) / metersPerDegLat
	metersPerDegLng := metersPerDegLat * math.Cos(officeLat*math.Pi/180)
	dLng := (offsetMeters * math.Sin(bearing)) / metersPerDegLng

	lat = officeLat + dLat
	lng = officeLng + dLng
	distanceMeters = haversineDistance(officeLat, officeLng, lat, lng)
	return
}

func haversineDistance(lat1, lng1, lat2, lng2 float64) float64 {
	const earthRadius = 6371000 // meters

	dLat := (lat2 - lat1) * math.Pi / 180
	dLng := (lng2 - lng1) * math.Pi / 180

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1*math.Pi/180)*math.Cos(lat2*math.Pi/180)*
			math.Sin(dLng/2)*math.Sin(dLng/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return earthRadius * c
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
		logger.Fatalf("[dummy] failed to connect to database: %v", err)
	}
	return db
}
