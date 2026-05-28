package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/reshap0318/hadirYuk/cmd/migration/seeders"
	"github.com/reshap0318/hadirYuk/internal/database"
	"github.com/reshap0318/hadirYuk/internal/helpers"
	"github.com/reshap0318/hadirYuk/internal/models"
	"gorm.io/gorm"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Get database configuration from environment
	dbConnection := helpers.GetEnv("DB_CONNECTION", "mysql")
	dbHost := helpers.GetEnv("DB_HOST", "127.0.0.1")
	dbPort := helpers.GetEnv("DB_PORT", "3306")
	dbUser := helpers.GetEnv("DB_USERNAME", "root")
	dbPassword := helpers.GetEnv("DB_PASSWORD", "")
	dbName := helpers.GetEnv("DB_DATABASE", "boilerplate")

	if dbHost == "" || dbPort == "" || dbUser == "" || dbName == "" {
		log.Fatal("Database configuration is incomplete. Please set DB_HOST, DB_PORT, DB_USERNAME, and DB_DATABASE")
	}

	// Initialize database connection
	db := initDB(dbConnection, dbHost, dbPort, dbUser, dbPassword, dbName)

	// Parse command
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "up":
		runMigration(db, "up")
	case "down":
		runMigration(db, "down")
	case "seed":
		runSeed(db)
	case "refresh":
		runMigration(db, "down")
		runMigration(db, "up")
		runSeed(db)
	default:
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Usage: go run cmd/migration/main.go <command>")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  up       Apply all migrations (create/update tables)")
	fmt.Println("  down     Drop all migrated tables")
	fmt.Println("  seed     Insert default data")
	fmt.Println("  refresh  Drop all tables, recreate, and seed data")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  go run cmd/migration/main.go up")
	fmt.Println("  go run cmd/migration/main.go down")
	fmt.Println("  go run cmd/migration/main.go seed")
	fmt.Println("  go run cmd/migration/main.go refresh")
}

func runMigration(db *gorm.DB, command string) {
	switch command {
	case "up":
		fmt.Println("Running migrations...")

		// AutoMigrate all models
		err := db.AutoMigrate(
			&models.User{},
			&models.UserProfile{},
			&models.PasswordReset{},
			&models.Permission{},
			&models.Role{},
			&models.RoleHasPermission{},
			&models.UserHasRole{},
			&models.Notification{},
			&models.Shift{},
			&models.OfficeLocation{},
			&models.LeaveType{},
			&models.UserShiftAssignment{},
		)
		if err != nil {
			log.Fatalf("Migration failed: %v", err)
		}

		fmt.Println("✓ Migration completed successfully!")

	case "down":
		fmt.Println("Rolling back migrations...")

		// Drop tables in correct order (foreign key constraints)
		err := db.Migrator().DropTable(
			&models.UserShiftAssignment{},
			&models.LeaveType{},
			&models.OfficeLocation{},
			&models.Shift{},
			&models.UserProfile{},
			&models.Notification{},
			&models.UserHasRole{},
			&models.RoleHasPermission{},
			&models.Role{},
			&models.Permission{},
			&models.PasswordReset{},
			&models.User{},
		)
		if err != nil {
			log.Fatalf("Rollback failed: %v", err)
		}

		fmt.Println("✓ Rollback completed successfully!")
	}
}

func runSeed(db *gorm.DB) {
	fmt.Println("\n🌱 Seeding default data...\n")

	// Seed in correct order
	permIDs := seeders.SeedPermissions(db)
	roleIDs := seeders.SeedRoles(db)
	userEmails := seeders.SeedUsers(db)
	seeders.SeedRolePermissions(db, roleIDs, permIDs)
	seeders.SeedUserRoles(db, userEmails, roleIDs)

	fmt.Println("\n✅ Seeding completed!")
}

func initDB(dbConnection, dbHost, dbPort, dbUser, dbPassword, dbName string) *gorm.DB {
	var db *gorm.DB
	var err error

	if dbConnection == "postgres" || dbConnection == "postgresql" {
		db, err = database.NewPostgreSQL(database.PostgreSQLConfig{
			Host:     dbHost,
			Port:     dbPort,
			User:     dbUser,
			Password: dbPassword,
			DBName:   dbName,
		})
	} else {
		db, err = database.NewMySQL(database.MySQLConfig{
			Host:     dbHost,
			Port:     dbPort,
			User:     dbUser,
			Password: dbPassword,
			DBName:   dbName,
		})
	}

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	return db
}
