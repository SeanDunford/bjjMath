package scraper

import (
	"encoding/csv"
	// "errors"
	// "fmt"
	"log"
	// "net/http"

	"os"
	// "path/filepath"
	// "strconv"
	// UrlResolver "github.com/SeanDunford/bjjMath/urlResolver"
	// "github.com/gocolly/colly"
)

// const bjjHeroesDomain = "www.bjjheroes.com/"
// const outputPath = "./output/"
// const csvOutputPath = outputPath + "csv/"
// const htmlOutputPath = outputPath + "html/"

// const relativeAthletesListLocation = csvOutputPath + "athletesList.csv"

// var athletesListLocation string

// const relativeAthletesHtmlLocation = htmlOutputPath + "athletesList.html"

// var athletesHtmlLocation string

// const athletesUrl = "https://" + bjjHeroesDomain + "a-z-bjj-fighters-list"

// const forceUpdateHtml = true

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
