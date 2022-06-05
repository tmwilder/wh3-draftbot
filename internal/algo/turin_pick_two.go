package algo

import (
	"fmt"
	. "github.com/tmwilder/wh3-draftbot/internal/common"
	"sort"
)

type P2Round struct {
	Picks   []Faction
	Matchup Matchup
	WhoWon  WhoWon
}

type P3Round struct {
	Picks      []Faction
	Ban        Faction
	CounterBan Faction
	Matchup    Matchup
}

type GameState struct {
	P2Rounds []P2Round
	P3Round  P3Round
}

type resultAndOdds struct {
	result bool
	odds   float64
}

func getSuccessors(tournamentInfo TournamentInfo, previousGameState GameState) []GameState {
	isFinalRound := len(previousGameState.P2Rounds) == tournamentInfo.RoundCount-1 &&
		previousGameState.P2Rounds[len(previousGameState.P2Rounds)-1].Matchup.P1 != EMPTY &&
		previousGameState.P2Rounds[len(previousGameState.P2Rounds)-1].Matchup.P2 != EMPTY
	if isFinalRound {
		return getSuccessorsP3(previousGameState)
	} else {
		return getSuccessorsP2(previousGameState)
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

	isP1Pick := len(previousGameState.P2Rounds)%2 == 1
	var lastRoundsPhase int
	if len(previousGameState.P2Rounds) > 0 {
		currentRound := previousGameState.P2Rounds[len(previousGameState.P2Rounds)-1]
		lastRoundsPhase = getP2RoundPhase(currentRound, isP1Pick)
	} else {
		lastRoundsPhase = -1
	}

	switch lastRoundsPhase {
	case -1, 2:
		pickCombos := getTwoCombos(previousGameState, isP1Pick)
		for _, v := range pickCombos {
			newGameState := deepcopy(previousGameState)
			newGameState.P2Rounds = append(newGameState.P2Rounds, P2Round{})
			newGameState.P2Rounds[len(newGameState.P2Rounds)-1].Picks = v
			successors = append(successors, newGameState)
		}
		return successors
	case 0:
		remainingPicks := getRemainingPicks(previousGameState, !isP1Pick)
		for _, v := range remainingPicks {
			newGameState := deepcopy(previousGameState)
			if isP1Pick {
				// If it is p1 pick, p2 is _counterpicking_
				newGameState.P2Rounds[len(newGameState.P2Rounds)-1].Matchup.P2 = v
			} else {
				// Otherwise p1 is counterpicking
				newGameState.P2Rounds[len(newGameState.P2Rounds)-1].Matchup.P1 = v
			}
			successors = append(successors, newGameState)
		}
		return successors
	case 1:
		// We're on final pick if not the above two
		currentRound := previousGameState.P2Rounds[len(previousGameState.P2Rounds)-1]
		for _, v := range currentRound.Picks {
			newGameState := deepcopy(previousGameState)
			if isP1Pick {
				newGameState.P2Rounds[len(newGameState.P2Rounds)-1].Matchup.P1 = v
			} else {
				newGameState.P2Rounds[len(newGameState.P2Rounds)-1].Matchup.P2 = v
			}
			successors = append(successors, newGameState)
		}
		// Detect if it is the last p2 round, and if so double the successors for winning/losing
		return successors
	default:
		panic(fmt.Sprintf("Illegal round phase: %d", lastRoundsPhase))
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
func getSuccessorsP3(previousGameState GameState) []GameState {
	isP1Pick := previousGameState.P2Rounds[len(previousGameState.P2Rounds)-1].WhoWon != P2

	roundPhase := getP3RoundPhase(previousGameState.P3Round, isP1Pick)

	var successors []GameState

	switch roundPhase {
	case -1:
		pickCombos := getThreeCombos(previousGameState, isP1Pick)
		for _, initialPicks := range pickCombos {
			newGameState := deepcopy(previousGameState)
			newGameState.P3Round = P3Round{
				Picks:      []Faction{},
				Ban:        EMPTY,
				CounterBan: EMPTY,
				Matchup: Matchup{
					P1: EMPTY,
					P2: EMPTY,
				}}
			newGameState.P3Round.Picks = initialPicks

			remainingBans := getRemainingPicks(previousGameState, !isP1Pick)
			for _, ban := range remainingBans {
				newGameStateWithBan := deepcopy(newGameState)
				newGameStateWithBan.P3Round.Ban = ban
				successors = append(successors, newGameStateWithBan)
			}
		}
		return successors
	case 0:
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
				successors = append(successors, newGameStateWithBan)
			}
		}
		return successors
	case 1:
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
			successors = append(successors, newGameState)
		}
		return successors
	default:
		panic(fmt.Sprintf("Illegal round phase: %d", roundPhase))
	}
}

func getRemainingPicks(previousGameState GameState, isP1 bool) []Faction {
	remainingFactions := map[Faction]bool{}
	for k, v := range Factions {
		remainingFactions[k] = v
	}
	for _, v := range previousGameState.P2Rounds {
		if isP1 {
			delete(remainingFactions, v.Matchup.P1)
		} else {
			delete(remainingFactions, v.Matchup.P2)
		}
	}
	var remainingFactionsList []Faction

	for k, _ := range remainingFactions {
		remainingFactionsList = append(remainingFactionsList, k)
	}
	sort.Slice(remainingFactionsList, func(i, j int) bool {
		return remainingFactionsList[i] < remainingFactionsList[j]
	})
	return remainingFactionsList
}

