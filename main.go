package main

import (
	"fmt"

	"github.com/SeanDunford/bjjMath/db"
	"github.com/SeanDunford/bjjMath/graph"
	"github.com/SeanDunford/bjjMath/scraper"
)

const forceUpdateAthleteListCsv = false
const forceUpdateAthleteRecordCsv = true
const LimitOfAthletes = -1

// const limitOfAthleteRecords = -1 TODO: Implement
const forceUpateTexOnly = true

// TODO: Make these flags configurable through binary params and add const forceUpdateHtml = true here
func main() {
	fmt.Println("go")

	db.ConnectToDb()
	dbFile := "db/bjjMath.db"

	// if forceUpateTexOnly {
	// 	const escapedName = "aaron-johnson"
	// 	const texUrl = "https://www.bjjheroes.com/bjj-fighters/aaron-johnson"
	// 	record := scraper.CreateAthleteRecord(escapedName, texUrl, forceUpateTexOnly)
	// 	fmt.Println("Created athlete record for " + escapedName)
	// 	fmt.Println(record)
	// 	fmt.Println("fin")
	// 	return
	// }

	if forceUpdateAthleteListCsv {
		fmt.Println("Force update athletes list csv bc of flag -forceUpdateCsv")
	}
	var athletes = scraper.ReadAthletesListCSV()
	if athletes == nil || len(athletes) < 1 {
		fmt.Println("Athletes list Csv empty or not found")
		athletes = scraper.CreateHeoresList(LimitOfAthletes)
	}

	for _, athlete := range athletes {
		escapedName := scraper.ParseEscapedNameFromUrl(athlete.Url)

		var athleteRecord = scraper.ReadAthleteRecordAsCsvByEscapedName(escapedName)
		if forceUpdateAthleteRecordCsv {
			fmt.Println("Force Update flag detected creating csv Record for" + escapedName)
			record := scraper.CreateAthleteRecord(escapedName, athlete.Url, false)
			fmt.Println("Created athlete record for " + escapedName)
			fmt.Println(record)
		} else if athleteRecord == nil || len(athleteRecord) < 2 {
			fmt.Println("Record for " + escapedName + " not found or empty scraping athlete page")
			record := scraper.CreateAthleteRecord(escapedName, athlete.Url, false)
			fmt.Println("Created athlete record for " + escapedName)
			fmt.Println(record)
		}

		graph.AddAthlete(athlete, dbFile)
	}

	fmt.Println("fin")
}
