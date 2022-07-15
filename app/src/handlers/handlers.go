package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type handler struct {
	DB *gorm.DB
}

func New(db *gorm.DB) handler {
	return handler{db}
}

func CheckResult(c *gin.Context, result interface{}, err error) (interface{}, error) {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.AbortWithStatus(http.StatusNotFound)
		return nil, err
	} else if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return nil, err
	} else {
		return result, nil
	}
}
