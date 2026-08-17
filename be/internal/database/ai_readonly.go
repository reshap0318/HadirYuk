package database

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// NewMySQLReadOnly opens a second GORM connection using the AI chat's
// dedicated readonly MySQL user (see plans/current/AI_CS_Chatbot.md §6.1 for
// the GRANT script). Kept separate from the main connection so a bug in
// AI-generated SQL can never touch a write-capable credential.
func NewMySQLReadOnly(cfg MySQLConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local&tls=false",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Error),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to AI readonly MySQL: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get AI readonly sql.DB: %w", err)
	}
	sqlDB.SetMaxOpenConns(5)
	sqlDB.SetMaxIdleConns(2)

	return db, nil
}
