package dbo

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

var hrefRegex, _ = regexp.Compile("=([^=]+).htm$")

func hrefToBOMID(href string) (string, error) {
	res := hrefRegex.FindStringSubmatch(href)
	if len(res) < 2 {
		return "", fmt.Errorf("Unable to get id from url %s", href)
	}
	return res[1], nil
}

func grossStringToInt(gross string) (int, error) {
	gross = strings.Replace(gross, ",", "", -1)
	gross = strings.Replace(gross, "$", "", -1)
	g, err := strconv.Atoi(gross)
	return g, err
}

type Details struct {
	ID    int
	BOMID string
	Title string
}

type DailyGross struct {
	ID    int
	Date  time.Time
	Gross int
}

type DayGross struct {
	ID    int
	BOMID string
	Title string
	Gross int
}

func GetDailyGross(day time.Time) ([]DayGross, error) {
	records := make([]DayGross, 0)

	selector := "#body > center > center > table > tbody > tr:nth-child(2) > td > table > tbody > tr"
	c := colly.NewCollector()
	c.OnHTML(selector, func(e *colly.HTMLElement) {
		// e.Request.Visit(e.Attr("href"))
		if e.ChildText("td:nth-child(3)") == "Title (Click to View)" {
			return
		}

		BOMID, err := hrefToBOMID(e.ChildAttr("td:nth-child(3) a", "href"))
		if err != nil {
			log.Println("Failed to get id:", err)
			return
		}

		gross, err := grossStringToInt(e.ChildText("td:nth-child(5)"))
		if err != nil {
			log.Println("Failed to get gross:", err)
			return
		}

		records = append(records, DayGross{
			BOMID: BOMID,
			Title: e.ChildText("td:nth-child(3)"),
			Gross: gross,
		})
	})

	c.OnRequest(func(r *colly.Request) {
		log.Println("Visiting", r.URL)
	})

	c.Visit(
		"https://www.boxofficemojo.com/daily/chart/?view=1day&sortdate=" + day.Format("2006-01-02"),
	)

	return records, nil
}
