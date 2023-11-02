package scraper

import (
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	"net/http"

	"os"
	"path/filepath"
	"strconv"

	UrlResolver "github.com/SeanDunford/bjjMath/urlResolver"
	"github.com/gocolly/colly"
)

const bjjHeroesDomain = "www.bjjheroes.com/"
const outputPath = "./output/"
const csvOutputPath = outputPath + "csv/"
const htmlOutputPath = outputPath + "html/"

const relativeAthletesListLocation = csvOutputPath + "athletesList.csv"

var athletesListLocation string

const relativeAthletesHtmlLocation = htmlOutputPath + "athletesList.html"

var athletesHtmlLocation string

const athletesUrl = "https://" + bjjHeroesDomain + "a-z-bjj-fighters-list"

const forceUpdateHtml = true

func CreateHeoresList() {
	getAbsoluteFilePaths()
	var urls [][]string
	if forceUpdateHtml {
		fmt.Println("Force update athletes list html bc of flag -forceUpdateHtml")
		urls = scrapeAthletesUrl()
	} else if athletesListCached() {
		urls = scrapeCachedHeroPage()
	} else {
		urls = scrapeAthletesUrl()
	}

	if len(urls) < 2 {
		log.Fatal("Unable to scrape athletes list")
	}

	writeAthletesListToCSv(urls)
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

func getAbsoluteFilePaths() {
	var err error
	athletesListLocation, err = filepath.Abs(relativeAthletesListLocation)
	if err != nil {
		log.Fatal(err)
	}
	athletesHtmlLocation, err = filepath.Abs(relativeAthletesHtmlLocation)
	if err != nil {
		log.Fatal(err)
	}
}

func athletesListCached() bool {
	if _, err := os.Stat(athletesHtmlLocation); errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true
}

func scrapeCachedHeroPage() [][]string {
	t := &http.Transport{}
	t.RegisterProtocol("file", http.NewFileTransport(http.Dir("/")))

	c := colly.NewCollector()
	c.WithTransport(t)

	athletesList := [][]string{}
	// TODO: Convert this to a type
	athletesList = append(athletesList, []string{"index", "firstName", "lastName", "nickName", "teamName", "url"})

	c.OnHTML("tbody.row-hover", func(e *colly.HTMLElement) {
		e.ForEach("tr", func(i int, rowEl *colly.HTMLElement) {
			firstName := rowEl.ChildText("td.column-1 > a")
			lastName := rowEl.ChildText("td.column-2 > a")
			nickName := rowEl.ChildText("td.column-3 > a")
			teamName := rowEl.ChildText("td.column-4")
			urlPath := rowEl.ChildAttrs("td.column-1 > a", "href")
			fullUrlPath := "https://" + bjjHeroesDomain + urlPath[0]
			resolvedUrl := UrlResolver.ResolveUrl(fullUrlPath)
			fmt.Println(resolvedUrl)

			athletesList = append(athletesList, []string{strconv.Itoa(i), firstName, lastName, nickName, teamName, resolvedUrl})
		})
	})

	c.Visit("file://" + athletesHtmlLocation)

	return athletesList
}

func scrapeAthletesUrl() [][]string {
	c := colly.NewCollector(
	// colly.AllowedDomains(bjjHeroesDomain),
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
	athletesList := [][]string{}
	athletesList = append(athletesList, []string{"index", "firstName", "lastName", "nickName", "teamName", "url"})

	c.OnHTML("tbody.row-hover", func(e *colly.HTMLElement) {
		e.ForEach("tr", func(i int, rowEl *colly.HTMLElement) {
			firstName := rowEl.ChildText("td.column-1 > a")
			lastName := rowEl.ChildText("td.column-2 > a")
			nickName := rowEl.ChildText("td.column-3 > a")
			teamName := rowEl.ChildText("td.column-4")
			urlPath := rowEl.ChildAttrs("td.column-1 > a", "href")
			fullUrlPath := "https://" + bjjHeroesDomain + urlPath[0]
			resolvedUrl := UrlResolver.ResolveUrl(fullUrlPath)
			fmt.Println(resolvedUrl)

			athletesList = append(athletesList, []string{strconv.Itoa(i), firstName, lastName, nickName, teamName, resolvedUrl})
		})
	})
	fmt.Println("Visiting " + athletesUrl)
	c.Visit(athletesUrl)

	return athletesList
}

func writeAthletesListToCSv(list [][]string) {
	fmt.Println("Creating athletes list csv" + athletesListLocation)
	csvFile, err := os.OpenFile(athletesListLocation, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}
	csvwriter := csv.NewWriter(csvFile)
	for _, listItem := range list {
		_ = csvwriter.Write(listItem)
	}
	csvwriter.Flush()
	csvFile.Close()

	fmt.Println("Updated athletes list can be found at " + athletesListLocation)
}