package app

import (
	"fmt"
	"github.com/gin-gonic/gin"
	. "github.com/tmwilder/wh3-draftbot/internal/algo"
	. "github.com/tmwilder/wh3-draftbot/internal/common"
	"net/http"
	"net/url"
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

func viewHandler(c *gin.Context) {
	// parseInputs(c)
	tournamentInfo := TournamentInfo{RoundCount: 3, MatchupOdds: MatchupsV1d2}
	pageData := pageData{
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
	}
	c.HTML(http.StatusOK, "draftbot.html", pageData)
}

func recommendHandler(c *gin.Context) {
	// Parse later
	tournamentInfo := TournamentInfo{RoundCount: 3, MatchupOdds: MatchupsV1d2}
	gameState := GameState{
		RoundNumber: 2,
		P2rounds:    []P2Round{{Picks: []Faction{NG, TZ}, Matchup: Matchup{P1: NG, P2: KI}}},
		P3Round:     P3Round{},
	}
	winRate, gameState := TurinMinimax(tournamentInfo, gameState, true, -1.0, 2.0)

	c.HTML(http.StatusOK, "draftbot.html", pageData{WinRate: winRate, GameState: gameState, TournamentInfo: tournamentInfo})
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
