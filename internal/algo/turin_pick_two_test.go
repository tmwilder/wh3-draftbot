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
		roundNumber: 3,
		p2rounds: []p2Round{
			{picks: []Faction{KH, KI}, matchup: Matchup{P1: KI, P2: KI}},
			{picks: []Faction{KH, TZ}, matchup: Matchup{P1: TZ, P2: TZ}},
		},
		p3Round: p3Round{},
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
		roundNumber: 1,
		p2rounds:    []p2Round{},
		p3Round: p3Round{
			picks:      []Faction{},
			ban:        EMPTY,
			counterBan: EMPTY,
			matchup:    Matchup{P1: EMPTY, P2: EMPTY},
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
		roundNumber: 3,
		p2rounds: []p2Round{
			{picks: []Faction{KH, KI}, matchup: Matchup{P1: KI, P2: KI}},
			{picks: []Faction{KH, TZ}, matchup: Matchup{P1: TZ, P2: TZ}},
			{picks: []Faction{KH, NG}, matchup: Matchup{P1: NG, P2: NG}},
			{picks: []Faction{KH, SL}, matchup: Matchup{P1: SL, P2: SL}}},
		p3Round: p3Round{},
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
		roundNumber: 3,
		p2rounds: []p2Round{
			{picks: []Faction{KH, SL}, matchup: Matchup{P1: KH, P2: SL}},
			{picks: []Faction{KH, TZ}, matchup: Matchup{P1: OK, P2: TZ}}},
		p3Round: p3Round{},
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
		roundNumber: 3,
		p2rounds:    []p2Round{},
		p3Round:     p3Round{},
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
		roundNumber: 3,
		p2rounds: []p2Round{
			{picks: []Faction{KH, SL}, matchup: Matchup{P1: KH, P2: SL}},
			{picks: []Faction{KH, TZ}, matchup: Matchup{P1: OK, P2: TZ}}},
		p3Round: p3Round{},
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
		roundNumber: 5,
		p2rounds: []p2Round{
			{picks: []Faction{KH, SL}, matchup: Matchup{P1: KH, P2: SL}},
			{picks: []Faction{TZ, OK}, matchup: Matchup{P1: OK, P2: TZ}},
			{picks: []Faction{NG, GC}, matchup: Matchup{P1: NG, P2: NG}}},
		p3Round: p3Round{},
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
		roundNumber: 3,
		p2rounds: []p2Round{
			{picks: []Faction{KH, SL}, matchup: Matchup{P1: KH, P2: SL}},
			{picks: []Faction{KH, SL}, matchup: Matchup{P1: KH, P2: SL}}},
		p3Round: p3Round{
			picks:      []Faction{NG, SL, KH},
			ban:        OK,
			counterBan: OK,
			matchup:    Matchup{P1: KH, P2: SL}},
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
		roundNumber: 3,
		p2rounds: []p2Round{
			{picks: []Faction{KH, SL}, matchup: Matchup{P1: KH, P2: SL}},
			{picks: []Faction{KH, SL}, matchup: Matchup{P1: KH, P2: SL}}},
		p3Round: p3Round{
			picks:      []Faction{NG, SL, KH},
			ban:        OK,
			counterBan: OK,
			matchup:    Matchup{P1: KH, P2: SL}},
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
		roundNumber: 3,
		p2rounds: []p2Round{
			{picks: []Faction{KH, SL}, matchup: Matchup{P1: KH, P2: SL}},
			{picks: []Faction{SL, KH}, matchup: Matchup{P1: SL, P2: KH}}},
		p3Round: p3Round{
			picks:      []Faction{NG, SL, KH},
			ban:        OK,
			counterBan: OK,
			matchup:    Matchup{P1: OK, P2: SL}},
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
		roundNumber: 5,
		p2rounds: []p2Round{
			{picks: []Faction{KH, SL}, matchup: Matchup{P1: NG, P2: NG}},
			{picks: []Faction{SL, KH}, matchup: Matchup{P1: TZ, P2: TZ}},
			{picks: []Faction{SL, KH}, matchup: Matchup{P1: KH, P2: KH}},
			{picks: []Faction{SL, KH}, matchup: Matchup{P1: SL, P2: SL}},
		},
		p3Round: p3Round{
			picks:      []Faction{NG, SL, KH},
			ban:        OK,
			counterBan: OK,
			matchup:    Matchup{P1: KI, P2: KI}},
	}

	winRate := computeWinRate(tournamentInfo, gameState)

	expected := .5
	if !(math.Abs(winRate-expected) < epsilon) {
		t.Errorf("Expected WR to be %f but it was %f", expected, winRate)
	}
}

func TestDeepCopy(t *testing.T) {
	gameState := GameState{
		roundNumber: 3,
		p2rounds: []p2Round{
			{picks: []Faction{KH, SL}, matchup: Matchup{P1: KH, P2: SL}},
			{picks: []Faction{SL, KH}, matchup: Matchup{P1: SL, P2: KH}}},
		p3Round: p3Round{
			picks:      []Faction{NG, SL, KH},
			ban:        OK,
			counterBan: OK,
			matchup:    Matchup{P1: OK, P2: SL}},
	}
	gameStateCopy := deepcopy(gameState)

	if !(gameStateCopy.p2rounds[0].matchup == gameState.p2rounds[0].matchup) {
		t.Errorf(fmt.Sprint("Cloned matchups not equal, failing."))
	}
	gameStateCopy.p2rounds[0].matchup.P1 = OK
	if gameStateCopy.p2rounds[0].matchup == gameState.p2rounds[0].matchup {
		t.Errorf(fmt.Sprint("Setting copy values impacted the non-copy, failing."))
	}

	gameState.p2rounds[0].matchup.P2 = KH
	if gameStateCopy.p2rounds[0].matchup.P2 == KH {
		t.Errorf(fmt.Sprint("Setting uncopied values impacted the original, failing."))
	}
}
