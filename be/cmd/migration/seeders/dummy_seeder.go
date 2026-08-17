package seeders

import (
	_ "embed"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

//go:embed dummy_data.sql
var dummyDataSQL string

// SeedDummyData inserts portfolio/demo data: sample employees, shift
// assignments and 30 days of attendance history. Idempotent, safe to re-run.
func SeedDummyData(db *gorm.DB) error {
	fmt.Println("Seeding dummy portfolio data...")

	statements := strings.Split(dummyDataSQL, ";")

	executed := 0
	for _, stmt := range statements {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" {
			continue
		}
		if err := db.Exec(stmt).Error; err != nil {
			return fmt.Errorf("dummy data statement failed: %w\n--- SQL ---\n%s", err, stmt)
		}
		executed++
	}

	fmt.Printf("✓ Seeded dummy data (%d statements)\n", executed)
	return nil
}
