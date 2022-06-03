package algo

import (
	"fmt"
	. "github.com/tmwilder/wh3-draftbot/internal/common"
	"math"
	"testing"
)

const epsilon = .00000001

func TestGetSuccessorsP3(t *testing.T) {
	tournamentInfo := TournamentInfo{RoundCount: 3, MatchupOdds: MatchupsV1d2}
	gameState := GameState{
		RoundNumber: 3,
		P2Rounds: []P2Round{
			{Picks: []Faction{KH, KI}, Matchup: Matchup{P1: KI, P2: KI}},
			{Picks: []Faction{KH, TZ}, Matchup: Matchup{P1: TZ, P2: TZ}},
		},
		P3Round: P3Round{},
	}
	gameStates := getSuccessors(tournamentInfo, gameState)

	fiveChooseThree := 5 * 4 / 2
	expected1 := fiveChooseThree * 5 // 5 for ban choices

	if !(len(gameStates) == expected1) {
		t.Errorf("Expected game states to %d items after initial pick but it had %d", expected1, len(gameStates))
	}

	gameStatesCounterPick := getSuccessors(tournamentInfo, gameStates[0])
	expected2 := (5 - 1) * 3 // 4 for opponent choices after the ban, 3 choices for the counterban

	if !(len(gameStatesCounterPick) == expected2) {
		t.Errorf("Expected game states to %d items after counterban but it had %d", expected2, len(gameStatesCounterPick))
	}

	gameStatesFinalPick := getSuccessors(tournamentInfo, gameStatesCounterPick[0])
	expected3 := 2

	if !(len(gameStatesFinalPick) == expected3) {
		t.Errorf("Expected game states to %d items after final pick but it had %d", expected3, len(gameStatesFinalPick))
	}
}

func TestGetSuccessorsP2Pregame(t *testing.T) {
	tournamentInfo := TournamentInfo{RoundCount: 3, MatchupOdds: MatchupsV1d2}
	gameState := GameState{
		RoundNumber: 1,
		P2Rounds:    []P2Round{},
		P3Round: P3Round{
			Picks:      []Faction{},
			Ban:        EMPTY,
			CounterBan: EMPTY,
			Matchup:    Matchup{P1: EMPTY, P2: EMPTY},
		},
	}
	gameStates := getSuccessors(tournamentInfo, gameState)
	expected1 := 7 * 6 / 2

	if !(len(gameStates) == expected1) {
		t.Errorf("Expected game states to %d items but it had %d", expected1, len(gameStates))
	}

	gameStatesCounterPick := getSuccessors(tournamentInfo, gameStates[0])
	expected2 := 7

	if !(len(gameStatesCounterPick) == expected2) {
		t.Errorf("Expected game states to %d items but it had %d", expected2, len(gameStatesCounterPick))
	}

	gameStatesFinalPick := getSuccessors(tournamentInfo, gameStatesCounterPick[0])
	expected3 := 2

	if !(len(gameStatesFinalPick) == expected3) {
		t.Errorf("Expected game states to %d items but it had %d", expected3, len(gameStatesFinalPick))
	}
}

func TestPick3Combo5thPick(t *testing.T) {
	gameState := GameState{
		RoundNumber: 3,
		P2Rounds: []P2Round{
			{Picks: []Faction{KH, KI}, Matchup: Matchup{P1: KI, P2: KI}},
			{Picks: []Faction{KH, TZ}, Matchup: Matchup{P1: TZ, P2: TZ}},
			{Picks: []Faction{KH, NG}, Matchup: Matchup{P1: NG, P2: NG}},
			{Picks: []Faction{KH, SL}, Matchup: Matchup{P1: SL, P2: SL}}},
		P3Round: P3Round{},
	}

	p1Combos := getThreeCombos(gameState, true)
	p2Combos := getThreeCombos(gameState, false)

	var expected = 1

	if !(len(p1Combos) == expected) {
		t.Errorf("Expected p1 combo list to have %d items but it had %d", expected, len(p1Combos))
	}

	if !(len(p1Combos) == expected) {
		t.Errorf("Expected p2 combo list to have %d items but it had %d", expected, len(p2Combos))
	}
}

func TestPick3Combo3rdPick(t *testing.T) {
	gameState := GameState{
		RoundNumber: 3,
		P2Rounds: []P2Round{
			{Picks: []Faction{KH, SL}, Matchup: Matchup{P1: KH, P2: SL}},
			{Picks: []Faction{KH, TZ}, Matchup: Matchup{P1: OK, P2: TZ}}},
		P3Round: P3Round{},
	}

	p1Combos := getThreeCombos(gameState, true)
	p2Combos := getThreeCombos(gameState, false)

	var expected = 10

	if !(len(p1Combos) == expected) {
		t.Errorf("Expected p1 combo list to have %d items but it had %d", expected, len(p1Combos))
	}

	if !(len(p1Combos) == expected) {
		t.Errorf("Expected p2 combo list to have %d items but it had %d", expected, len(p2Combos))
	}
}

