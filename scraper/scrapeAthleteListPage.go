package scraper

import (
	"encoding/csv"
	"errors"
	"fmt"
	"log"

	"os"
	"strconv"

	UrlResolver "github.com/SeanDunford/bjjMath/urlResolver"
	"github.com/gocolly/colly"
)

func CreateHeoresList(limit int) []Athlete {
	getAbsoluteFilePaths()
	var athletes []Athlete
	if forceUpdateHtml {
		fmt.Println("Force update athletes list html bc of flag -forceUpdateHtml")
		athletes = scrapeAthletesUrl(limit)
	} else if athletesListCached() {
		athletes = scrapeCachedHeroPage(limit)
	} else {
		athletes = scrapeAthletesUrl(limit)
	}

	urlMapping := resolveAthleteUrls(athletes)
	writeUrlMappingToCsv(urlMapping)

	if len(athletes) < 2 {
		log.Fatal("Unable to scrape athletes list")
	}

	writeAthletesListToCSv(athletes)
	return athletes
}

func writeUrlMappingToCsv(urlMapping map[string]string) {
	fmt.Println("Creating urlMapping csv" + urlMappingLocation)
	// 0644 means we can read and write the file or directory but other users can only read it.
	csvFile, err := os.OpenFile(urlMappingLocation, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}
	csvwriter := csv.NewWriter(csvFile)
	_ = csvwriter.Write([]string{"originalUrl", "resolvedUrl"})
	for fullUrlPath, resolvedUrl := range urlMapping {
		_ = csvwriter.Write([]string{fullUrlPath, resolvedUrl})
	}
	fmt.Println("Updated athletes list can be found at " + urlMappingLocation)
	csvwriter.Flush()
	csvFile.Close()
}

func resolveAthleteUrls(athletes []Athlete) map[string]string {
	urlMapping := make(map[string]string)
	for i, a := range athletes[1:] {
		fullUrlPath := a.Url
		resolvedUrl := UrlResolver.ResolveUrl(fullUrlPath)
		urlMapping[fullUrlPath] = resolvedUrl
		fmt.Println(strconv.Itoa(i), ") ", resolvedUrl)
		a.Url = resolvedUrl
	}
	return urlMapping
}

func athletesListCached() bool {
	if _, err := os.Stat(athletesListHtmlLocation); errors.Is(err, os.ErrNotExist) {
		return false
	}

	return true
}

func scrapeCachedHeroPage(limit int) []Athlete {
	athletesList := []Athlete{}

	elements := ScrapeCachedPageProcessChildrenOfTag(
		athletesListHtmlLocation,
		"tbody.row-hover",
		childSelector,
	)

	for i, el := range elements{
		if (i > limit) {
			break
		}
		athletesList = append(athletesList, *athleteFromTableRow(i, el))
	}

	return athletesList
}

func athleteFromTableRow(i int, rowEl *colly.HTMLElement) *Athlete {
	firstName := rowEl.ChildText("td.column-1 > a")
	lastName := rowEl.ChildText("td.column-2 > a")
	nickName := rowEl.ChildText("td.column-3 > a")
	teamName := rowEl.ChildText("td.column-4")
	urlPath := rowEl.ChildAttrs("td.column-1 > a", "href")
	fullUrlPath := "https://" + bjjHeroesDomain + urlPath[0]

	return &Athlete{
		Index:     strconv.Itoa(i),
		FirstName: firstName,
		LastName:  lastName,
		NickName:  nickName,
		TeamName:  teamName,
		Url:       fullUrlPath,
	}
}

func scrapeAthletesUrl(limit int) []Athlete {

	elements := ScrapeUrlProcessChildrenOfTag(
		athletesUrl,
		bjjHeroesDomain,
		"tbody.row-hover",
		childSelector,
		athletesListHtmlLocation,
	)

	athletesList := []Athlete{}
	for i, el := range elements {
		if limit != -1 && i >= limit {
			return athletesList
		}

		athletesList = append(athletesList, *athleteFromTableRow(i, el))
	}

	return athletesList
}

func writeAthletesListToCSv(list []Athlete) {
	fmt.Println("Creating athletes list csv" + athletesListLocation)
	// 0644 means we can read and write the file or directory but other users can only read it.
	csvFile, err := os.OpenFile(athletesListLocation, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}
	csvwriter := csv.NewWriter(csvFile)
	csvwriter.Write(AthleteKeys)
	for _, listItem := range list {
		row := listItem.toCsvRow()
		_ = csvwriter.Write(row)
	}
	csvwriter.Flush()
	csvFile.Close()

	fmt.Println("Updated athletes list can be found at " + athletesListLocation)
}

func ReadAthletesListCSV() []Athlete {
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
	athletes := []Athlete{}
	for _, row := range records[1:] {
		athlete := NewAthleteFromCsvRow(row)
		athletes = append(athletes, *athlete)
	}
	defer file.Close()
	return athletes
}
