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
const parentSelector = "tbody"
const childSelector = "tr"

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

func ScrapeCachedPageProcessChildrenOfTag(
	htmlLocation string,
	parentSelector string,
	childSelector string,
) []*colly.HTMLElement {
	// TODO: return []*colly.HTMLElement is ineficient and should have a callback
	// Unsure on how to create generic interfaces for that

	t := &http.Transport{}
	t.RegisterProtocol("file", http.NewFileTransport(http.Dir("/")))

	c := colly.NewCollector()
	c.WithTransport(t)

	result := []*colly.HTMLElement{}
	c.OnHTML(parentSelector, func(e *colly.HTMLElement) {
		e.ForEach(childSelector, func(i int, rowEl *colly.HTMLElement) {
			result = append(result, rowEl)
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
	) []*colly.HTMLElement {
	// TODO: return []*colly.HTMLElement is ineficient and should have a callback
	// Unsure on how to type that


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

	result := []*colly.HTMLElement{}
	c.OnHTML(parentSelector, func(e *colly.HTMLElement) {
		e.ForEach(childSelector, func(i int, rowEl *colly.HTMLElement) {
			result = append(result, rowEl)
		})
	})
	c.Visit(url)
	return result
}
