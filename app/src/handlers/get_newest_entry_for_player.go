package handlers

import (
	"errors"
	"net/http"

	"funglejunk.com/kick-api/src/models"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func (h handler) GetCurrentEntryForPlayer(c *gin.Context) {
	slug := c.Param("slug")
	var entry models.ValueEntry
	if err := h.DB.Where("player_slug = ?", slug).Order("Day desc").First(&entry).Error; err != nil {
		log.Error(err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.AbortWithStatus(http.StatusNotFound)
		}
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, entry)
}
