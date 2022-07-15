package handlers

import (
	"net/http"

	"funglejunk.com/kick-api/src/db"
	"github.com/gin-gonic/gin"
)

func (h handler) GetAllEntriesForPlayer(c *gin.Context) {
	slug := c.Param("slug")
	r1, e1 := db.GetEntriesForPlayer(h.DB, slug)
	if r, e := CheckResult(c, r1, e1); e == nil {
		c.JSON(http.StatusOK, r)
	}
}