func TestPick2Combo1stPick(t *testing.T) {
	gameState := GameState{
		RoundNumber: 3,
		P2Rounds:    []P2Round{},
		P3Round:     P3Round{},
	}

	p1Combos := getTwoCombos(gameState, true)
	p2Combos := getTwoCombos(gameState, false)

	var expected = 21

	if !(len(p1Combos) == expected) {
		t.Errorf("Expected p1 combo list to have %d items but it had %d", expected, len(p1Combos))
	}

	if !(len(p1Combos) == expected) {
		t.Errorf("Expected p2 combo list to have %d items but it had %d", expected, len(p2Combos))
	}
}

func TestPick2Combo3rdPick(t *testing.T) {
	gameState := GameState{
		RoundNumber: 3,
		P2Rounds: []P2Round{
			{Picks: []Faction{KH, SL}, Matchup: Matchup{P1: KH, P2: SL}},
			{Picks: []Faction{KH, TZ}, Matchup: Matchup{P1: OK, P2: TZ}}},
		P3Round: P3Round{},
	}

	p1Combos := getTwoCombos(gameState, true)
	p2Combos := getTwoCombos(gameState, false)

	var expected = 10

	if !(len(p1Combos) == expected) {
		t.Errorf("Expected p1 combo list to have %d items but it had %d", expected, len(p1Combos))
	}

	if !(len(p1Combos) == expected) {
		t.Errorf("Expected p2 combo list to have %d items but it had %d", expected, len(p2Combos))
	}
}

func TestPick2Combo4thPick(t *testing.T) {
	gameState := GameState{
		RoundNumber: 5,
		P2Rounds: []P2Round{
			{Picks: []Faction{KH, SL}, Matchup: Matchup{P1: KH, P2: SL}},
			{Picks: []Faction{TZ, OK}, Matchup: Matchup{P1: OK, P2: TZ}},
			{Picks: []Faction{NG, GC}, Matchup: Matchup{P1: NG, P2: NG}}},
		P3Round: P3Round{},
	}

	p1Combos := getTwoCombos(gameState, true)
	p2Combos := getTwoCombos(gameState, false)

	var expected = 6

	if !(len(p1Combos) == expected) {
		t.Errorf("Expected p1 combo list to have %d items but it had %d", expected, len(p1Combos))
	}

	if !(len(p1Combos) == expected) {
		t.Errorf("Expected p2 combo list to have %d items but it had %d", expected, len(p2Combos))
	}
}

func TestComputeWinRateSimple(t *testing.T) {
	tournamentInfo := TournamentInfo{
		RoundCount:  3,
		MatchupOdds: map[Matchup]float64{Matchup{P1: KH, P2: SL}: .5}}

	gameState := GameState{
		RoundNumber: 3,
		P2Rounds: []P2Round{
			{Picks: []Faction{KH, SL}, Matchup: Matchup{P1: KH, P2: SL}},
			{Picks: []Faction{KH, SL}, Matchup: Matchup{P1: KH, P2: SL}}},
		P3Round: P3Round{
			Picks:      []Faction{NG, SL, KH},
			Ban:        OK,
			CounterBan: OK,
			Matchup:    Matchup{P1: KH, P2: SL}},
	}

	winRate := computeWinRate(tournamentInfo, gameState)

	expected := .5
	if !(math.Abs(winRate-expected) < epsilon) {
		t.Errorf("Expected WR to be %f but it was %f", expected, winRate)
	}
}

func TestComputeWinRateSimple2(t *testing.T) {
	tournamentInfo := TournamentInfo{
		RoundCount:  3,
		MatchupOdds: map[Matchup]float64{Matchup{P1: KH, P2: SL}: 1.0}}

	gameState := GameState{
		RoundNumber: 3,
		P2Rounds: []P2Round{
			{Picks: []Faction{KH, SL}, Matchup: Matchup{P1: KH, P2: SL}},
			{Picks: []Faction{KH, SL}, Matchup: Matchup{P1: KH, P2: SL}}},
		P3Round: P3Round{
			Picks:      []Faction{NG, SL, KH},
			Ban:        OK,
			CounterBan: OK,
			Matchup:    Matchup{P1: KH, P2: SL}},
	}

	winRate := computeWinRate(tournamentInfo, gameState)

	expected := 1.0
	if !(math.Abs(winRate-expected) < epsilon) {
		t.Errorf("Expected WR to be %f but it was %f", expected, winRate)
	}
}

