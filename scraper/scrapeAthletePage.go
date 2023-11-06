package scraper

import (
	"encoding/csv"
	"fmt"

	// "errors"
	// "fmt"
	"log"
	// "net/http"

	"os"
	// "path/filepath"

	// UrlResolver "github.com/SeanDunford/bjjMath/urlResolver"
	"github.com/gocolly/colly"
)

func scrapeAthletesPage(athleteUrl string) {
	c := colly.NewCollector(
		colly.AllowedDomains(bjjHeroesDomain),
	)

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.OnError(func(_ *colly.Response, err error) {
		log.Println("Something went wrong:", err)
	})

	c.OnScraped(func(r *colly.Response) {
		fmt.Println("Finished", r.Request.URL)
	})

	c.OnResponse(func(r *colly.Response) {
		err := r.Save(athletesHtmlLocation)
		if err != nil {
			log.Fatal(err)
		}
	})

	c.OnHTML("tbody.row-hover", func(e *colly.HTMLElement) {
		e.ForEach("tr", func(i int, rowEl *colly.HTMLElement) {
			Opponent	W/L	Method	Competition	Weight	Stage	Year
		})
	})
	c.Visit(athleteUrl)

	return athletesList
}

func ReadAthletesListCSV() [][]string {
	getAbsoluteFilePaths()
	file, err := os.Open(athletesListLocation)
	if err != nil {
		return nil
	}
	reader := csv.NewReader(file)
	// TODO: This is broken and not reading my csv input correctly
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	if len(records) < 1 {
		return nil
	}
	defer file.Close()
	return records
}
