package algo

import (
	. "github.com/tmwilder/wh3-draftbot/internal/common"
	"math"
)

func TurinMinimax(tournamentInfo TournamentInfo, gameState GameState, isMaximizingPlayer bool, alpha float64, beta float64) (float64, GameState) {
	draftIsComplete := (gameState.roundNumber == tournamentInfo.RoundCount) &&
		(gameState.p3Round.matchup.P1 != EMPTY && gameState.p3Round.matchup.P2 != EMPTY)
	if draftIsComplete {
		return computeWinRate(tournamentInfo, gameState), gameState
	}

	if isMaximizingPlayer {
		bestVal := -1.0
		var bestGameState GameState
		successors := getSuccessors(tournamentInfo, gameState)
		for _, v := range successors {
			// Make sure the false alternation works here
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
			// Make sure the false alternation works here
			value, candidateGameState := TurinMinimax(tournamentInfo, v, true, alpha, beta)

			if value < bestVal {
				bestGameState = candidateGameState
			}

			bestVal = math.Min(bestVal, value)
			beta = math.Min(beta, bestVal)
			if beta <= alpha {
				break
			}
		}
		return bestVal, bestGameState
	}
}
