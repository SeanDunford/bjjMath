package scraper

import "github.com/gocolly/colly"

var AthleteKeys = []string{ // How to make const?
	"index",
	"firstName",
	"lastName",
	"nickName",
	"teamName",
	"url",
}

func (a Athlete) toCsvRow() []string {
	return []string{
		a.Index,
		a.FirstName,
		a.LastName,
		a.NickName,
		a.TeamName,
		a.Url,
	}
}

type Athlete struct {
	Index     string
	FirstName string
	LastName  string
	NickName  string
	TeamName  string
	Url       string
}

func NewAthleteFromCsvRow(csvRow []string) *Athlete {
	return &Athlete{
		Index:     csvRow[0],
		FirstName: csvRow[1],
		LastName:  csvRow[2],
		NickName:  csvRow[3],
		TeamName:  csvRow[4],
		Url:       csvRow[5],
	}
}

func (a Athlete) processAthleteListTableChild(i int, rowEl *colly.HTMLElement, athletesList []Athlete) []Athlete {
	return append(athletesList, *athleteFromTableRow(i, rowEl))
}
