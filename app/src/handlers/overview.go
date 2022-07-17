package handlers

import (
	"net/http"

	"funglejunk.com/kick-api/src/db"
	"funglejunk.com/kick-api/src/models"
	"github.com/gin-gonic/gin"
)

func (h handler) GetOverview(c *gin.Context) {

	ps, e := db.GetAllPlayers(h.DB)
	if _, e := CheckResult(c, ps, e); e == nil {
		mapped := make(map[string][]models.Player)
		for _, e := range ps {
			c := e.Club
			if val, present := mapped[c]; present {
				val = append(val, e)
				mapped[c] = val
			} else {
				mapped[c] = []models.Player{}
				mapped[c] = append(mapped[c], e)
			}
		}
		c.HTML(http.StatusOK, "overview.html", mapped)
	}

}
