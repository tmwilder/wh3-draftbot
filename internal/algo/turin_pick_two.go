package algo

import (
	"fmt"
	"github.com/tmwilder/wh3-draftbot/internal/common"
	"sort"
)

type P2Round struct {
	Picks       []common.Faction
	Matchup     common.Matchup
	CounterPick common.Faction
	FinalPick   common.Faction
}

type P3Round struct {
	Picks      []common.Faction
	Ban        common.Faction
	CounterBan common.Faction
	Matchup    common.Matchup
}

type GameState struct {
	RoundNumber int
	P2rounds    []P2Round
	P3Round     P3Round
	RoundPhase  int
}

type resultAndOdds struct {
	result bool
	odds   float64
}

func getSuccessors(tournamentInfo common.TournamentInfo, previousGameState GameState) []GameState {
	if previousGameState.RoundNumber < tournamentInfo.RoundCount {
		return getSuccessorsP2(previousGameState)
	} else {
		// TODO swap this to a pobabalistic handler in the recursive routine + by forking here by _both_ first and second pick and and assigning odds weighting in the caller
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
func getSuccessorsP2(previousGameState GameState) []GameState {
	var successors []GameState

	isP1Pick := previousGameState.RoundNumber%2 == 1
	isInitialPicks := previousGameState.RoundPhase == 0

	if isInitialPicks {
		pickCombos := getTwoCombos(previousGameState, isP1Pick)
		for _, v := range pickCombos {
			newGameState := deepcopy(previousGameState)
			newGameState.P2rounds = append(newGameState.P2rounds, P2Round{})
			newGameState.P2rounds[previousGameState.RoundNumber-1].Picks = v
			newGameState.RoundPhase = 1
			successors = append(successors, newGameState)
		}
		return successors
	}

	currentRound := previousGameState.P2rounds[previousGameState.RoundNumber-1]

	var isCounterPick = previousGameState.RoundPhase == 1

	if isCounterPick {
		remainingPicks := getRemainingPicks(previousGameState, !isP1Pick)
		for _, v := range remainingPicks {
			newGameState := deepcopy(previousGameState)
			if isP1Pick {
				// If it is p1 pick, p2 is _counterpicking_
				newGameState.P2rounds[previousGameState.RoundNumber-1].Matchup.P2 = v
			} else {
				// Otherwise p1 is counterpicking
				newGameState.P2rounds[previousGameState.RoundNumber-1].Matchup.P1 = v
			}
			newGameState.RoundPhase = 2
			successors = append(successors, newGameState)
		}
		return successors
	} else {
		// We're on final pick if not the above two
		for _, v := range currentRound.Picks {
			newGameState := deepcopy(previousGameState)
			if isP1Pick {
				newGameState.P2rounds[previousGameState.RoundNumber-1].Matchup.P1 = v
			} else {
				newGameState.P2rounds[previousGameState.RoundNumber-1].Matchup.P2 = v
			}
			newGameState.RoundNumber += 1
			newGameState.RoundPhase = 0
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
func getSuccessorsP3(previousGameState GameState, isP1Pick bool) []GameState {
	isInitialPick := previousGameState.RoundPhase == 0

	var successors []GameState

	if isInitialPick {
		pickCombos := getThreeCombos(previousGameState, isP1Pick)
		for _, initialPicks := range pickCombos {
			newGameState := deepcopy(previousGameState)
			newGameState.P3Round = P3Round{
				Picks:      []common.Faction{},
				Ban:        common.EMPTY,
				CounterBan: common.EMPTY,
				Matchup: common.Matchup{
					P1: common.EMPTY,
					P2: common.EMPTY,
				}}
			newGameState.P3Round.Picks = initialPicks

			remainingBans := getRemainingPicks(previousGameState, !isP1Pick)
			for _, ban := range remainingBans {
				newGameStateWithBan := deepcopy(newGameState)
				newGameStateWithBan.P3Round.Ban = ban
				newGameStateWithBan.RoundPhase = 1
				successors = append(successors, newGameStateWithBan)
			}
		}
		return successors
	}
	isCounterPick := previousGameState.RoundPhase == 1
	if isCounterPick {
		counterBans := previousGameState.P3Round.Picks
		for _, counterBan := range counterBans {
			remainingPicks := getRemainingPicks(previousGameState, !isP1Pick)

			for _, pick := range remainingPicks {
				if pick == previousGameState.P3Round.Ban {
					continue
				}
				newGameStateWithBan := deepcopy(previousGameState)
				newGameStateWithBan.P3Round.CounterBan = counterBan
				if isP1Pick {
					newGameStateWithBan.P3Round.Matchup.P2 = pick
				} else {
					newGameStateWithBan.P3Round.Matchup.P1 = pick
				}
				newGameStateWithBan.RoundPhase = 2
				successors = append(successors, newGameStateWithBan)
			}
		}
		return successors
	} else {
		// We're on final pick if not the above two
		for _, v := range previousGameState.P3Round.Picks {
			if v == previousGameState.P3Round.CounterBan {
				continue
			}
			newGameState := deepcopy(previousGameState)
			if isP1Pick {
				newGameState.P3Round.Matchup.P1 = v
			} else {
				newGameState.P3Round.Matchup.P2 = v
			}
			newGameState.RoundPhase = 0
			successors = append(successors, newGameState)
		}
		return successors
	}
}

func getRemainingPicks(previousGameState GameState, isP1 bool) []common.Faction {
	remainingFactions := map[common.Faction]bool{}
	for k, v := range common.Factions {
		remainingFactions[k] = v
	}
	for _, v := range previousGameState.P2rounds {
		if isP1 {
			delete(remainingFactions, v.Matchup.P1)
		} else {
			delete(remainingFactions, v.Matchup.P2)
		}
	}
	var remainingFactionsList []common.Faction

	for k, _ := range remainingFactions {
		remainingFactionsList = append(remainingFactionsList, k)
	}
	sort.Slice(remainingFactionsList, func(i, j int) bool {
		return remainingFactionsList[i] < remainingFactionsList[j]
	})
	return remainingFactionsList
}

func getTwoCombos(state GameState, isP1Pick bool) [][]common.Faction {
	remainingFactions := getRemainingPicks(state, isP1Pick)
	var combos [][]common.Faction
	for i, v := range remainingFactions {
		for j := i + 1; j < len(remainingFactions); j++ {
			combos = append(combos, []common.Faction{v, remainingFactions[j]})
		}
	}
	return combos
}

func getThreeCombos(state GameState, isP1Pick bool) [][]common.Faction {
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

func deepcopy(state GameState) GameState {
	p2Rounds := make([]P2Round, len(state.P2rounds))
	copy(p2Rounds, state.P2rounds)

	var p3RoundMatchup = common.Matchup{
		P1: state.P3Round.Matchup.P1,
		P2: state.P3Round.Matchup.P2}

	var p3RoundPicks []common.Faction
	copy(state.P3Round.Picks, p3RoundPicks)

	var p3RoundCopy = P3Round{
		state.P3Round.Picks,
		state.P3Round.Ban,
		state.P3Round.CounterBan,
		p3RoundMatchup}

	return GameState{
		RoundNumber: state.RoundNumber,
		P2rounds:    p2Rounds,
		P3Round:     p3RoundCopy,
	}
}

/**
For one specific gamestate consisting of a full set of games, compute the odds of player one winning.
*/
func computeWinRate(tournamentInfo common.TournamentInfo, gameState GameState) float64 {
	// Validate input sanity
	if gameState.P3Round.Matchup.P1 == common.EMPTY {
		panic(fmt.Sprintf("Expected p3Round to be set but it was not"))
	}
	eventLength := len(gameState.P2rounds) + 1
	if eventLength != tournamentInfo.RoundCount {
		panic(fmt.Sprintf("Expected: %d rounds but got: %d rounds instead.", tournamentInfo.RoundCount, eventLength))
	}

	var matchups []common.Matchup
	for _, v := range gameState.P2rounds {
		matchups = append(matchups, v.Matchup)
	}
	matchups = append(matchups, gameState.P3Round.Matchup)

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
			nextWin := make([]resultAndOdds, len(current))
			nextLoss := make([]resultAndOdds, len(current))
			copy(nextWin, current)
			copy(nextLoss, current)
			nextWin = append(nextWin, resultAndOdds{true, rNextOdds})
			nextLoss = append(nextLoss, resultAndOdds{false, 1.0 - rNextOdds})
			stack = append(stack, nextWin)
			stack = append(stack, nextLoss)
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
