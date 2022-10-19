package handlers

import (
	"net/http"
	"time"

	"funglejunk.com/kick-api/src/db"
	"funglejunk.com/kick-api/src/models"
	"funglejunk.com/kick-api/src/util"
	"github.com/gin-gonic/gin"
)

type PlayerValue struct {
	PlayerName     string
	PlayerSlug     string
	TotalPoints    int
	Price          int
	PricePerPoints float32
	PointRank      int
	Coeff          float32
}

type PriceValueResponse struct {
	PlayerPrices      []PlayerValue
	PlayerPriceValues []PlayerValue
	PlayerCoeffs      []PlayerValue
}

func (h handler) GetPriceValue(c *gin.Context) {
	allPlayers, err := db.GetAllPlayers(h.DB)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	pvs := make([]PlayerValue, 0, len(allPlayers))

	for _, p := range allPlayers {
		// not invalid
		pricePerPoint := models.GetPlayerPointsPerPrice(p)
		if pricePerPoint == -1 || pricePerPoint < 0 {
			continue
		}
		// skip if older than 2 days
		lastEntry := p.ValueEntries[len(p.ValueEntries)-1]
		now := time.Now()
		diff := now.Sub(time.Time(lastEntry.Day))
		if int64(diff.Hours()/24) > 3 {
			continue
		}

		pvs = append(pvs, PlayerValue{
			PlayerName:     p.Name,
			PlayerSlug:     p.Slug,
			Price:          p.ValueEntries[len(p.ValueEntries)-1].Value,
			PricePerPoints: models.GetPlayerPointsPerPrice(p),
			TotalPoints:    p.TotalPoints,
		})
	}

	playersByTotalPoints := util.SortBy2(pvs, func(pv PlayerValue) int {
		return pv.TotalPoints
	})
	util.Reverse(playersByTotalPoints)

	playersByPriceValue := util.SortBy2(pvs, func(pv PlayerValue) float32 {
		return pv.PricePerPoints
	})

	for i1, _ := range playersByPriceValue {
		p := &playersByPriceValue[i1]
		p.PointRank = -1
		for i, p2 := range playersByTotalPoints {
			if p.PlayerSlug == p2.PlayerSlug {
				p.PointRank = i + 1
				break
			}
		}
		p.Coeff = float32(i1) + 1.75*float32(p.PointRank)
	}

	playersByCoeff := util.SortBy2(playersByPriceValue, func(pv PlayerValue) float32 {
		return pv.Coeff
	})

	c.HTML(http.StatusOK, "pricevalues.html", PriceValueResponse{
		PlayerPrices: playersByTotalPoints, PlayerPriceValues: playersByPriceValue, PlayerCoeffs: playersByCoeff,
	})
}
