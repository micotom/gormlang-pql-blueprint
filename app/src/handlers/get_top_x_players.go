package handlers

import (
	"net/http"
	"strconv"
	"time"

	"funglejunk.com/kick-api/src/db"
	"funglejunk.com/kick-api/src/models"
	"funglejunk.com/kick-api/src/util"
	"github.com/gin-gonic/gin"
	"github.com/micotom/gfuncs"
)

type ShortDate struct {
	day   int
	month time.Month
	year  int
}

type TopResultX struct {
	Name     string `json:"player"`
	Slug     string
	Diff     int `json:"diff"`
	Position int `json:"position"`
}

func (h handler) GetTopPlayersX(c *gin.Context) {
	if queryParam, ok := c.GetQuery("days"); ok {
		daysMinus, e := strconv.Atoi(queryParam)
		if e != nil {
			c.AbortWithStatus(http.StatusBadRequest)
		} else {
			ps, e := db.GetAllPlayers(h.DB)
			if _, e := CheckResult(c, ps, e); e == nil {
				y, m, d := time.Now().Date()
				last := ShortDate{
					day: d, month: m, year: y,
				}
				lY, lM, lD := time.Now().AddDate(0, 0, -daysMinus).Date()
				first := ShortDate{
					day: lD, month: lM, year: lY,
				}

				sort := gfuncs.SortBy(ps, func(p models.Player) int {
					lSum := util.SumByInt(p.ValueEntries, func(ve models.ValueEntry) int {
						if isInDateRange(shortDateFromTime(time.Time(ve.Day)), first, last) {
							return ve.RaiseDiff
						}
						return 0
					})
					return lSum
				})

				r := []TopResultX{}
				for i, p := range sort {
					diff := util.SumByInt(p.ValueEntries, func(ve models.ValueEntry) int {
						return ve.RaiseDiff
					})
					r = append(r, TopResultX{
						Name:     p.Name,
						Slug:     p.Slug,
						Diff:     diff,
						Position: i,
					})
				}

				c.JSON(http.StatusOK, r)
			}
		}
	} else {
		c.AbortWithStatus(http.StatusBadRequest)
	}
}

func shortDateFromTime(t time.Time) ShortDate {
	y, m, d := t.Date()
	return ShortDate{
		year: y, month: m, day: d,
	}
}

func (sh ShortDate) toTime() time.Time {
	return time.Date(sh.year, sh.month, sh.day, 0, 0, 0, 0, time.UTC)
}

func (sh ShortDate) equal(other ShortDate) bool {
	return sh.year == other.year && sh.month == other.month && sh.year == other.year
}

func isInDateRange(d ShortDate, first ShortDate, last ShortDate) bool {
	return d.equal(first) || d.equal(last) || (d.toTime().After(first.toTime()) && d.toTime().Before(first.toTime()))
}
