package scraper

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"

	"github.com/gocolly/colly"
)

const bjjHeroesDomain = "www.bjjheroes.com"
const outputPath = "./output/"
const csvOutputPath = outputPath + "csv/"
const htmlOutputPath = outputPath + "html/"
const relativeAthletesListLocation = csvOutputPath + "athletesList.csv"
const relativeUrlMappingLocation = csvOutputPath + "urlMapping.csv"
const relativeAthletesListHtmlLocation = htmlOutputPath + "athletesList.html"
const athletesUrl = "https://" + bjjHeroesDomain + "/a-z-bjj-fighters-list"

var absoluteHtmlOutputPath string
var absoluteCsvOutputPath string
var athletesListLocation string
var urlMappingLocation string
var athletesListHtmlLocation string

const forceUpdateHtml = false

func getAbsoluteFilePaths() {
	var err error
	athletesListLocation, err = filepath.Abs(relativeAthletesListLocation)
	if err != nil {
		log.Fatal(err)
	}
	athletesListHtmlLocation, err = filepath.Abs(relativeAthletesListHtmlLocation)
	if err != nil {
		log.Fatal(err)
	}
	urlMappingLocation, err = filepath.Abs(relativeUrlMappingLocation)
	if err != nil {
		log.Fatal(err)
	}

	absoluteHtmlOutputPath, err = filepath.Abs(htmlOutputPath)
	if err != nil {
		log.Fatal(err)
	}
	absoluteCsvOutputPath, err = filepath.Abs(csvOutputPath)
	if err != nil {
		log.Fatal(err)
	}
}

type parseable interface {
	[]Athlete | AthleteRecord
	processAthleteListTableChild()
}

func ScrapeCachedPageProcessChildrenOfTag[T parseable](
	htmlLocation string,
	parentSelector string,
	childSelector string,
	callback func(int, *colly.HTMLElement, []T) []T,
) []T {

	t := &http.Transport{}
	t.RegisterProtocol("file", http.NewFileTransport(http.Dir("/")))

	c := colly.NewCollector()
	c.WithTransport(t)

	result := make([]T, 0)
	c.OnHTML(parentSelector, func(e *colly.HTMLElement) {
		e.ForEach(childSelector, func(i int, rowEl *colly.HTMLElement) {
			// result = append(result, callback(i, rowEl, result))
		})
	})

	c.Visit("file://" + athletesListHtmlLocation)
	return result
}

func ScrapeUrlProcessChildrenOfTag(
	url string,
	allowedDomain string,
	parentSelector string,
	childSelector string,
	htmlLocation string,
	callback func(int, *colly.HTMLElement)) {

	c := colly.NewCollector(
		colly.AllowedDomains(allowedDomain),
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
		err := r.Save(htmlLocation)
		if err != nil {
			log.Fatal(err)
		}
	})

	c.OnHTML(parentSelector, func(e *colly.HTMLElement) {
		e.ForEach(childSelector, callback)
	})
	c.Visit(url)

}
