package database

import (
	"log"

	"github.com/jaydeadlondon/project_na_go/internal/models"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	err := db.AutoMigrate(
		&models.User{},
		&models.Monitor{},
		&models.MonitorCheck{},
		&models.Incident{},
	)
	if err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	log.Println("✅ Migrations completed")
}