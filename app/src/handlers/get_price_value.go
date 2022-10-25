package handlers

import (
	"net/http"
	"time"

	"funglejunk.com/kick-api/src/db"
	"funglejunk.com/kick-api/src/models"
	"funglejunk.com/kick-api/src/util"
	"github.com/gin-gonic/gin"
	"github.com/micotom/gfuncs"
	log "github.com/sirupsen/logrus"
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

func indexOf(slice []PlayerValue, p PlayerValue) int {
	for i, p2 := range slice {
		if p.PlayerSlug == p2.PlayerSlug {
			return i
		}
	}
	return -1
}

func (h handler) GetPriceValue(c *gin.Context) {
	allPlayers, err := db.GetAllPlayers(h.DB)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	pvs := make([]PlayerValue, 0, len(allPlayers))

	for _, p := range allPlayers {
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

	playersByTotalPoints := gfuncs.SortBy(pvs, func(pv PlayerValue) int {
		return pv.TotalPoints
	})
	util.Reverse(playersByTotalPoints)

	playersByPriceValue := util.SortBy2(pvs, func(pv PlayerValue) float32 {
		return pv.PricePerPoints
	})

	log.Info("best points: ", playersByTotalPoints[0].PlayerName)

	maxLen := 200
	coeffPlayers := make([]PlayerValue, maxLen)

	for i := 0; i < len(playersByTotalPoints) && i < 200; i++ {
		priceValueIndex := indexOf(playersByPriceValue, playersByTotalPoints[i])
		p := &playersByPriceValue[priceValueIndex]
		p.PointRank = i + 1
		p.Coeff = float32(priceValueIndex) + 2.*float32(p.PointRank)
		coeffPlayers[i] = *p
	}

	coeffPlayers = util.SortBy2(coeffPlayers, func(pv PlayerValue) float32 {
		return pv.Coeff
	})

	c.HTML(http.StatusOK, "pricevalues.html", PriceValueResponse{
		PlayerPrices: playersByTotalPoints, PlayerPriceValues: playersByPriceValue, PlayerCoeffs: coeffPlayers,
	})
}
