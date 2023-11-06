package scraper

import (
	"encoding/csv"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	// "errors"
	// "fmt"
	"log"
	// "net/http"

	"os"
	// "path/filepath"

	"github.com/gocolly/colly"
)

func ScrapeCachedAthletePage(fileLocation string) []Athlete {
	t := &http.Transport{}
	t.RegisterProtocol("file", http.NewFileTransport(http.Dir("/")))

	c := colly.NewCollector()
	c.WithTransport(t)

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.OnError(func(_ *colly.Response, err error) {
		log.Println("Something went wrong:", err)
	})

	c.OnScraped(func(r *colly.Response) {
		fmt.Println("Finished", r.Request.URL)
	})

	athletesList := []Athlete{}
	athletesList = append(athletesList, []string{
		"SortId",
		"Opponent",
		"OpponentLink",
		"W/L",
		"Method",
		"MethodLink",
		"Competition",
		"Weight",
		"Stage",
		"Year",
	})

	c.OnHTML("tbody", func(e *colly.HTMLElement) {
		e.ForEach("tr", func(i int, rowEl *colly.HTMLElement) {
			sortText := rowEl.ChildText("td:nth-child(1)")
			opponentText := rowEl.ChildText("td:nth-child(2)")
			oponentLink := rowEl.ChildAttr("td:nth-child(2) > a", "href")
			wlText := rowEl.ChildText("td:nth-child(3)")
			methodText := rowEl.ChildText("td:nth-child(4)")
			methodLink := rowEl.ChildAttr("td:nth-child(4) > a", "href")
			competitionText := rowEl.ChildText("td:nth-child(5)")
			weightText := rowEl.ChildText("td:nth-child(6)")
			stageText := rowEl.ChildText("td:nth-child(7)")
			yearText := rowEl.ChildText("td:nth-child(8)")

			athletesList = append(athletesList, []string{
				sortText,
				opponentText,
				oponentLink,
				wlText,
				methodText,
				methodLink,
				competitionText,
				weightText,
				stageText,
				yearText,
			})
		})
	})

	c.Visit("file://" + fileLocation)

	return athletesList
}

func ScrapeAthletesPage(athleteUrl string) []Athlete {
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
		athleteProfileHtmlLocation := AthletesHtmlLocationFromUrl(athleteUrl)
		err := r.Save(athleteProfileHtmlLocation)
		if err != nil {
			log.Fatal(err)
		}
	})

	athletesList := []Athlete{}
	athletesList = append(athletesList, []string{
		"SortId",
		"Opponent",
		"OpponentLink",
		"W/L",
		"Method",
		"MethodLink",
		"Competition",
		"Weight",
		"Stage",
		"Year",
	})

	c.OnHTML("tbody", func(e *colly.HTMLElement) {
		e.ForEach("tr", func(i int, rowEl *colly.HTMLElement) {
			sortText := rowEl.ChildText("td:nth-child(1)")
			opponentText := rowEl.ChildText("td:nth-child(2)")
			oponentLink := rowEl.ChildAttr("td:nth-child(2) > a", "href")
			wlText := rowEl.ChildText("td:nth-child(3)")
			methodText := rowEl.ChildText("td:nth-child(4)")
			methodLink := rowEl.ChildAttr("td:nth-child(4) > a", "href")
			competitionText := rowEl.ChildText("td:nth-child(5)")
			weightText := rowEl.ChildText("td:nth-child(6)")
			stageText := rowEl.ChildText("td:nth-child(7)")
			yearText := rowEl.ChildText("td:nth-child(8)")

			athletesList = append(athletesList, []string{
				sortText,
				opponentText,
				oponentLink,
				wlText,
				methodText,
				methodLink,
				competitionText,
				weightText,
				stageText,
				yearText,
			})
		})
	})
	c.Visit(athleteUrl)

	return athletesList
}

func ParseEscapedNameFromUrl(athleteUrl string) string {
	u, err := url.Parse(athleteUrl)
	if err != nil {
		log.Fatal(err)
	}
	pathPieces := strings.Split(u.Path, "/")
	return pathPieces[len(pathPieces)-1]
}

func AthletesHtmlLocationFromUrl(athleteUrl string) string {
	name := ParseEscapedNameFromUrl(athleteUrl)
	return absoluteHtmlOutputPath + "/" + name + ".html"
}

func AthletesHtmlLocationFromEscapedName(escapedName string) string {
	return absoluteHtmlOutputPath + "/" + escapedName + ".html"
}

func athleteRecordCached(escapedName string) bool {
	athleteRecordLocation := absoluteHtmlOutputPath + "/" + escapedName + ".html"
	if _, err := os.Stat(athleteRecordLocation); errors.Is(err, os.ErrNotExist) {
		return false
	}

	return true
}

func CreateAthleteRecord(escapedName string, athleteProfileUrl string) []Athlete {
	getAbsoluteFilePaths()
	var record []Athlete
	if athleteRecordCached(escapedName) {
		htmlLocation := AthletesHtmlLocationFromEscapedName(escapedName)
		record = ScrapeCachedAthletePage(htmlLocation)
	} else {
		record = ScrapeAthletesPage(athleteProfileUrl)
	}

	if len(record) < 1 {
		log.Fatal("Unable to scrape athletes list")
	}

	writeAthletesRecordToCSv(escapedName, record)
	return record
}

func writeAthletesRecordToCSv(escapedName string, record []Athlete) {
	athletesRecordLocation := absoluteCsvOutputPath + "/" + escapedName + ".csv"
	fmt.Println("Creating athletes record as csv" + athletesRecordLocation)

	// 0644 means we can read and write the file or directory but other users can only read it.
	csvFile, err := os.OpenFile(athletesRecordLocation, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}
	csvwriter := csv.NewWriter(csvFile)
	for _, listItem := range record {
		_ = csvwriter.Write(listItem)
	}
	csvwriter.Flush()
	csvFile.Close()

	fmt.Println("Updated athletes list can be found at " + athletesRecordLocation)
}

func ReadAthleteRecordAsCsvByEscapedName(athleteName string) []Athlete {
	getAbsoluteFilePaths()
	athleteRecordLocation := absoluteCsvOutputPath + "/" + athleteName + ".csv"
	file, err := os.Open(athleteRecordLocation)
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
