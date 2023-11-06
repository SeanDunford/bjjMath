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

	fmt.Println("fin")
}
