package handlers

import (
	"net/http"

	"funglejunk.com/kick-api/src/db"
	"funglejunk.com/kick-api/src/models"
	"funglejunk.com/kick-api/src/util"
	"github.com/gin-gonic/gin"
)

type PositonsResult struct {
	Keepers    []MinPlayerInfo
	Defenders  []MinPlayerInfo
	Midfielder []MinPlayerInfo
	Attackers  []MinPlayerInfo
}

type MinPlayerInfo struct {
	Name  string
	Slug  string
	Club  string
	Value int
}

func (h handler) GetPlayersByPosition(c *gin.Context) {

	all, e := db.GetAllPlayers(h.DB)
	if _, e := CheckResult(c, all, e); e == nil {
		allByPosition := util.GroupBy(all, func(p models.Player) string {
			return p.Position
		})
		for key, posSlice := range allByPosition {
			posByValue := util.SortBy(posSlice, func(p1 models.Player, p2 models.Player) bool {
				return p1.ValueEntries[len(p1.ValueEntries)-1].Value > p2.ValueEntries[len(p2.ValueEntries)-1].Value
			})
			allByPosition[key] = posByValue
		}

		newMap := make(map[string][]MinPlayerInfo)
		for key, players := range allByPosition {
			newMap[key] = util.Map(players, func(p models.Player) MinPlayerInfo {
				return MinPlayerInfo{
					Name:  p.Name,
					Slug:  p.Slug,
					Club:  p.Club,
					Value: p.ValueEntries[len(p.ValueEntries)-1].Value,
				}
			})
		}

		r := PositonsResult{
			Keepers:    newMap["Torhüter"],
			Defenders:  newMap["Abwehrspieler"],
			Midfielder: newMap["Mittelfeldspieler"],
			Attackers:  newMap["Stürmer"],
		}
		c.HTML(http.StatusOK, "positions.html", r)
		// c.JSON(http.StatusOK, r)
		return
	}
}
