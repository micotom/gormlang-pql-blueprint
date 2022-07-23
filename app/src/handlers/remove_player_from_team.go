package handlers

import (
	"net/http"

	"funglejunk.com/kick-api/src/db"
	"github.com/gin-gonic/gin"
)

func (h handler) RemovePlayerFromTeam(c *gin.Context) {
	teamSlug := c.Param("team_slug")
	if playerSlug, ok := c.GetQuery("p"); ok {
		e := db.DeletePlayerFromTeam(h.DB, teamSlug, playerSlug)
		if e != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
		} else {
			c.Redirect(301, "/teams/t1")
			// c.Status(http.StatusOK)
		}
	}
}
