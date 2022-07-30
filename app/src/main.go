package main

import (
	"fmt"
	"html/template"
	"os"
	"time"

	"funglejunk.com/kick-api/src/db"
	"funglejunk.com/kick-api/src/handlers"
	"funglejunk.com/kick-api/src/models"
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"gorm.io/datatypes"
)

func init() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
}

func main() {
	DB := db.Init()

	h := handlers.New(DB)

	c := cron.New()
	c.AddFunc("0 8 * * *", func() {
		h.DoScrape()
	})

	r := gin.Default()
	r.SetFuncMap(template.FuncMap{
		"dateStr":                    dateStr,
		"moneyStr":                   moneyStr,
		"trend":                      models.PlayerValueTrend,
		"teamValue":                  models.TeamCurrentValue,
		"teamRaise":                  models.TeamTotalRaise,
		"playerCurrentValue":         models.PlayerGetCurrentValue,
		"playerCurrentRaiseDiff":     models.PlayerGetCurrentRaiseDiff,
		"playerCurrentRaiseDiffPerc": models.PlayerGetCurrentRaisePerc,
	})

	r.Static("/assets", "./assets")
	r.LoadHTMLGlob("./tmpl/*.html")

	r.GET("/overview", h.GetOverview)

	r.GET("/players", h.GetAllPlayers)
	r.GET("/players/:slug", h.GetPlayer)
	r.GET("/players/compare", h.ComparePlayers)
	r.GET("/players/turnarounds", h.GetPlayersWithTurnoarounds)
	r.GET("/players/top", h.GetTopPlayersX)
	r.GET("/players/positions", h.GetPlayersByPosition)

	r.GET("/teams/:slug", h.GetTeam)
	r.POST("/teams", h.CreateTeam)
	r.GET("/teams/add/:team_slug", h.AddPlayerToTeam)
	r.GET("/teams/delete/:team_slug", h.RemovePlayerFromTeam)

	r.Run() // listen and serve on 0.0.0.0:8080

}

func dateStr(date datatypes.Date) string {
	y, m, d := time.Time(date).Date()
	return fmt.Sprintf("%d-%d-%d", d, m, y)
}

func moneyStr(i int) string {
	p := message.NewPrinter(language.German)
	return p.Sprintf("%d", i)
}
