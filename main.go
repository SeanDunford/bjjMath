package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gocolly/colly"
)

const relativeHeroesListLocation = "./output/heroesList.csv"

var heroesListLocation string

const heroesUrl = "https://www.bjjheroes.com/a-z-bjj-fighters-list"

func getAbsoluteFilePaths() {
	var err error
	heroesListLocation, err = filepath.Abs(relativeHeroesListLocation)
	if err != nil {
		log.Fatal(err)
	}
}

func readHeroesListCSV() [][]string {
	file, err := os.Open(heroesListLocation)
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

func scrapeHeroesList() [][]string {
	// https://www.bjjheroes.com/
	// https://www.scrapingbee.com/blog/web-scraping-go/

	// TODO: Add a short circuit to check if the file was cached recently
	c := colly.NewCollector(
		colly.AllowedDomains("www.bjjheroes.com"),
	)

	heroesData := [][]string{}
	heroesData = append(heroesData, []string{"index", "firstName", "lastName", "nickName", "teamName", "url"})

	c.OnHTML("tbody.row-hover", func(e *colly.HTMLElement) {
		e.ForEach("tr", func(i int, rowEl *colly.HTMLElement) {
			firstName := rowEl.ChildText("td.column-1 > a")
			lastName := rowEl.ChildText("td.column-2 > a")
			nickName := rowEl.ChildText("td.column-3 > a")
			teamName := rowEl.ChildText("td.column-4")
			urlPath := rowEl.ChildAttrs("td.column-1 > a", "href")
			fullUrlPath := "https://www.bjjheroes.com" + urlPath[0]

			heroesData = append(heroesData, []string{strconv.Itoa(i), firstName, lastName, nickName, teamName, fullUrlPath})
		})
	})

	c.Visit(heroesUrl)
	return heroesData
}

func writeHeroesListToCSv(list [][]string) {
	csvFile, err := os.OpenFile(heroesListLocation, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
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

func main() {
	fmt.Println("go")
	getAbsoluteFilePaths()
	var list = readHeroesListCSV()
	if list == nil {
		var urls = scrapeHeroesList()
		writeHeroesListToCSv(urls)
	}

	fmt.Println("fin")
}
