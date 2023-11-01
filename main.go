package main

import (
	"encoding/csv"
<<<<<<< HEAD
	"fmt"
	"log"
=======
	"errors"
	"fmt"
	"log"
	"net/http"
>>>>>>> 4ca3353 (Improvement: Caching .csv file)
	"os"
	"path/filepath"
	"strconv"

	"github.com/gocolly/colly"
)

<<<<<<< HEAD
const relativeHeroesListLocation = "./output/heroesList.csv"

var heroesListLocation string

const heroesUrl = "https://www.bjjheroes.com/a-z-bjj-fighters-list"

func getAbsoluteFilePaths() {
	var err error
	heroesListLocation, err = filepath.Abs(relativeHeroesListLocation)
=======
const bjjHeroesDomain = "www.bjjheroes.com/"
const outputPath = "./output/"
const csvOutputPath = outputPath + "csv/"
const htmlOutputPath = outputPath + "html/"

const relativeAthletesListLocation = csvOutputPath + "athletesList.csv"

var athletesListLocation string

const relativeAthletesHtmlLocation = htmlOutputPath + "athletesList.html"

var athletesHtmlLocation string

const athletesUrl = bjjHeroesDomain + "a-z-bjj-fighters-list"

const forceUpdateHtml = true
const forceUpdateCsv = true

func getAbsoluteFilePaths() {
	var err error
	athletesListLocation, err = filepath.Abs(relativeAthletesListLocation)
	if err != nil {
		log.Fatal(err)
	}
	athletesHtmlLocation, err = filepath.Abs(relativeAthletesHtmlLocation)
>>>>>>> 4ca3353 (Improvement: Caching .csv file)
	if err != nil {
		log.Fatal(err)
	}
}

<<<<<<< HEAD
func readHeroesListCSV() [][]string {
	file, err := os.Open(heroesListLocation)
=======
func readAthletesListCSV() [][]string {
	file, err := os.Open(athletesListLocation)
>>>>>>> 4ca3353 (Improvement: Caching .csv file)
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

<<<<<<< HEAD
func scrapeHeroesList() [][]string {
	// https://www.bjjheroes.com/
	// https://www.scrapingbee.com/blog/web-scraping-go/

	// TODO: Add a short circuit to check if the file was cached recently
	c := colly.NewCollector(
		colly.AllowedDomains("www.bjjheroes.com"),
	)

	heroesData := [][]string{}
	heroesData = append(heroesData, []string{"index", "firstName", "lastName", "nickName", "teamName", "url"})
=======
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
	athletesList = append(athletesList, []string{"index", "firstName", "lastName", "nickName", "teamName", "url"})
>>>>>>> 4ca3353 (Improvement: Caching .csv file)

	c.OnHTML("tbody.row-hover", func(e *colly.HTMLElement) {
		e.ForEach("tr", func(i int, rowEl *colly.HTMLElement) {
			firstName := rowEl.ChildText("td.column-1 > a")
			lastName := rowEl.ChildText("td.column-2 > a")
			nickName := rowEl.ChildText("td.column-3 > a")
			teamName := rowEl.ChildText("td.column-4")
			urlPath := rowEl.ChildAttrs("td.column-1 > a", "href")
<<<<<<< HEAD
			fullUrlPath := "https://www.bjjheroes.com" + urlPath[0]

			heroesData = append(heroesData, []string{strconv.Itoa(i), firstName, lastName, nickName, teamName, fullUrlPath})
		})
	})

	c.Visit(heroesUrl)
	return heroesData
}

func writeHeroesListToCSv(list [][]string) {
	csvFile, err := os.OpenFile(heroesListLocation, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
=======
			fullUrlPath := bjjHeroesDomain + urlPath[0]

			athletesList = append(athletesList, []string{strconv.Itoa(i), firstName, lastName, nickName, teamName, fullUrlPath})
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
			fullUrlPath := bjjHeroesDomain + urlPath[0]

			athletesList = append(athletesList, []string{strconv.Itoa(i), firstName, lastName, nickName, teamName, fullUrlPath})
		})
	})
	fmt.Println("Visiting " + athletesUrl)
	c.Visit(athletesUrl)

	return athletesList
}

func writeAthletesListToCSv(list [][]string) {
	csvFile, err := os.OpenFile(athletesListLocation, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
>>>>>>> 4ca3353 (Improvement: Caching .csv file)
	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}
	csvwriter := csv.NewWriter(csvFile)
	for _, listItem := range list {
		_ = csvwriter.Write(listItem)
	}
	csvwriter.Flush()
	csvFile.Close()
}

<<<<<<< HEAD
func main() {
	fmt.Println("go")
	getAbsoluteFilePaths()
	var list = readHeroesListCSV()
	if list == nil {
		var urls = scrapeHeroesList()
		writeHeroesListToCSv(urls)
	}

=======
func createHeoresList() {
	var urls [][]string
	if athletesListCached() && !forceUpdateHtml {
		urls = scrapeCachedHeroPage()
	} else {
		urls = scrapeAthletesUrl()
	}
	if len(urls) < 2 {
		log.Fatal("Unable to scrape athletes list")
	}
	writeAthletesListToCSv(urls)
}

func main() {
	fmt.Println("go")
	getAbsoluteFilePaths()
	var list = readAthletesListCSV()
	if list == nil || forceUpdateCsv {
		createHeoresList()
	}

	fmt.Println("Updated athletes list can be found at " + athletesListLocation)

>>>>>>> 4ca3353 (Improvement: Caching .csv file)
	fmt.Println("fin")
}
