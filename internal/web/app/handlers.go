package app

import (
	"fmt"
	. "github.com/tmwilder/wh3-draftbot/internal/algo"
	. "github.com/tmwilder/wh3-draftbot/internal/common"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

type pageData struct {
	WinRate          float64
	GameState        GameState
	TournamentInfo   TournamentInfo
	SuggestedLine    GameState
	RenderNewP2      bool
	RenderNewP3      bool
	RenderExistingP3 bool
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	t, err := loadTemplate("draftbot.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	parseInputs(r)

	tournamentInfo := TournamentInfo{RoundCount: 3, MatchupOdds: MatchupsV1d2}
	err = t.Execute(w,
		pageData{
			TournamentInfo: tournamentInfo,
			WinRate:        0.5,
			GameState: GameState{
				RoundNumber: 3,
				P2rounds: []P2Round{
					{Picks: []Faction{NG, TZ}, Matchup: Matchup{P1: NG, P2: KI}},
					{Picks: []Faction{OK, KI}, Matchup: Matchup{P1: KH, P2: KI}}},
				P3Round: P3Round{
					Picks:      []Faction{NG, SL, KH},
					Ban:        OK,
					CounterBan: OK,
					Matchup:    Matchup{P1: KH, P2: SL}}},
		},
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func recommendHandler(w http.ResponseWriter, r *http.Request) {
	// Parse later
	tournamentInfo := TournamentInfo{RoundCount: 3, MatchupOdds: MatchupsV1d2}
	gameState := GameState{
		RoundNumber: 2,
		P2rounds:    []P2Round{{Picks: []Faction{NG, TZ}, Matchup: Matchup{P1: NG, P2: KI}}},
		P3Round:     P3Round{},
	}
	winRate, gameState := TurinMinimax(tournamentInfo, gameState, true, -1.0, 2.0)

	t, err := loadTemplate("draftbot.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = t.Execute(w, pageData{WinRate: winRate, GameState: gameState, TournamentInfo: tournamentInfo})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func parseInputs(r *http.Request) (TournamentInfo, GameState) {
	queryParams, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		panic("Cannot parse input: " + err.Error())
	}
	// Extract matchup odds
	// We'll do it the gross way so we can remember life without tools ;)
	matchupOdds := map[Matchup]float64{}
	for k, v := range queryParams {
		if len(k) == 5 && k[2] == '-' {
			split := strings.Split(k, "-")
			odds, err := strconv.ParseFloat(v[0], 64)
			if err != nil {
				panic("Cannot parse input: " + err.Error())
			}
			matchupOdds[Matchup{P1: Faction(split[0]), P2: Faction(split[1])}] = odds
		}
	}

	roundsStr := queryParams.Get("rounds")
	roundCount, err := strconv.ParseInt(roundsStr, 0, 64)
	if err != nil {
		fmt.Println("Cannot parse input: " + err.Error())
		roundCount = 3
	}
	tournamentInfo := TournamentInfo{RoundCount: int(roundCount), MatchupOdds: matchupOdds}
	// Add gameInfo parse - maybe go to web framework/Gin here ;)
	return tournamentInfo, GameState{}
}

func loadTemplate(templateName string) (*template.Template, error) {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	return template.ParseFiles(wd + "/internal/web/template/" + templateName)
}
