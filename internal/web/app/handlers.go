package app

import (
	. "github.com/tmwilder/wh3-draftbot/internal/algo"
	. "github.com/tmwilder/wh3-draftbot/internal/common"
	"html/template"
	"log"
	"net/http"
	"os"
)

type pageData struct {
	WinRate          float64
	GameState        GameState
	TournamentInfo   TournamentInfo
	SuggestedLine    GameState
	RenderFreshP2    bool
	RenderFreshP3    bool
	RenderExistingP3 bool
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	t, err := loadTemplate("draftbot.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
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
			RenderFreshP2:    false,
			RenderFreshP3:    false,
			RenderExistingP3: true,
		},
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func evaluateHandler(w http.ResponseWriter, r *http.Request) {
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

func loadTemplate(templateName string) (*template.Template, error) {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	return template.ParseFiles(wd + "/internal/web/template/" + templateName)
}
