package database

import (
	"fmt"
	"os"
	"scheduler/internal/app/models"

	"github.com/jackc/pgx/v5/pgconn"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Database struct {
	Client *gorm.DB
}

// NewDatabase - returns a pointer to a database object
func NewDatabase() (*Database, error) {
	log.Info("Setting up new database connection")
	dbURL := os.Getenv("DB_URL")
	log.Info("DB URL", dbURL)

	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{
		Logger: nil,
	})
	if err != nil {
		return &Database{}, fmt.Errorf("failed to connect to database: %w", err)
	}
	db.Logger.LogMode(logger.Warn)
	err = ping(db)
	if err != nil {
		return &Database{}, err
	}

	return &Database{
		Client: db,
	}, nil
}

func ping(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("Failed to get underlying DB instance:", err)
	}
	return sqlDB.Ping()
}

func DBMigration(store *Database) error {
	err := store.Client.AutoMigrate(models.Account{}, models.Schedule{})
	if err != nil {
		log.Error("failed to setup database")
		return err
	}
	uniqueConstraintSQL := `
		DO $$ 
		BEGIN
			IF NOT EXISTS (
				SELECT 1
				FROM pg_constraint
				WHERE conname = 'uc_account_work_date'
			) THEN
				EXECUTE 'ALTER TABLE schedules ADD CONSTRAINT uc_account_work_date UNIQUE (account_id, work_date)';
			END IF;
		END $$;
	`
	err = store.Client.Exec(uniqueConstraintSQL).Error
	if err != nil {
		log.Error("failed to add CONSTRAINT to schedules table")
		return err
	}
	return nil
}

func IsUniqueConstraintViolation(err error) bool {
	pgErr, ok := err.(*pgconn.PgError)
	if !ok {
		return false
	}
	return pgErr.Code == "23505" // PostgreSQL error code for unique violation
}
