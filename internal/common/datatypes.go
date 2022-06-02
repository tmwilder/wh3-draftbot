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
	P1 Faction
	P2 Faction
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
	Matchup{GC, KH}: .4,
	Matchup{GC, KI}: .55,
	Matchup{GC, NG}: .4,
	Matchup{GC, OK}: .4,
	Matchup{GC, SL}: .6,
	Matchup{GC, TZ}: .4,

	Matchup{KH, KH}: .5,
	Matchup{KH, KI}: .55,
	Matchup{KH, NG}: .3,
	Matchup{KH, OK}: .6,
	Matchup{KH, SL}: .65,
	Matchup{KH, TZ}: .35,

	Matchup{KI, KI}: .5,
	Matchup{KI, NG}: .4,
	Matchup{KI, OK}: .6,
	Matchup{KI, SL}: .65,
	Matchup{KI, TZ}: .6,

	Matchup{NG, NG}: .5,
	Matchup{NG, OK}: .65,
	Matchup{NG, SL}: .7,
	Matchup{NG, TZ}: .3,

	Matchup{OK, OK}: .5,
	Matchup{OK, SL}: .6,
	Matchup{OK, TZ}: .4,

	Matchup{SL, SL}: .5,
	Matchup{SL, TZ}: .3,

	Matchup{TZ, TZ}: .5,
}

func GetMatchupValue(matchup Matchup, tournamentInfo TournamentInfo) float64 {
	if val, ok := tournamentInfo.MatchupOdds[matchup]; ok {
		return val
	} else {
		// Search for the opposite.
		oppositeKey := Matchup{P1: matchup.P2, P2: matchup.P1}
		if val2, ok2 := tournamentInfo.MatchupOdds[oppositeKey]; ok2 {
			return 1.0 - val2
		} else {
			panic(fmt.Sprintf("Could not find matchup results for: %s v. %s ", matchup.P1, matchup.P2))
		}
	}
}
