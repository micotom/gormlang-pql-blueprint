package main

import (
	"os"

	"funglejunk.com/kick-api/src/db"
	"funglejunk.com/kick-api/src/handlers"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
}

func main() {
	log.WithFields(log.Fields{
		"animal": "walrus",
	}).Info("A walrus appears")

	DB := db.Init()
	h := handlers.New(DB)

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "puh",
		})
	})
	r.GET("/players", h.GetAllPlayers)
	r.Run() // listen and serve on 0.0.0.0:8080

}
