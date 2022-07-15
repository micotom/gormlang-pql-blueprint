package db

import (
	"fmt"
	"os"

	"funglejunk.com/kick-api/src/models"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Init() *gorm.DB {
	host := os.Getenv("DB_HOST")
	user := os.Getenv("POSTGRES_USER")
	dbName := os.Getenv("POSTGRES_DB")
	pass := os.Getenv("POSTGRES_PASSWORD")
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=5432 sslmode=disable",
		host, user, pass, dbName,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Error(err)
	}

	models := []interface{}{&models.Player{}, &models.ValueEntry{}}
	db.AutoMigrate(models...)

	return db
}

func GetAllPlayers(db *gorm.DB) ([]models.Player, error) {
	var ps []models.Player
	err := db.Find(&ps).Error
	return ps, err
}

func GetPlayerBySlug(db *gorm.DB, slug string) (models.Player, error) {
	var p models.Player
	err := db.Where("slug = ?", slug).First(&p).Error
	return p, err
}

func GetAllPlayersWithEntries(db *gorm.DB) ([]models.Player, error) {
	var players []models.Player
	err := db.Model(&models.Player{}).Preload("ValueEntries").Find(&players).Error
	return players, err
}

func GetEntriesForPlayer(db *gorm.DB, slug string) ([]models.ValueEntry, error) {
	var currentEntries []models.ValueEntry
	err := db.Where("player_slug = ?", slug).Find(&currentEntries).Error
	return currentEntries, err
}

func GetCurrentEntry(db *gorm.DB, slug string) (models.ValueEntry, error) {
	var entry models.ValueEntry
	err := db.Where("player_slug = ?", slug).Order("Day desc").First(&entry).Error
	return entry, err
}
