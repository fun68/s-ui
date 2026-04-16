package migration

import (
	"fmt"
	"github.com/alireza0/s-ui/config"
	"github.com/alireza0/s-ui/database"
	"github.com/alireza0/s-ui/database/model"
	"log"
)

func MigrateDb() {
	// void running on first install
	err := database.OpenDB()
	if err != nil {
		log.Fatal("Open db error: ", err)
		return
	}
	db := database.GetDB()

	tx := db.Begin()
	defer func() {
		if err == nil {
			tx.Commit()
		} else {
			tx.Rollback()
		}
	}()
	currentVersion := config.GetVersion()
	dbVersion := ""
	err = tx.Model(&model.Setting{}).
		Select("value").
		Where("`key` = ?", "version").
		Scan(&dbVersion).Error

	if err != nil {
		log.Printf("Failed to get DB version: %v", err)
	}
	fmt.Println("Current version:", currentVersion, "\nDatabase version:", dbVersion)

	if currentVersion == dbVersion {
		fmt.Println("Database is up to date, no need to migrate")
		return
	}
	fmt.Println("Start migrating database...")

	// Set version
	err = tx.Model(&model.Setting{}).
		Where("`key` = ?", "version").
		Update("value", currentVersion).Error
	if err != nil {
		log.Fatal("Update version failed: ", err)
		return
	}
	fmt.Println("Migration done!")
}
