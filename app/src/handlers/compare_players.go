package handlers

import (
	"net/http"
	"strings"

	"funglejunk.com/kick-api/src/db"
	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
)

type ComparisonResult struct {
	PlayerResults []ComparisonPlayerResult
	Winner        string  `json:"winner"`
	WinningDiff   float32 `json:"diff"`
}

type ComparisonPlayerResult struct {
	PlayerName string          `json:"name"`
	Value      ComparisonValue `json:"value"`
}

type ComparisonValue struct {
	Diff float32        `json:"diff"`
	Date datatypes.Date `json:"date"`
}

func (h handler) ComparePlayers(c *gin.Context) {
	if queryParam, ok := c.GetQuery("ps"); ok {
		names := strings.Split(queryParam, ",")
		result := ComparisonResult{}

		for _, n := range names {
			pResult := ComparisonPlayerResult{}
			p, _ := db.GetPlayerBySlug(h.DB, n)
			entry, _ := db.GetCurrentEntry(h.DB, n)
			v := ComparisonValue{
				Diff: float32(entry.RaiseDiff),
				Date: entry.Day,
			}
			pResult.Value = v
			pResult.PlayerName = p.Name
			result.PlayerResults = append(result.PlayerResults, pResult)
		}

		var winner string
		var highestVal float32
		for i, r := range result.PlayerResults {
			if i == 0 {
				winner = r.PlayerName
				highestVal = r.Value.Diff
			} else {
				if r.Value.Diff > highestVal {
					highestVal = r.Value.Diff
					winner = r.PlayerName
				}
			}
		}
		result.Winner = winner
		result.WinningDiff = highestVal

		c.JSON(http.StatusOK, result)
	} else {
		c.AbortWithStatus(http.StatusBadRequest)
	}
}
