package algo

import (
	. "github.com/tmwilder/wh3-draftbot/internal/common"
	"math"
)

func TurinMinimax(tournamentInfo TournamentInfo, gameState gameState, isMaximizingPlayer bool, alpha float64, beta float64) float64 {
	draftIsComplete := (gameState.depth == tournamentInfo.RoundCount) &&
		(gameState.p3Round.matchup.P1Pick != EMPTY)
	if draftIsComplete {
		return computeWinRate(tournamentInfo, gameState)
	}

	if isMaximizingPlayer {
		bestVal := -1.0
		successors := getSuccessors(tournamentInfo, gameState)
		for _, v := range successors {
			// Make sure the false alternation works here
			value := TurinMinimax(tournamentInfo, v, false, alpha, beta)
			bestVal = math.Max(bestVal, value)
			alpha = math.Max(alpha, bestVal)
			if beta <= alpha {
				break
			}
		}
		return bestVal
	} else {
		bestVal := 2.0
		successors := getSuccessors(tournamentInfo, gameState)
		for _, v := range successors {
			// Make sure the false alternation works here
			value := TurinMinimax(tournamentInfo, v, true, alpha, beta)
			bestVal = math.Min(bestVal, value)
			beta = math.Min(beta, bestVal)
			if beta <= alpha {
				break
			}
		}
		return bestVal
	}
}
