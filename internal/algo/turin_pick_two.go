package algo

import (
	"fmt"
	"github.com/tmwilder/wh3-draftbot/internal/common"
)

type p2Round struct {
	initialPicks []common.Faction
	matchup      common.Matchup
}

type p3Round struct {
	initialPicks []common.Faction
	ban          common.Faction
	counterBan   common.Faction
	matchup      common.Matchup
}

type gameState struct {
	roundNumber int
	p2rounds    []p2Round
	p3Round     p3Round
}

type resultAndOdds struct {
	result bool
	odds   float64
}

func getSuccessors(tournamentInfo common.TournamentInfo, previousGameState gameState) []gameState {
	if previousGameState.roundNumber < tournamentInfo.RoundCount {
		return getSuccessorsP2(previousGameState)
	} else {
		return getSuccessorsP3(previousGameState, true)
	}
}

/**
Evaluate a pre-last pick, so all rounds but the last in Turin rules.
There are 3 cases:
1. We need to do initial picks
2. We need to counterpick
3. We need to do final pick
We also must figure out who is picking first, p1, or p2 from gamestate.
We also need to determine what factions remain for each player if we are in 1 or 3.
*/
func getSuccessorsP2(previousGameState gameState) []gameState {
	var successors []gameState

	isP1Pick := previousGameState.roundNumber%2 == 1
	isInitialPicks := len(previousGameState.p2rounds) < previousGameState.roundNumber

	if isInitialPicks {
		pickCombos := getTwoCombos(previousGameState, isP1Pick)
		for _, v := range pickCombos {
			newGameState := deepcopy(previousGameState)
			newGameState.p2rounds = append(newGameState.p2rounds, p2Round{})
			newGameState.p2rounds[previousGameState.roundNumber-1].initialPicks = v
			successors = append(successors, newGameState)
		}
		return successors
	}

	currentRound := previousGameState.p2rounds[previousGameState.roundNumber-1]

	var isCounterPick bool
	if isP1Pick {
		isCounterPick = currentRound.matchup.P2Pick == common.EMPTY
	} else {
		isCounterPick = currentRound.matchup.P1Pick == common.EMPTY
	}

	if isCounterPick {
		remainingPicks := getRemainingPicks(previousGameState, isP1Pick)
		for _, v := range remainingPicks {
			newGameState := deepcopy(previousGameState)
			if isP1Pick {
				// If it is p1 pick, p2 is _counterpicking_
				newGameState.p2rounds[previousGameState.roundNumber-1].matchup.P2Pick = v
			} else {
				// Otherwise p1 is counterpicking
				newGameState.p2rounds[previousGameState.roundNumber-1].matchup.P1Pick = v
			}
			successors = append(successors, newGameState)
		}
		return successors
	} else {
		// We're on final pick if not the above two
		for _, v := range currentRound.initialPicks {
			newGameState := deepcopy(previousGameState)
			if isP1Pick {
				newGameState.p2rounds[previousGameState.roundNumber-1].matchup.P1Pick = v
			} else {
				newGameState.p2rounds[previousGameState.roundNumber-1].matchup.P2Pick = v
			}
			newGameState.roundNumber += 1
			successors = append(successors, newGameState)
		}
		return successors
	}
}

/**
Evaluate the last pick.
There are 3 cases:
1. We need to do initial picks and ban
2. We need to do counterpick and ban
3. We need to do final pick
We also must figure out who is picking first, p1 or p2 from gamestate.
We also need to determine what factions remain for each player in 1 or 3.
*/
func getSuccessorsP3(previousGameState gameState, isP1Pick bool) []gameState {
	isInitialPick := len(previousGameState.p3Round.initialPicks) == 0

	var successors []gameState

	if isInitialPick {
		pickCombos := getThreeCombos(previousGameState, isP1Pick)
		for _, initialPicks := range pickCombos {
			newGameState := deepcopy(previousGameState)
			newGameState.p3Round = p3Round{
				initialPicks: []common.Faction{},
				ban:          common.EMPTY,
				counterBan:   common.EMPTY,
				matchup: common.Matchup{
					P1Pick: common.EMPTY,
					P2Pick: common.EMPTY,
				}}
			newGameState.p3Round.initialPicks = initialPicks

			remainingBans := getRemainingPicks(previousGameState, !isP1Pick)
			for _, ban := range remainingBans {
				newGameStateWithBan := deepcopy(newGameState)
				newGameStateWithBan.p3Round.ban = ban
				successors = append(successors, newGameStateWithBan)
			}
			return successors
		}
	}
	// TODO need to ensure we properly initialize nullity story across board for this to jive
	isCounterPick := previousGameState.p3Round.counterBan == common.EMPTY
	if isCounterPick {
		counterBans := previousGameState.p3Round.initialPicks
		for _, counterBan := range counterBans {
			newGameState := deepcopy(previousGameState)
			newGameState.p3Round.counterBan = counterBan
			remainingPicks := getRemainingPicks(previousGameState, !isP1Pick)

			for _, pick := range remainingPicks {
				newGameStateWithBan := deepcopy(newGameState)
				newGameStateWithBan.p3Round.counterBan = counterBan
				if !isP1Pick {
					newGameStateWithBan.p3Round.matchup.P2Pick = pick
				} else {
					newGameStateWithBan.p3Round.matchup.P1Pick = pick
				}
				successors = append(successors, newGameState)
			}
		}
		return successors
	} else {
		// We're on final pick if not the above two
		for _, v := range previousGameState.p3Round.initialPicks {
			newGameState := deepcopy(previousGameState)
			if isP1Pick {
				newGameState.p3Round.matchup.P1Pick = v
			} else {
				newGameState.p3Round.matchup.P2Pick = v
			}
			newGameState.roundNumber += 1
			successors = append(successors, newGameState)
		}
		return successors
	}
}

