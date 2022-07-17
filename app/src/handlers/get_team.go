package handlers

import (
	"net/http"

	"funglejunk.com/kick-api/src/db"
	"github.com/gin-gonic/gin"
)

func (h handler) GetTeam(c *gin.Context) {
	slug := c.Param("slug")
	t, e := db.GetTeam(h.DB, slug)
	if _, e := CheckResult(c, t, e); e == nil {
		c.HTML(http.StatusOK, "team.html", t)
	}
}
