package algo

import (
	. "github.com/tmwilder/wh3-draftbot/internal/common"
	"math"
)

func TurinMinimax(tournamentInfo TournamentInfo, gameState GameState, isMaximizingPlayer bool, alpha float64, beta float64) (float64, GameState) {
	if draftIsComplete(tournamentInfo, gameState) {
		return computeWinRate(tournamentInfo, gameState), gameState
	}

	if isMaximizingPlayer {
		bestVal := -1.0
		var bestGameState GameState
		successors := getSuccessors(tournamentInfo, gameState)
		for _, v := range successors {
			value, candidateGameState := TurinMinimax(tournamentInfo, v, false, alpha, beta)

			if value > bestVal {
				bestGameState = candidateGameState
			}

			bestVal = math.Max(bestVal, value)
			alpha = math.Max(alpha, bestVal)
			if beta <= alpha {
				break
			}
		}
		return bestVal, bestGameState
	} else {
		bestVal := 2.0
		var bestGameState GameState
		successors := getSuccessors(tournamentInfo, gameState)
		for _, v := range successors {
			isMaximizingPlayerNext := true
			if finalRoundIsNext(tournamentInfo, v) && whoWonTheLastRound(v) == P2 {
				// If it's the round before last and P2 just won, p2 goes again.
				// This only happens for p2 because 2nd to last round is always even.
				isMaximizingPlayerNext = false
			}
			value, candidateGameState := TurinMinimax(tournamentInfo, v, isMaximizingPlayerNext, alpha, beta)

			if value < bestVal {
				bestGameState = candidateGameState
			}

			bestVal = math.Min(bestVal, value)
			beta = math.Min(beta, bestVal)
			if beta <= alpha {
				if isMaximizingPlayerNext == false {
					// Don't alpha/beta prune if we're taking two turns in a row because it's the second to last round.
					continue
				}
				break
			}
		}
		return bestVal, bestGameState
	}
}

func draftIsComplete(tournamentInfo TournamentInfo, gameState GameState) bool {
	return ((len(gameState.P2Rounds) + 1) == tournamentInfo.RoundCount) &&
		(gameState.P3Round.Matchup.P1 != EMPTY && gameState.P3Round.Matchup.P2 != EMPTY)
}

func finalRoundIsNext(tournamentInfo TournamentInfo, gameState GameState) bool {
	return (len(gameState.P2Rounds) == tournamentInfo.RoundCount-1) && (len(gameState.P3Round.Picks) == 0)
}

func whoWonTheLastRound(gameState GameState) WhoWon {
	return gameState.P2Rounds[len(gameState.P2Rounds)-1].WhoWon
}
