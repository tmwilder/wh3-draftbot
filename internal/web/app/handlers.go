package app

import (
	"fmt"
	"github.com/gin-gonic/gin"
	. "github.com/tmwilder/wh3-draftbot/internal/algo"
	. "github.com/tmwilder/wh3-draftbot/internal/common"
	"net/http"
	"strconv"
	"strings"
)

type pageData struct {
	WinRate        float64
	GameState      GameState
	TournamentInfo TournamentInfo
	SuggestedLine  GameState
}

func viewHandler(c *gin.Context) {
	tournamentInfo, gameState := applyDefaults(parseInputs(c))
	pageData := pageData{
		TournamentInfo: tournamentInfo,
		WinRate:        0.0,
		GameState:      gameState,
	}
	c.HTML(http.StatusOK, "draftbot.html", pageData)
}

func recommendHandler(c *gin.Context) {
	// Parse later
	tournamentInfo := TournamentInfo{RoundCount: 3, MatchupOdds: MatchupsV1d2}
	gameState := GameState{
		RoundNumber: 2,
		P2Rounds:    []P2Round{{Picks: []Faction{NG, TZ}, Matchup: Matchup{P1: NG, P2: KI}}},
		P3Round:     P3Round{},
	}
	winRate, gameState := TurinMinimax(tournamentInfo, gameState, true, -1.0, 2.0)

	c.HTML(http.StatusOK, "draftbot.html", pageData{WinRate: winRate, GameState: gameState, TournamentInfo: tournamentInfo})
}

func parseInputs(c *gin.Context) (TournamentInfo, GameState) {
	queryParams := c.Request.URL.Query()
	// Extract matchup odds
	// We'll do it the gross way so we can remember life without tools ;)
	matchupOdds := map[Matchup]float64{}
	for k, v := range queryParams {
		if k[0:5] == "odds-" {
			f1 := Faction(k[5:7])
			f2 := Faction(k[7:])
			odds, err := strconv.ParseFloat(v[0], 64)
			if err != nil {
				panic("Cannot parse input: " + err.Error())
			}
			matchupOdds[Matchup{P1: f1, P2: f2}] = odds
		}
	}

	roundsStr := queryParams.Get("rounds")
	roundCount, err := strconv.ParseInt(roundsStr, 0, 64)
	if err != nil {
		fmt.Println("Cannot parse input: " + err.Error())
		roundCount = 3
	}
	tournamentInfo := TournamentInfo{RoundCount: int(roundCount), MatchupOdds: matchupOdds}

	picks := c.QueryArray("picks")
	p1picks := c.QueryArray("p1pick")
	p2picks := c.QueryArray("p2pick")

	// Populate the P2Round info
	var p2Rounds []P2Round
	for i, v := range picks {
		if v == "" {
			continue
		}
		// TODO figure out validation and user input story
		roundPicks := parsePicks(v)

		matchup := Matchup{}
		if len(p1picks) >= i {
			matchup.P1 = Faction(p1picks[i])
		}

		if len(p2picks) >= i {
			matchup.P2 = Faction(p2picks[i])
		}
		round := P2Round{Picks: roundPicks, Matchup: matchup}
		p2Rounds = append(p2Rounds, round)
	}

	// Populate P3Round
	p3Round := P3Round{}

	p3Round.Picks = parsePicks(c.Query("last-picks"))
	p3Round.Ban = Faction(c.Query("last-ban"))
	p3Round.CounterBan = Faction(c.Query("last-counter-ban"))
	p3Round.Matchup.P1 = Faction(c.Query("last-p1pick"))
	p3Round.Matchup.P2 = Faction(c.Query("last-p2pick"))

	roundNumber := len(p2Rounds)
	// TODO scrutiny here - maybe move to re-usable fn
	roundPhase := 0
	//if p2Rounds[len(p2Rounds)-1].Matchup.P2 != EMPTY && p2Rounds[len(p2Rounds)-1].Matchup.P1 != EMPTY {
	//	roundPhase = 2
	//} else if p2Rounds[len(p2Rounds)-1].Matchup.P2 != EMPTY || p2Rounds[len(p2Rounds)-1].Matchup.P1 != EMPTY {
	//	roundPhase = 1
	//}

	gameState := GameState{
		RoundNumber: roundNumber,
		P2Rounds:    p2Rounds,
		P3Round:     p3Round,
		RoundPhase:  roundPhase,
	}

	return tournamentInfo, gameState
}

func parsePicks(factionStr string) []Faction {
	factionLetters := strings.Split(factionStr, " ")
	var factions []Faction
	for _, v := range factionLetters {
		factions = append(factions, Faction(v))
	}
	if factions[0] == EMPTY {
		return []Faction{}
	} else {
		return factions
	}
}

func applyDefaults(info TournamentInfo, state GameState) (TournamentInfo, GameState) {
	// Pad the p2rounds out to full event
	p2PickDeficit := info.RoundCount - len(state.P2Rounds) - 1
	if p2PickDeficit > 0 {
		for i := 0; i < p2PickDeficit; i++ {
			state.P2Rounds = append(state.P2Rounds, P2Round{})
		}
	}
	return info, state
}