func getTwoCombos(state GameState, isP1Pick bool) [][]Faction {
	remainingFactions := getRemainingPicks(state, isP1Pick)
	var combos [][]Faction
	for i, v := range remainingFactions {
		for j := i + 1; j < len(remainingFactions); j++ {
			combos = append(combos, []Faction{v, remainingFactions[j]})
		}
	}
	return combos
}

func getThreeCombos(state GameState, isP1Pick bool) [][]Faction {
	remainingFactions := getRemainingPicks(state, isP1Pick)
	var combos [][]Faction
	for i, v := range remainingFactions {
		for j := i + 1; j < len(remainingFactions); j++ {
			for k := j + 1; k < len(remainingFactions); k++ {
				combos = append(combos, []Faction{v, remainingFactions[j], remainingFactions[k]})
			}
		}
	}
	return combos
}

func deepcopy(state GameState) GameState {
	p2Rounds := make([]P2Round, len(state.P2Rounds))
	copy(p2Rounds, state.P2Rounds)

	var p3RoundMatchup = Matchup{
		P1: state.P3Round.Matchup.P1,
		P2: state.P3Round.Matchup.P2}

	var p3RoundPicks []Faction
	copy(state.P3Round.Picks, p3RoundPicks)

	var p3RoundCopy = P3Round{
		state.P3Round.Picks,
		state.P3Round.Ban,
		state.P3Round.CounterBan,
		p3RoundMatchup}

	return GameState{
		P2Rounds: p2Rounds,
		P3Round:  p3RoundCopy,
	}
}

/**
For one specific gamestate consisting of a full set of games, compute the odds of player one winning.
*/
func computeWinRate(tournamentInfo TournamentInfo, gameState GameState) float64 {
	// Validate input sanity
	if gameState.P3Round.Matchup.P1 == EMPTY {
		panic(fmt.Sprintf("Expected p3Round to be set but it was not"))
	}
	eventLength := len(gameState.P2Rounds) + 1
	if eventLength != tournamentInfo.RoundCount {
		panic(fmt.Sprintf("Expected: %d rounds but got: %d rounds instead.", tournamentInfo.RoundCount, eventLength))
	}

	var matchups []Matchup
	for _, v := range gameState.P2Rounds {
		matchups = append(matchups, v.Matchup)
	}
	matchups = append(matchups, gameState.P3Round.Matchup)

	// Expand the result tree
	var results [][]resultAndOdds
	var stack [][]resultAndOdds

	r1Odds := GetMatchupValue(matchups[0], tournamentInfo)
	stack = append(stack, []resultAndOdds{{true, r1Odds}})
	stack = append(stack, []resultAndOdds{{false, 1.0 - r1Odds}})

	for len(stack) > 0 {
		current := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		if len(current) == len(matchups) {
			results = append(results, current)
		} else {
			rNextOdds := GetMatchupValue(matchups[len(current)], tournamentInfo)
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

/**
From game state and tournament info, interprets what round number we are on, what phase of the round we are in, and
whether it is P1s turn.

Used to infer this info from rec requests so the client does not have explicitly input it.
*/
func InterpretRoundInfo(gameState GameState, tournamentInfo TournamentInfo) (int, int, bool) {
	p2Rounds := gameState.P2Rounds
	p3Round := gameState.P3Round
	startedRound3 := len(p3Round.Picks) != 0
	var gameRound int
	var roundPhase int
	var isP1Pick bool

	if startedRound3 {
		gameRound = tournamentInfo.RoundCount
	} else {
		gameRound = len(p2Rounds)
		if gameRound == 0 {
			gameRound = 1
		} else {
			currentRound := p2Rounds[len(p2Rounds)-1]
			isP1Pick = gameRound%2 == 1
			roundPhase = getP2RoundPhase(currentRound, isP1Pick)
		}
	}
	return gameRound, roundPhase, isP1Pick
}

func getP3RoundPhase(currentRound P3Round, isP1Pick bool) int {
	var roundPhase int
	if isP1Pick {
		if currentRound.Matchup.P1 != EMPTY {
			roundPhase = 2
		} else if currentRound.Matchup.P2 != EMPTY {
			roundPhase = 1
		} else if currentRound.Ban != EMPTY {
			roundPhase = 0
		} else {
			roundPhase = -1
		}
	} else {
		if currentRound.Matchup.P2 != EMPTY {
			roundPhase = 2
		} else if currentRound.Matchup.P1 != EMPTY {
			roundPhase = 1
		} else if currentRound.Ban != EMPTY {
			roundPhase = 0
		} else {
			roundPhase = -1
		}
	}
	return roundPhase
}

func getP2RoundPhase(currentRound P2Round, isP1Pick bool) int {
	var roundPhase int
	if isP1Pick {
		if currentRound.Matchup.P1 != EMPTY {
			roundPhase = 2
		} else if currentRound.Matchup.P2 != EMPTY {
			roundPhase = 1
		} else if len(currentRound.Picks) == 2 {
			roundPhase = 0
		} else {
			roundPhase = -1
		}
	} else {
		if currentRound.Matchup.P2 != EMPTY {
			roundPhase = 2
		} else if currentRound.Matchup.P1 != EMPTY {
			roundPhase = 1
		} else if len(currentRound.Picks) == 2 {
			roundPhase = 0
		} else {
			roundPhase = -1
		}
	}
	return roundPhase
}
