package common

import "fmt"

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

// MatchupsV1d2
// The matchup values - we only need to express half.
// Later we will make this a dynamic input but hardcoding for now.
var MatchupsV1d2 = map[Matchup]float64{
	Matchup{GC, GC}: .5,
	Matchup{GC, KH}: .5,
	Matchup{GC, KI}: .5,
	Matchup{GC, NG}: .5,
	Matchup{GC, OK}: .5,
	Matchup{GC, SL}: .5,
	Matchup{GC, TZ}: .5,

	Matchup{KH, KH}: .5,
	Matchup{KH, KI}: .5,
	Matchup{KH, NG}: .5,
	Matchup{KH, OK}: .5,
	Matchup{KH, SL}: .5,
	Matchup{KH, TZ}: .5,

	Matchup{KI, KI}: .5,
	Matchup{KI, NG}: .5,
	Matchup{KI, OK}: .5,
	Matchup{KI, SL}: .5,
	Matchup{KI, TZ}: .5,

	Matchup{NG, NG}: .5,
	Matchup{NG, OK}: .5,
	Matchup{NG, SL}: .5,
	Matchup{NG, TZ}: .5,

	Matchup{OK, OK}: .5,
	Matchup{OK, SL}: .5,
	Matchup{OK, TZ}: .5,

	Matchup{SL, SL}: .5,
	Matchup{SL, TZ}: .5,

	Matchup{TZ, TZ}: .5,
}

func GetMatchupValue(matchup Matchup, info TournamentInfo) float64 {
	if val, ok := info.MatchupOdds[matchup]; ok {
		return val
	} else {
		// Search for the opposite.
		oppositeKey := Matchup{P1Pick: matchup.P2Pick, P2Pick: matchup.P1Pick}
		if val2, ok2 := info.MatchupOdds[oppositeKey]; ok2 {
			return val2
		} else {
			panic(fmt.Sprintf("Could not find matchup results for: %s v. %s ", matchup.P1Pick, matchup.P2Pick))
		}
	}
}
