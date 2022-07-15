package handlers

import (
	"net/http"

	"funglejunk.com/kick-api/src/db"
	"github.com/gin-gonic/gin"
)

func (h handler) GetAllPlayers(c *gin.Context) {
	r1, e1 := db.GetAllPlayers(h.DB)
	if r, e := CheckResult(c, r1, e1); e == nil {
		c.JSON(http.StatusOK, r)
	}
}
