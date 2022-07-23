package handlers

import (
	"net/http"

	"funglejunk.com/kick-api/src/db"
	"github.com/gin-gonic/gin"
)

func (h handler) AddPlayerToTeam(c *gin.Context) {
	teamSlug := c.Param("team_slug")
	if playerSlug, ok := c.GetQuery("p"); ok {
		t, e := db.GetTeam(h.DB, teamSlug)
		if _, e := CheckResult(c, t, e); e == nil {
			p, e := db.GetPlayerBySlug(h.DB, playerSlug)
			if _, e := CheckResult(c, p, e); e == nil {
				t.Players = append(t.Players, p)
				if e := h.DB.Save(&t).Error; e == nil {
					c.Redirect(301, "/teams/t1")
				} else {
					c.AbortWithStatus(http.StatusInternalServerError)
				}
			} else {
				c.AbortWithStatus(http.StatusBadRequest)
			}
		}
	} else {
		c.AbortWithStatus(http.StatusBadRequest)
	}
}
