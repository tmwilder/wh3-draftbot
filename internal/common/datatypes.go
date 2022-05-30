package common

type Faction string

const (
	DC Faction = "DC"
	GC         = "GC"
	KH         = "KH"
	KI         = "KI"
	NG         = "NG"
	OK         = "OK"
	SL         = "SL"
	TZ         = "TZ"
)

type Matchup struct {
	P1Pick Faction
	P2Pick Faction
}

type TournamentInfo struct {
	RoundCount  int
	MatchupOdds map[Matchup]float64
}
