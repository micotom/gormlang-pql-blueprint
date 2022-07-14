package handlers

import (
	"net/http"

	"funglejunk.com/kick-api/src/models"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func (h handler) GetAllPlayers(c *gin.Context) {
	var players []models.Player
	if result := h.DB.Find(&players); result.Error != nil {
		log.Error(result.Error)
	}
	c.JSON(http.StatusOK, gin.H{"players": players})
}
