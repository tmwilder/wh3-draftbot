package algo

import (
	. "github.com/tmwilder/wh3-draftbot/internal/common"
	"math"
	"testing"
)

const epsilon = .00000001

func TestPick3Combo5thPick(t *testing.T) {
	gameState := gameState{
		depth:        3,
		currentRound: 3,
		p2rounds: []p2Round{
			{initialPicks: []Faction{KH, KI}, matchup: Matchup{P1Pick: KI, P2Pick: KI}},
			{initialPicks: []Faction{KH, TZ}, matchup: Matchup{P1Pick: TZ, P2Pick: TZ}},
			{initialPicks: []Faction{KH, NG}, matchup: Matchup{P1Pick: NG, P2Pick: NG}},
			{initialPicks: []Faction{KH, SL}, matchup: Matchup{P1Pick: SL, P2Pick: SL}}},
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
	gameState := gameState{
		depth:        3,
		currentRound: 3,
		p2rounds: []p2Round{
			{initialPicks: []Faction{KH, SL}, matchup: Matchup{P1Pick: KH, P2Pick: SL}},
			{initialPicks: []Faction{KH, TZ}, matchup: Matchup{P1Pick: OK, P2Pick: TZ}}},
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
	gameState := gameState{
		depth:        3,
		currentRound: 1,
		p2rounds:     []p2Round{},
		p3Round:      p3Round{},
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
	gameState := gameState{
		depth:        3,
		currentRound: 3,
		p2rounds: []p2Round{
			{initialPicks: []Faction{KH, SL}, matchup: Matchup{P1Pick: KH, P2Pick: SL}},
			{initialPicks: []Faction{KH, TZ}, matchup: Matchup{P1Pick: OK, P2Pick: TZ}}},
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
	gameState := gameState{
		depth:        5,
		currentRound: 1,
		p2rounds: []p2Round{
			{initialPicks: []Faction{KH, SL}, matchup: Matchup{P1Pick: KH, P2Pick: SL}},
			{initialPicks: []Faction{TZ, OK}, matchup: Matchup{P1Pick: OK, P2Pick: TZ}},
			{initialPicks: []Faction{NG, GC}, matchup: Matchup{P1Pick: NG, P2Pick: NG}}},
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
		MatchupOdds: map[Matchup]float64{Matchup{P1Pick: KH, P2Pick: SL}: .5}}

	gameState := gameState{
		depth:        3,
		currentRound: 3,
		p2rounds: []p2Round{
			{initialPicks: []Faction{KH, SL}, matchup: Matchup{P1Pick: KH, P2Pick: SL}},
			{initialPicks: []Faction{KH, SL}, matchup: Matchup{P1Pick: KH, P2Pick: SL}}},
		p3Round: p3Round{
			initialPicks: []Faction{NG, SL, KH},
			ban:          OK,
			counterBan:   OK,
			matchup:      Matchup{P1Pick: KH, P2Pick: SL}},
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
		MatchupOdds: map[Matchup]float64{Matchup{P1Pick: KH, P2Pick: SL}: 1.0}}

	gameState := gameState{
		depth:        3,
		currentRound: 3,
		p2rounds: []p2Round{
			{initialPicks: []Faction{KH, SL}, matchup: Matchup{P1Pick: KH, P2Pick: SL}},
			{initialPicks: []Faction{KH, SL}, matchup: Matchup{P1Pick: KH, P2Pick: SL}}},
		p3Round: p3Round{
			initialPicks: []Faction{NG, SL, KH},
			ban:          OK,
			counterBan:   OK,
			matchup:      Matchup{P1Pick: KH, P2Pick: SL}},
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
			Matchup{P1Pick: KH, P2Pick: SL}: .6,
			Matchup{P1Pick: SL, P2Pick: KH}: .4,
			Matchup{P1Pick: OK, P2Pick: SL}: .5}}

	gameState := gameState{
		depth:        3,
		currentRound: 3,
		p2rounds: []p2Round{
			{initialPicks: []Faction{KH, SL}, matchup: Matchup{P1Pick: KH, P2Pick: SL}},
			{initialPicks: []Faction{SL, KH}, matchup: Matchup{P1Pick: SL, P2Pick: KH}}},
		p3Round: p3Round{
			initialPicks: []Faction{NG, SL, KH},
			ban:          OK,
			counterBan:   OK,
			matchup:      Matchup{P1Pick: OK, P2Pick: SL}},
	}

	winRate := computeWinRate(tournamentInfo, gameState)

	expected := .5
	if !(math.Abs(winRate-expected) < epsilon) {
		t.Errorf("Expected WR to be %f but it was %f", expected, winRate)
	}
}
