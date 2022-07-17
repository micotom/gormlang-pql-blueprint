package db

import (
	"errors"
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

	models := []interface{}{&models.Player{}, &models.ValueEntry{}, &models.Team{}}
	db.AutoMigrate(models...)

	return db
}

func GetOrCreatePlayer(db *gorm.DB, p models.Player) (models.Player, error) {
	err := db.FirstOrCreate(&p).Preload("ValueEntries", func(db *gorm.DB) *gorm.DB {
		return db.Order("value_entries.Day ASC")
	}).Error
	return p, err
}

func SavePlayer(db *gorm.DB, p models.Player) error {
	return db.Save(p).Error
}

func GetPlayerBySlug(db *gorm.DB, slug string) (models.Player, error) {
	var p models.Player
	err := db.Where("slug = ?", slug).Preload("ValueEntries", func(db *gorm.DB) *gorm.DB {
		return db.Order("value_entries.Day ASC")
	}).First(&p).Error
	return p, err
}

func GetAllPlayers(db *gorm.DB) ([]models.Player, error) {
	var players []models.Player
	err := db.Model(&models.Player{}).Preload("ValueEntries", func(db *gorm.DB) *gorm.DB {
		return db.Order("value_entries.Day ASC")
	}).Find(&players).Error
	return players, err
}

func AddPlayerToTeam(db *gorm.DB, t *models.Team, slug string) error {
	if p, e := GetPlayerBySlug(db, slug); e == nil {
		t.Players = append(t.Players, p)
		db.Save(t)
		return nil
	} else {
		return e
	}
}

func DeletePlayerFromTeam(db *gorm.DB, t *models.Team, slug string) error {
	var pos = -1
	for i, p2 := range t.Players {
		if p2.Slug == slug {
			pos = i
			break
		}
	}
	if pos != -1 {
		t.Players = append(t.Players[:pos], t.Players[pos+1:]...)
		return db.Save(t).Error
	}
	return errors.New("No such player")
}

func GetTeam(db *gorm.DB, slug string) (models.Team, error) {
	var t models.Team
	e := db.Model(&models.Team{}).Where("slug = ?", slug).Preload("Players").Preload("Players.ValueEntries", func(db *gorm.DB) *gorm.DB {
		return db.Order("value_entries.Day ASC")
	}).First(&t).Error
	return t, e
}

func CreateTeam(db *gorm.DB, slug string) (models.Team, error) {
	var t models.Team = models.Team{
		Slug: slug,
	}
	e := db.FirstOrCreate(&t).Error
	return t, e
}
