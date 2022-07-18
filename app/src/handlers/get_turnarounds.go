package handlers

import (
	"net/http"
	"time"

	"funglejunk.com/kick-api/src/db"
	"funglejunk.com/kick-api/src/models"
	"github.com/gin-gonic/gin"
)

func (h handler) GetPlayersWithTurnoarounds(c *gin.Context) {

	ps, e := db.GetAllPlayers(h.DB)
	if _, e := CheckResult(c, ps, e); e == nil {
		turnarounds := []models.Player{}
		today := time.Now()
		yesterday := today.AddDate(0, 0, -1)
		tY, tM, tD := today.Date()
		yY, yM, yD := yesterday.Date()
		for _, p := range ps {
			if tE, e := models.ValueEntryAt(p.ValueEntries, tY, tM, tD); e == nil {
				if yE, e := models.ValueEntryAt(p.ValueEntries, yY, yM, yD); e == nil {
					if yE.RaiseDiff < 0 && tE.RaiseDiff > 0 {
						turnarounds = append(turnarounds, p)
					}
				}
			}
		}
		c.HTML(http.StatusOK, "turnarounds.html", turnarounds)
		//c.JSON(http.StatusOK, turnarounds)
	}

}
