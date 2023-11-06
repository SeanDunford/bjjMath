package scraper

var MatchKeys = []string{
	"sortId",
	"opponent",
	"opponentLink",
	"W/L",
	"method",
	"methodLink",
	"competition",
	"weight",
	"stage",
	"year",
}

type Match struct {
	SortId       string
	Opponent     string
	OpponentLink string
	winLoss      string
	Method       string
	MethodLink   string
	Competition  string
	Weight       string
	Stage        string
	Year         string
}

type AthleteRecord []Match

func NewMatchFromCsvRow(csvRow []string) *Match {
	return &Match{
		SortId:       csvRow[0],
		Opponent:     csvRow[1],
		OpponentLink: csvRow[2],
		winLoss:      csvRow[3],
		Method:       csvRow[4],
		MethodLink:   csvRow[5],
		Competition:  csvRow[6],
		Weight:       csvRow[7],
		Stage:        csvRow[8],
		Year:         csvRow[9],
	}
}

func (m Match) toCsvRow() []string {
	return []string{
		m.SortId,
		m.Opponent,
		m.OpponentLink,
		m.winLoss,
		m.Method,
		m.MethodLink,
		m.Competition,
		m.Weight,
		m.Stage,
		m.Year,
	}
}
