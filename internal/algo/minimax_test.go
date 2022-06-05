package algo

import (
	"fmt"
	. "github.com/tmwilder/wh3-draftbot/internal/common"
	"math"
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
		P2Rounds: []P2Round{},
		P3Round:  P3Round{},
	}

	value, gameState := TurinMinimax(tournamentInfo, gameState, true, -1.0, 2.0)

	fmt.Printf("%f\n", value)
	fmt.Printf("%+v\n", gameState)
}

func TestMinimaxFullGame52ndPick(t *testing.T) {
	tournamentInfo := TournamentInfo{RoundCount: 5, MatchupOdds: MatchupsV1d2}
	gameState := GameState{
		P2Rounds: []P2Round{
			{Picks: []Faction{NG, TZ}, Matchup: Matchup{P1: NG, P2: KI}},
		},
		P3Round: P3Round{},
	}
	value, gameState := TurinMinimax(tournamentInfo, gameState, true, -1.0, 2.0)

	fmt.Printf("%f\n", value)
	fmt.Printf("%+v\n", gameState)
}

func TestMinimaxFullGame3rdPick(t *testing.T) {
	tournamentInfo := TournamentInfo{RoundCount: 5, MatchupOdds: MatchupsV1d2}
	gameState := GameState{
		P2Rounds: []P2Round{
			{Picks: []Faction{NG, TZ}, Matchup: Matchup{P1: NG, P2: KI}},
			{Picks: []Faction{OK, KI}, Matchup: Matchup{P1: KH, P2: KI}},
		},
		P3Round: P3Round{},
	}
	value, gameState := TurinMinimax(tournamentInfo, gameState, true, -1.0, 2.0)

	fmt.Printf("%f\n", value)
	fmt.Printf("%+v\n", gameState)
}

func TestMinimaxFullGame3(t *testing.T) {
	tournamentInfo := TournamentInfo{RoundCount: 3, MatchupOdds: MatchupsV1d2}
	gameState := GameState{
		P2Rounds: []P2Round{},
		P3Round:  P3Round{},
	}
	value, gameState := TurinMinimax(tournamentInfo, gameState, true, -1.0, 2.0)

	fmt.Printf("%f\n", value)
	fmt.Printf("%+v\n", gameState)
}

func TestMinimaxR3(t *testing.T) {
	tournamentInfo := TournamentInfo{RoundCount: 3, MatchupOdds: MatchupsV1d2}
	gameState := GameState{
		P2Rounds: []P2Round{
			{Picks: []Faction{SL, TZ}, Matchup: Matchup{P1: TZ, P2: GC}},
			{Picks: []Faction{KH, TZ}, Matchup: Matchup{P1: OK, P2: KH}},
		},
		P3Round: P3Round{},
	}

	value, gameState := TurinMinimax(tournamentInfo, gameState, true, -1.0, 2.0)

	fmt.Printf("%f\n", value)
	fmt.Printf("%+v\n", gameState)
}

func TestMinimaxR3P2Wins(t *testing.T) {
	tournamentInfo := TournamentInfo{RoundCount: 3, MatchupOdds: MatchupsV1d2}
	gameState := GameState{
		P2Rounds: []P2Round{
			{Picks: []Faction{GC, TZ}, Matchup: Matchup{P1: GC, P2: GC}, WhoWon: P1},
			{Picks: []Faction{KH, TZ}, Matchup: Matchup{P1: KH, P2: KI}, WhoWon: P2},
		},
		P3Round: P3Round{},
	}

	value, gameState := TurinMinimax(tournamentInfo, gameState, false, -1.0, 2.0)

	fmt.Printf("%f\n", value)
	fmt.Printf("%+v\n", gameState)
}

func TestMinimaxR3Polarized(t *testing.T) {
	tournamentInfo := TournamentInfo{RoundCount: 3, MatchupOdds: matchupsPolarized}
	gameState := GameState{
		P2Rounds: []P2Round{},
		P3Round:  P3Round{},
	}

	winRate, gameState := TurinMinimax(tournamentInfo, gameState, true, -1.0, 2.0)

	if !(gameState.P2Rounds[0].Matchup == Matchup{P1: GC, P2: GC}) {
		t.Errorf("First matchup was not the expected GC GC")
	}

	if !(gameState.P2Rounds[1].Matchup == Matchup{P1: KH, P2: KH}) {
		t.Errorf("First matchup was not the expected KH KH")
	}

	if !(gameState.P3Round.Matchup == Matchup{P1: NG, P2: NG}) {
		t.Errorf("First matchup was not the expected NG NG")
	}

	expected := .5
	if !(math.Abs(winRate-expected) < epsilon) {
		t.Errorf("Expected WR to be %f but it was %f", expected, winRate)
	}
}
