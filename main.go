package main

import (
	"fmt"

	"github.com/SeanDunford/bjjMath/scraper"
)

const forceUpdateCsv = false
const limit = -1

// TODO: Make these flags configurable through binary params and add const forceUpdateHtml = true here

func main() {
	fmt.Println("go")
	if forceUpdateCsv {
		fmt.Println("Force update athletes list csv bc of flag -forceUpdateCsv")
	}
	var list = scraper.ReadAthletesListCSV()
	if list == nil {
		fmt.Println("Athletes list Csv empty or not found")
		list = scraper.CreateHeoresList(limit)
	}

	for _, row := range list[1:] {
		url := row[5]
		escapedName := scraper.ParseEscapedNameFromUrl(url)

		var athleteRecord = scraper.ReadAthleteRecordAsCsvByEscapedName(escapedName)
		if athleteRecord == nil || len(athleteRecord) < 2 {
			fmt.Println("Record for" + escapedName + "not found or empty scraping athlete page")
			list = scraper.CreateAthleteRecord(escapedName, url)
			fmt.Println("Created athlete record for " + escapedName)
			fmt.Println(list)
		}
	}

	fmt.Println("fin")
}
