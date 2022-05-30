package algo

import (
	. "github.com/tmwilder/wh3-draftbot/internal/common"
	"math"
	"testing"
)

const epsilon = .00000001

func TestComputeWinRateSimple(t *testing.T) {
	tournamentInfo := TournamentInfo{
		RoundCount:  3,
		MatchupOdds: map[Matchup]float64{Matchup{P1Pick: KH, P2Pick: SL}: .5}}

	gameState := gameState{
		depth:        3,
		currentRound: 3,
		p2rounds: []p2Round{
			{prospectivePicks: [2]Faction{KH, SL}, matchup: &Matchup{P1Pick: KH, P2Pick: SL}},
			{prospectivePicks: [2]Faction{KH, SL}, matchup: &Matchup{P1Pick: KH, P2Pick: SL}}},
		p3Round: &p3Round{
			prospectivePicks: [3]Faction{NG, SL, KH},
			ban:              OK,
			counterBan:       OK,
			matchup:          &Matchup{P1Pick: KH, P2Pick: SL}},
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
			{prospectivePicks: [2]Faction{KH, SL}, matchup: &Matchup{P1Pick: KH, P2Pick: SL}},
			{prospectivePicks: [2]Faction{KH, SL}, matchup: &Matchup{P1Pick: KH, P2Pick: SL}}},
		p3Round: &p3Round{
			prospectivePicks: [3]Faction{NG, SL, KH},
			ban:              OK,
			counterBan:       OK,
			matchup:          &Matchup{P1Pick: KH, P2Pick: SL}},
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
			{prospectivePicks: [2]Faction{KH, SL}, matchup: &Matchup{P1Pick: KH, P2Pick: SL}},
			{prospectivePicks: [2]Faction{SL, KH}, matchup: &Matchup{P1Pick: SL, P2Pick: KH}}},
		p3Round: &p3Round{
			prospectivePicks: [3]Faction{NG, SL, KH},
			ban:              OK,
			counterBan:       OK,
			matchup:          &Matchup{P1Pick: OK, P2Pick: SL}},
	}

	winRate := computeWinRate(tournamentInfo, gameState)

	expected := .5
	if !(math.Abs(winRate-expected) < epsilon) {
		t.Errorf("Expected WR to be %f but it was %f", expected, winRate)
	}
}
