package main

import (
	"os"

	"funglejunk.com/kick-api/src/db"
	"funglejunk.com/kick-api/src/handlers"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
}

func main() {
	DB := db.Init()
	h := handlers.New(DB)

	r := gin.Default()
	r.GET("/scrape", h.DoScrape)
	r.GET("/players", h.GetAllPlayers)
	r.GET("/players/:slug", h.GetAllEntriesForPlayer)
	r.GET("/players/:slug/current", h.GetCurrentEntryForPlayer)
	r.Run() // listen and serve on 0.0.0.0:8080
}