func getRemainingPicks(previousGameState gameState, isP1 bool) []common.Faction {
	remainingFactions := map[common.Faction]bool{}
	for k, v := range common.Factions {
		remainingFactions[k] = v
	}
	for _, v := range previousGameState.p2rounds {
		if isP1 {
			delete(remainingFactions, v.matchup.P1Pick)
		} else {
			delete(remainingFactions, v.matchup.P2Pick)
		}
	}
	var remainingFactionsList []common.Faction

	for k, _ := range remainingFactions {
		remainingFactionsList = append(remainingFactionsList, k)
	}
	return remainingFactionsList
}

func getTwoCombos(state gameState, isP1Pick bool) [][]common.Faction {
	remainingFactions := getRemainingPicks(state, isP1Pick)
	var combos [][]common.Faction
	for i, v := range remainingFactions {
		for j := i + 1; j < len(remainingFactions); j++ {
			combos = append(combos, []common.Faction{v, remainingFactions[j]})
		}
	}
	return combos
}

func getThreeCombos(state gameState, isP1Pick bool) [][]common.Faction {
	remainingFactions := getRemainingPicks(state, isP1Pick)
	var combos [][]common.Faction
	for i, v := range remainingFactions {
		for j := i + 1; j < len(remainingFactions); j++ {
			for k := j + 1; k < len(remainingFactions); k++ {
				combos = append(combos, []common.Faction{v, remainingFactions[j], remainingFactions[k]})
			}
		}
	}
	return combos
}

func deepcopy(state gameState) gameState {
	// TODO after walk - might have copy bugs here
	p2Rounds := make([]p2Round, len(state.p2rounds))
	copy(p2Rounds, state.p2rounds)

	var p3RoundMatchup = common.Matchup{
		P1Pick: state.p3Round.matchup.P1Pick,
		P2Pick: state.p3Round.matchup.P2Pick}

	var p3RoundPicks []common.Faction
	copy(state.p3Round.initialPicks, p3RoundPicks)

	var p3RoundCopy = p3Round{
		state.p3Round.initialPicks,
		state.p3Round.ban,
		state.p3Round.counterBan,
		p3RoundMatchup}

	return gameState{
		roundNumber: state.roundNumber,
		p2rounds:    p2Rounds,
		p3Round:     p3RoundCopy,
	}
}

/**
For one specific gamestate consisting of a full set of games, compute the odds of player one winning.
*/
func computeWinRate(tournamentInfo common.TournamentInfo, gameState gameState) float64 {
	// Validate input sanity
	if gameState.p3Round.matchup.P1Pick == common.EMPTY {
		panic(fmt.Sprintf("Expected p3Round to be set but it was not"))
	}
	eventLength := len(gameState.p2rounds) + 1
	if eventLength != tournamentInfo.RoundCount {
		panic(fmt.Sprintf("Expected: %d rounds but got: %d rounds instead.", tournamentInfo.RoundCount, eventLength))
	}

	var matchups []common.Matchup
	for _, v := range gameState.p2rounds {
		matchups = append(matchups, v.matchup)
	}
	matchups = append(matchups, gameState.p3Round.matchup)

	// Expand the result tree
	var results [][]resultAndOdds
	var stack [][]resultAndOdds

	r1Odds := common.GetMatchupValue(matchups[0], tournamentInfo)
	stack = append(stack, []resultAndOdds{{true, r1Odds}})
	stack = append(stack, []resultAndOdds{{false, 1.0 - r1Odds}})

	for len(stack) > 0 {
		current := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		if len(current) == len(matchups) {
			results = append(results, current)
		} else {
			rNextOdds := common.GetMatchupValue(matchups[len(current)], tournamentInfo)
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
