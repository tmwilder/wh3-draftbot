package common

type Faction string

const (
	GC    Faction = "GC"
	KH            = "KH"
	KI            = "KI"
	NG            = "NG"
	OK            = "OK"
	SL            = "SL"
	TZ            = "TZ"
	EMPTY         = ""
)

var Factions = map[Faction]bool{
	GC: true,
	KH: true,
	KI: true,
	NG: true,
	OK: true,
	SL: true,
	TZ: true}

type Matchup struct {
	P1Pick Faction
	P2Pick Faction
}

type TournamentInfo struct {
	RoundCount  int
	MatchupOdds map[Matchup]float64
}
