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
			{initialPicks: []Faction{KH, KI}, matchup: Matchup{P1Pick: KI, P2Pick: KI}},
			{initialPicks: []Faction{KH, TZ}, matchup: Matchup{P1Pick: TZ, P2Pick: TZ}},
		},
		p3Round: p3Round{},
	}

	value, gameState := TurinMinimax(tournamentInfo, gameState, true, -1.0, 2.0)

	fmt.Printf("%f\n", value)
	fmt.Printf("%+v\n", gameState)
}