func TestComputeWinRateLessSimple(t *testing.T) {
	tournamentInfo := TournamentInfo{
		RoundCount: 3,
		MatchupOdds: map[Matchup]float64{
			Matchup{P1: KH, P2: SL}: .6,
			Matchup{P1: SL, P2: KH}: .4,
			Matchup{P1: OK, P2: SL}: .5}}

	gameState := GameState{
		RoundNumber: 3,
		P2Rounds: []P2Round{
			{Picks: []Faction{KH, SL}, Matchup: Matchup{P1: KH, P2: SL}},
			{Picks: []Faction{SL, KH}, Matchup: Matchup{P1: SL, P2: KH}}},
		P3Round: P3Round{
			Picks:      []Faction{NG, SL, KH},
			Ban:        OK,
			CounterBan: OK,
			Matchup:    Matchup{P1: OK, P2: SL}},
	}

	winRate := computeWinRate(tournamentInfo, gameState)

	expected := .5
	if !(math.Abs(winRate-expected) < epsilon) {
		t.Errorf("Expected WR to be %f but it was %f", expected, winRate)
	}
}

func TestComputeWinRateRegression(t *testing.T) {
	tournamentInfo := TournamentInfo{
		RoundCount: 5,
		MatchupOdds: map[Matchup]float64{
			Matchup{P1: NG, P2: NG}: .5,
			Matchup{P1: TZ, P2: TZ}: .5,
			Matchup{P1: KH, P2: KH}: .5,
			Matchup{P1: SL, P2: SL}: .5,
			Matchup{P1: KI, P2: KI}: .5,
		}}

	gameState := GameState{
		RoundNumber: 5,
		P2Rounds: []P2Round{
			{Picks: []Faction{KH, SL}, Matchup: Matchup{P1: NG, P2: NG}},
			{Picks: []Faction{SL, KH}, Matchup: Matchup{P1: TZ, P2: TZ}},
			{Picks: []Faction{SL, KH}, Matchup: Matchup{P1: KH, P2: KH}},
			{Picks: []Faction{SL, KH}, Matchup: Matchup{P1: SL, P2: SL}},
		},
		P3Round: P3Round{
			Picks:      []Faction{NG, SL, KH},
			Ban:        OK,
			CounterBan: OK,
			Matchup:    Matchup{P1: KI, P2: KI}},
	}

	winRate := computeWinRate(tournamentInfo, gameState)

	expected := .5
	if !(math.Abs(winRate-expected) < epsilon) {
		t.Errorf("Expected WR to be %f but it was %f", expected, winRate)
	}
}

func TestDeepCopy(t *testing.T) {
	gameState := GameState{
		RoundNumber: 3,
		P2Rounds: []P2Round{
			{Picks: []Faction{KH, SL}, Matchup: Matchup{P1: KH, P2: SL}},
			{Picks: []Faction{SL, KH}, Matchup: Matchup{P1: SL, P2: KH}}},
		P3Round: P3Round{
			Picks:      []Faction{NG, SL, KH},
			Ban:        OK,
			CounterBan: OK,
			Matchup:    Matchup{P1: OK, P2: SL}},
	}
	gameStateCopy := deepcopy(gameState)

	if !(gameStateCopy.P2Rounds[0].Matchup == gameState.P2Rounds[0].Matchup) {
		t.Errorf(fmt.Sprint("Cloned matchups not equal, failing."))
	}
	gameStateCopy.P2Rounds[0].Matchup.P1 = OK
	if gameStateCopy.P2Rounds[0].Matchup == gameState.P2Rounds[0].Matchup {
		t.Errorf(fmt.Sprint("Setting copy values impacted the non-copy, failing."))
	}

	gameState.P2Rounds[0].Matchup.P2 = KH
	if gameStateCopy.P2Rounds[0].Matchup.P2 == KH {
		t.Errorf(fmt.Sprint("Setting uncopied values impacted the original, failing."))
	}
}

/**
We were computing the wrong odds. Add that scenario.
*/
func TestComputeWinrateRegression(t *testing.T) {
	gameState := GameState{
		RoundNumber: 3,
		P2Rounds: []P2Round{
			{Picks: []Faction{}, Matchup: Matchup{P1: TZ, P2: GC}},
			{Picks: []Faction{}, Matchup: Matchup{P1: KH, P2: KH}}},
		P3Round: P3Round{
			Picks:      []Faction{},
			Ban:        OK,
			CounterBan: OK,
			Matchup:    Matchup{P1: SL, P2: NG}},
	}

	result := computeWinRate(TournamentInfo{RoundCount: 3, MatchupOdds: matchupsPolarized}, gameState)
	expected := 0.0
	if !(math.Abs(result-expected) < epsilon) {
		t.Errorf("Expected WR to be %f but it was %f", expected, result)
	}
}
