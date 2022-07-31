package handlers

import (
	"fmt"
	"net/http"
	"os"
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

var scrapeDomain = os.Getenv("SCRAPE_DOMAIN")
var scrapeLink1 = os.Getenv("SCRAPE_LINK_WINNERS")
var scrapeLink2 = os.Getenv("SCRAPE_LINK_LOSERS")

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
	fmt.Println(fmt.Sprintf("check: %d, %d, %d", nY, int(nM), nD))
	for _, entry := range olds {
		y, m, d := time.Time(entry.Day).Date()
		isDuplicate := y == nY && m == nM && d == nD
		fmt.Println(fmt.Sprintf("\tagainst: %d, %d, %d", y, int(m), d))
		if isDuplicate {
			return false
		}
	}
	return true
}

func (h handler) DoScrapeRemote(c *gin.Context) {
	h.DoScrape()
	c.Status(http.StatusOK)
}

func (h handler) DoScrape() {
	cy := colly.NewCollector()

	cy.Limit(&colly.LimitRule{
		DomainGlob:  scrapeDomain,
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

					dbP, e := db.GetOrCreatePlayer(h.DB, *p)
					if e != nil {
						panic(e)
					}

					if shouldAddEntry(*vE, dbP.ValueEntries) {
						allEntries := append(dbP.ValueEntries, *vE)
						dbP.ValueEntries = allEntries
						if e := db.SavePlayer(h.DB, dbP); e != nil {
							panic(e)
						}
					} else {
						log.Info("Duplicate")
					}
				})
			})
		})

	cy.Visit(scrapeLink1)
	cy.Visit(scrapeLink2)

}
