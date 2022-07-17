package handlers

import (
	"net/http"

	"funglejunk.com/kick-api/src/db"
	"github.com/gin-gonic/gin"
)

type AddPlayerToTeamBody struct {
	PlayerSlug string `json:"player"`
}

func (h handler) AddPlayerToTeam(c *gin.Context) {
	teamSlug := c.Param("slug")
	var body AddPlayerToTeamBody
	if err := c.BindJSON(&body); err == nil {
		t, e := db.GetTeam(h.DB, teamSlug)
		if _, e := CheckResult(c, t, e); e == nil {
			p, e := db.GetPlayerBySlug(h.DB, body.PlayerSlug)
			if _, e := CheckResult(c, p, e); e == nil {
				t.Players = append(t.Players, p)
				if e := h.DB.Save(t).Error; e == nil {
					c.Status(http.StatusCreated)
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
