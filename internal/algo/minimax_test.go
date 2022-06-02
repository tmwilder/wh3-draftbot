package algo

import (
	"fmt"
	. "github.com/tmwilder/wh3-draftbot/internal/common"
	"testing"
)

var matchupsPolarized = map[Matchup]float64{
	Matchup{GC, GC}: .5,
	Matchup{GC, KH}: 1.0,
	Matchup{GC, KI}: 1.0,
	Matchup{GC, NG}: 1.0,
	Matchup{GC, OK}: 1.0,
	Matchup{GC, SL}: 1.0,
	Matchup{GC, TZ}: 1.0,

	Matchup{KH, KH}: 0.5,
	Matchup{KH, KI}: 1.0,
	Matchup{KH, NG}: 1.0,
	Matchup{KH, OK}: 1.0,
	Matchup{KH, SL}: 1.0,
	Matchup{KH, TZ}: 1.0,

	Matchup{KI, KI}: 0.5,
	Matchup{KI, NG}: 1.0,
	Matchup{KI, OK}: 1.0,
	Matchup{KI, SL}: 1.0,
	Matchup{KI, TZ}: 1.0,

	Matchup{NG, NG}: 0.5,
	Matchup{NG, OK}: 1.0,
	Matchup{NG, SL}: 1.0,
	Matchup{NG, TZ}: 1.0,

	Matchup{OK, OK}: 0.5,
	Matchup{OK, SL}: 0.0,
	Matchup{OK, TZ}: 0.0,

	Matchup{SL, SL}: 0.5,
	Matchup{SL, TZ}: 0.0,

	Matchup{TZ, TZ}: 0.5,
}

func TestMinimaxFullGame5(t *testing.T) {
	tournamentInfo := TournamentInfo{RoundCount: 5, MatchupOdds: MatchupsV1d2}
	gameState := GameState{
		RoundNumber: 1,
		P2rounds:    []P2Round{},
		P3Round:     P3Round{},
		RoundPhase:  0,
	}

	value, gameState := TurinMinimax(tournamentInfo, gameState, true, -1.0, 2.0)

	fmt.Printf("%f\n", value)
	fmt.Printf("%+v\n", gameState)
}

func TestMinimaxFullGame52ndPick(t *testing.T) {
	tournamentInfo := TournamentInfo{RoundCount: 5, MatchupOdds: MatchupsV1d2}
	gameState := GameState{
		RoundNumber: 2,
		P2rounds: []P2Round{
			{Picks: []Faction{NG, TZ}, Matchup: Matchup{P1: NG, P2: KI}},
		},
		P3Round:    P3Round{},
		RoundPhase: 0,
	}
	value, gameState := TurinMinimax(tournamentInfo, gameState, true, -1.0, 2.0)

	fmt.Printf("%f\n", value)
	fmt.Printf("%+v\n", gameState)
}

func TestMinimaxFullGame3rdPick(t *testing.T) {
	tournamentInfo := TournamentInfo{RoundCount: 5, MatchupOdds: MatchupsV1d2}
	gameState := GameState{
		RoundNumber: 3,
		P2rounds: []P2Round{
			{Picks: []Faction{NG, TZ}, Matchup: Matchup{P1: NG, P2: KI}},
			{Picks: []Faction{OK, KI}, Matchup: Matchup{KH, KI}},
		},
		P3Round:    P3Round{},
		RoundPhase: 0,
	}
	value, gameState := TurinMinimax(tournamentInfo, gameState, true, -1.0, 2.0)

	fmt.Printf("%f\n", value)
	fmt.Printf("%+v\n", gameState)
}

func TestMinimaxFullGame3(t *testing.T) {
	tournamentInfo := TournamentInfo{RoundCount: 3, MatchupOdds: MatchupsV1d2}
	gameState := GameState{
		RoundNumber: 1,
		P2rounds:    []P2Round{},
		P3Round:     P3Round{},
		RoundPhase:  0,
	}
	value, gameState := TurinMinimax(tournamentInfo, gameState, true, -1.0, 2.0)

	fmt.Printf("%f\n", value)
	fmt.Printf("%+v\n", gameState)
}

func TestMinimaxR3(t *testing.T) {
	tournamentInfo := TournamentInfo{RoundCount: 3, MatchupOdds: MatchupsV1d2}
	gameState := GameState{
		RoundNumber: 3,
		P2rounds: []P2Round{
			{Picks: []Faction{SL, TZ}, Matchup: Matchup{P1: TZ, P2: GC}},
			{Picks: []Faction{KH, TZ}, Matchup: Matchup{P1: OK, P2: KH}},
		},
		P3Round:    P3Round{},
		RoundPhase: 0,
	}

	value, gameState := TurinMinimax(tournamentInfo, gameState, true, -1.0, 2.0)

	fmt.Printf("%f\n", value)
	fmt.Printf("%+v\n", gameState)
}

func TestMinimaxR3Polarized(t *testing.T) {
	tournamentInfo := TournamentInfo{RoundCount: 3, MatchupOdds: matchupsPolarized}
	gameState := GameState{
		RoundNumber: 1,
		P2rounds:    []P2Round{},
		P3Round:     P3Round{},
		RoundPhase:  0,
	}

	value, gameState := TurinMinimax(tournamentInfo, gameState, true, -1.0, 2.0)

	fmt.Printf("%f\n", value)
	fmt.Printf("%+v\n", gameState)
}
