package handlers

import (
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"funglejunk.com/kick-api/src/db"
	"funglejunk.com/kick-api/src/models"
	"github.com/gin-gonic/gin"
	"github.com/gocolly/colly"
	"github.com/gosimple/slug"
	"gorm.io/datatypes"

	log "github.com/sirupsen/logrus"
)

var intRegex *regexp.Regexp = regexp.MustCompile("[^-?(?:\\d+)+]")
var floatRegex *regexp.Regexp = regexp.MustCompile("[^-?(?:\\d+.?)+]")

func parseMoney(str string) int {
	numStr := intRegex.ReplaceAllString(str, "")
	m, err := strconv.Atoi(numStr)
	if err != nil {
		panic(err)
	}
	return m
}

func parsePercent(str string) float32 {
	percString := floatRegex.ReplaceAllString(strings.Replace(str, ",", ".", 1), "")
	v, err := strconv.ParseFloat(percString, 32)
	if err != nil {
		panic(err)
	}
	return float32(v)
}

func parsePoints(str string) int {
	numStr := intRegex.ReplaceAllString(str, "")
	p, err := strconv.Atoi(numStr)
	if err != nil {
		panic(err)
	}
	return p
}

func shouldAddEntry(new models.ValueEntry, olds []models.ValueEntry) bool {
	nY, nM, nD := time.Time(new.Day).Date()
	for _, entry := range olds {
		y, m, d := time.Time(entry.Day).Date()
		add := y != nY && m != nM && d != nD
		if !add {
			return false
		}
	}
	return true
}

func (h handler) DoScrape(c *gin.Context) {
	cy := colly.NewCollector()

	cy.Limit(&colly.LimitRule{
		DomainGlob:  "ligainsider.de/*",
		Delay:       1 * time.Second,
		RandomDelay: 1 * time.Second,
	})

	cy.OnHTML("table[class=\"annotated-list ranking_table noten_table mvr_table table-striped\"]",
		func(el *colly.HTMLElement) {
			el.ForEach("tbody", func(_ int, body *colly.HTMLElement) {
				body.ForEach("tr", func(_ int, row *colly.HTMLElement) {
					var p = &models.Player{}
					var vE = &models.ValueEntry{
						Day: datatypes.Date(time.Now()),
					}
					row.ForEach("td", func(index int, column *colly.HTMLElement) {
						switch index {
						case 2:
							p.Name = strings.TrimSpace(column.Text)
							p.Slug = slug.MakeLang(p.Name, "de")
						case 3:
							p.Club = column.Text
						case 4:
							p.Position = column.Text
						case 5:
							p.TotalPoints = parsePoints(column.Text)
						case 6:
							vE.Value = parseMoney(column.Text)
						case 7:
							vE.RaisePerc = parsePercent(column.Text)
						case 8:
							vE.RaiseDiff = parseMoney(column.Text)
						}
					})

					h.DB.FirstOrCreate(p).Preload("ValueEntries")
					currentEntries, err := db.GetEntriesForPlayer(h.DB, p.Slug)

					if err != nil {
						panic(err)
					}

					if shouldAddEntry(*vE, currentEntries) {
						allEntries := append(p.ValueEntries, *vE)
						p.ValueEntries = allEntries
						h.DB.Save(p)
					} else {
						log.Info("Duplicate day entry - ignore")
					}
				})
			})
		})

	cy.Visit("https://www.ligainsider.de/stats/kickbase/marktwerte/tag/gewinner/")
	cy.Visit("https://www.ligainsider.de/stats/kickbase/marktwerte/tag/verlierer/")

	result, err := db.GetAllPlayersWithEntries(h.DB)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
	} else {
		c.JSON(http.StatusOK, result)
	}
}
