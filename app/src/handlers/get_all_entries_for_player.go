package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"funglejunk.com/kick-api/src/models"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func (h handler) GetAllEntriesForPlayer(c *gin.Context) {
	slug := c.Param("slug")
	var entries []models.ValueEntry
	log.Info(fmt.Printf("got slug: %s", slug))
	if err := h.DB.Where("player_slug = ?", slug).Find(&entries).Error; err != nil {
		log.Error(err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.AbortWithStatus(http.StatusNotFound)
		}
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, entries)
}
