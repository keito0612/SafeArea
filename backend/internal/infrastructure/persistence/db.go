package persistence

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB

// InitDB はデータベース接続を初期化する
func InitDB() (*gorm.DB, error) {
	if db != nil {
		return db, nil
	}

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		getEnv("DB_HOST", "localhost"),
		getEnv("DB_PORT", "5432"),
		getEnv("DB_USER", "postgres"),
		getEnv("DB_PASSWORD", "postgres"),
		getEnv("DB_NAME", "safearea"),
		getEnv("DB_SSLMODE", "disable"),
	)

	// GORMロガー設定
	gormLogger := logger.New(
		slog.NewLogLogger(slog.Default().Handler(), slog.LevelInfo),
		logger.Config{
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  getLogLevel(),
			IgnoreRecordNotFoundError: true,
			Colorful:                  false,
		},
	)

	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}

	// コネクションプール設定
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB: %w", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(20)
	sqlDB.SetConnMaxLifetime(time.Hour)

	slog.Info("Database connection established")
	return db, nil
}

// GetDB はデータベース接続を取得する
func GetDB() *gorm.DB {
	return db
}

// CloseDB はデータベース接続を閉じる
func CloseDB() error {
	if db == nil {
		return nil
	}

	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB: %w", err)
	}

	if err := sqlDB.Close(); err != nil {
		return fmt.Errorf("failed to close database: %w", err)
	}

	db = nil
	slog.Info("Database connection closed")
	return nil
}

// HealthCheck はデータベースの接続状態を確認する
func HealthCheck() error {
	if db == nil {
		return fmt.Errorf("database connection is not initialized")
	}

	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB: %w", err)
	}

	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	return nil
}

// getEnv は環境変数を取得し、存在しない場合はデフォルト値を返す
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getLogLevel は環境に応じたログレベルを返す
func getLogLevel() logger.LogLevel {
	env := os.Getenv("ENV")
	switch env {
	case "prod", "production":
		return logger.Warn
	case "dev", "development":
		return logger.Info
	default:
		return logger.Info
	}
}
