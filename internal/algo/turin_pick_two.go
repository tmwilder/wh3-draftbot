package algo

import (
	"fmt"
	"github.com/tmwilder/wh3-draftbot/internal/common"
)

type p2Round struct {
	prospectivePicks [2]common.Faction
	matchup          *common.Matchup
}

type p3Round struct {
	prospectivePicks [3]common.Faction
	ban              common.Faction
	counterBan       common.Faction
	matchup          *common.Matchup
}

type gameState struct {
	depth        int
	currentRound int
	p2rounds     []p2Round
	p3Round      *p3Round
}

type resultAndOdds struct {
	result bool
	odds   float64
}

func getSuccessors(state gameState) []gameState {
	return []gameState{}
}

/**
For one specific gamestate consisting of a full set of games, compute the odds of player one winning.
*/
func computeWinRate(tournamentInfo common.TournamentInfo, gameState gameState) float64 {
	// Validate input sanity
	if gameState.p3Round == nil {
		panic(fmt.Sprintf("Expected p3Round to be defined but it was not"))
	}
	eventLength := len(gameState.p2rounds) + 1
	if eventLength != tournamentInfo.RoundCount {
		panic(fmt.Sprintf("Expected: %d rounds but got: %d rounds instead.", tournamentInfo.RoundCount, eventLength))
	}

	var matchups []common.Matchup
	for _, v := range gameState.p2rounds {
		matchups = append(matchups, *v.matchup)
	}
	matchups = append(matchups, *gameState.p3Round.matchup)

	// Expand the result tree
	var results [][]resultAndOdds
	var stack [][]resultAndOdds

	r1Odds := tournamentInfo.MatchupOdds[matchups[0]]
	stack = append(stack, []resultAndOdds{{true, r1Odds}})
	stack = append(stack, []resultAndOdds{{false, 1.0 - r1Odds}})

	for len(stack) > 0 {
		current := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		if len(current) == len(matchups) {
			results = append(results, current)
		} else {
			rNextOdds := tournamentInfo.MatchupOdds[matchups[len(current)]]
			stack = append(stack, append(current, resultAndOdds{true, rNextOdds}))
			stack = append(stack, append(current, resultAndOdds{false, 1.0 - rNextOdds}))
		}
	}

	// Compute the P1 match winrate expected value from the result tree as a straight through multiplication
	var p1WinOdds []float64

	for _, resAndOdds := range results {
		p1WinTotal := 0
		resultProbability := 1.0
		for _, v := range resAndOdds {
			if v.result {
				p1WinTotal += 1
			}
			resultProbability *= v.odds
		}
		// > half the number of matchups is a win, so append the odds of this variant occurring
		if p1WinTotal > len(matchups)/2 {
			p1WinOdds = append(p1WinOdds, resultProbability)
		}
	}

	p1WinProbability := 0.0

	for _, v := range p1WinOdds {
		p1WinProbability += v
	}

	return p1WinProbability
}
