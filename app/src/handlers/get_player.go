package handlers

import (
	"net/http"

	"funglejunk.com/kick-api/src/db"
	"github.com/gin-gonic/gin"
)

func (h handler) GetPlayer(c *gin.Context) {
	slug := c.Param("slug")
	p, e1 := db.GetPlayerBySlug(h.DB, slug)
	if _, e := CheckResult(c, p, e1); e == nil {
		c.HTML(http.StatusOK, "player.html", p)
	}
}
