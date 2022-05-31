package algo

import (
	"fmt"
	. "github.com/tmwilder/wh3-draftbot/internal/common"
	"testing"
)

func TestMinimax(t *testing.T) {

	tournamentInfo := TournamentInfo{RoundCount: 3, MatchupOdds: MatchupsV1d2}
	gameState := GameState{
		roundNumber: 3,
		p2rounds: []p2Round{
			{picks: []Faction{KH, KI}, matchup: Matchup{P1: KI, P2: KI}},
			{picks: []Faction{KH, TZ}, matchup: Matchup{P1: TZ, P2: TZ}},
		},
		p3Round: p3Round{},
	}

	value, gameState := TurinMinimax(tournamentInfo, gameState, true, -1.0, 2.0)

	fmt.Printf("%f\n", value)
	fmt.Printf("%+v\n", gameState)
}
