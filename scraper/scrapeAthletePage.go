package scraper

import (
	"encoding/csv"
	"errors"
	"fmt"
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

func matchFromRecordTableRow(i int, rowEl *colly.HTMLElement) *Match {
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

	match := Match{
		SortId:       sortText,
		Opponent:     opponentText,
		OpponentLink: oponentLink,
		winLoss:      wlText,
		Method:       methodText,
		MethodLink:   methodLink,
		Competition:  competitionText,
		Weight:       weightText,
		Stage:        stageText,
		Year:         yearText,
	}

	return &match
}

func ScrapeCachedAthletePage(fileLocation string) AthleteRecord {
	htmlLocation := "file://" + fileLocation

	elements := ScrapeCachedPageProcessChildrenOfTag(
		htmlLocation,
		parentSelector,
		childSelector,
	)

	record := AthleteRecord{}
	for i, el := range elements {
		record = append(record, *matchFromRecordTableRow(i, el))
	}

	return record
}

func ScrapeAthletesPage(athleteUrl string) AthleteRecord {
	elements := ScrapeUrlProcessChildrenOfTag(
		athleteUrl,
		AthletesHtmlLocationFromUrl(athleteUrl),
		bjjHeroesDomain,
		parentSelector,
		childSelector,
	)

	record := AthleteRecord{}
	for i, row := range elements {
		record = append(record, *matchFromRecordTableRow(i, row))
	}

	return record
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

func athleteRecordPageCached(escapedName string) bool {
	athleteRecordLocation := absoluteHtmlOutputPath + "/" + escapedName + ".html"
	if _, err := os.Stat(athleteRecordLocation); errors.Is(err, os.ErrNotExist) {
		return false
	}

	return true
}

func athleteRecordCached(escapedName string) AthleteRecord {
	athleteRecordLocation := absoluteCsvOutputPath + "/" + escapedName + ".csv"
	if _, err := os.Stat(athleteRecordLocation); errors.Is(err, os.ErrNotExist) {
		return nil
	}

	return ReadAthleteRecordAsCsvByEscapedName(escapedName)
}

func CreateAthleteRecord(escapedName string, athleteProfileUrl string, force bool) AthleteRecord {
	getAbsoluteFilePaths()

	var record AthleteRecord = athleteRecordCached(escapedName)
	if force {
		record = ScrapeAthletesPage(athleteProfileUrl)
	} else if len(record) > 1 {
		return record
	} else if athleteRecordPageCached(escapedName) {
		htmlLocation := AthletesHtmlLocationFromEscapedName(escapedName)
		record = ScrapeCachedAthletePage(htmlLocation)
	} else {
		record = ScrapeAthletesPage(athleteProfileUrl)
	}

	// TODO: Find a way to detect empty athlete records based on html content against errors
	// if len(record) < 1 {
	// 	log.Fatal("Unable to scrape athletes list")
	// }

	writeAthletesRecordToCSv(escapedName, record)
	return record
}

func writeAthletesRecordToCSv(escapedName string, record AthleteRecord) {
	athletesRecordLocation := absoluteCsvOutputPath + "/" + escapedName + ".csv"
	fmt.Println("Creating athletes record as csv" + athletesRecordLocation)

	// 0644 means we can read and write the file or directory but other users can only read it.
	csvFile, err := os.OpenFile(athletesRecordLocation, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}
	csvwriter := csv.NewWriter(csvFile)
	for _, match := range record {
		_ = csvwriter.Write(match.toCsvRow())
	}
	csvwriter.Flush()
	csvFile.Close()

	fmt.Println("Updated athletes list can be found at " + athletesRecordLocation)
}

func ReadAthleteRecordAsCsvByEscapedName(athleteName string) AthleteRecord {
	getAbsoluteFilePaths()
	athleteRecordLocation := absoluteCsvOutputPath + "/" + athleteName + ".csv"
	file, err := os.Open(athleteRecordLocation)
	if err != nil {
		return nil
	}
	reader := csv.NewReader(file)
	// TODO: This is broken and not reading my csv input correctly
	csvItems, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	if len(csvItems) < 1 {
		return nil
	}
	records := AthleteRecord{}
	for _, row := range csvItems {
		match := NewMatchFromCsvRow(row)
		records = append(records, *match)
	}
	defer file.Close()
	return records
}
