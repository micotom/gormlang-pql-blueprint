package handlers

import (
	"net/http"

	"funglejunk.com/kick-api/src/db"
	"github.com/gin-gonic/gin"
)

type CreateTeamBody struct {
	Name string `json:"name"`
}

func (h handler) CreateTeam(c *gin.Context) {
	var body CreateTeamBody
	if err := c.BindJSON(&body); err == nil {
		t, e := db.CreateTeam(h.DB, body.Name)
		if _, e := CheckResult(c, t, e); e == nil {
			c.Status(http.StatusCreated)
		}
	} else {
		c.AbortWithStatus(http.StatusBadRequest)
	}
}
