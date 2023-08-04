package calc

import (
	"chess-profile/types"
	"math"
	"sort"
	"strings"
)

type Calculator struct {
	openingLength int
	openingCount  int
}

func (c *Calculator) Calc(user string, games []types.Game) types.Profile {
	games = getLegitGames(games)
	sortGamesByOpening(games, c.openingLength)
	whiteOpenings := getTopOpenings(games, user, true, c.openingLength, c.openingCount)
	blackOpenings := getTopOpenings(games, user, false, c.openingLength, c.openingCount)
	mateWinRatio := getRatio(games, func(g types.Game) bool {
		return g.Outcome == types.Checkmate && g.Winner != nil && *g.Winner == user
	}, func(g types.Game) bool {
		return g.Winner != nil && *g.Winner == user
	})
	resignLossRatio := getRatio(games, func(g types.Game) bool {
		return g.Outcome == types.Resignation && g.Winner != nil && *g.Winner != user
	}, func(g types.Game) bool {
		return g.Winner != nil && *g.Winner != user
	})
	durationPercentiles := getDurationPercentiles(games)
	return types.Profile{
		User:                user,
		OpeningsWhite:       whiteOpenings,
		OpeningsBlack:       blackOpenings,
		MateWinRatio:        mateWinRatio,
		ResignLossRatio:     resignLossRatio,
		DurationPercentiles: durationPercentiles,
	}
}

func NewCalculator(openingLength int, openingCount int) *Calculator {
	return &Calculator{openingLength, openingCount}
}

func getLegitGames(games []types.Game) []types.Game {
	return filter(games, func(g types.Game) bool {
		return g.Outcome != types.Unknown
	})
}

func sortGamesByOpening(games []types.Game, openingLength int) {
	getOpeningSortKey := func(game types.Game) string {
		return strings.Join(game.Moves[:(openingLength*2)], "")
	}
	sort.Slice(games, func(i, j int) bool {
		iOpenings := filter(games, func(g types.Game) bool {
			return getOpeningSortKey(games[i]) == getOpeningSortKey(g)
		})
		jOpenings := filter(games, func(g types.Game) bool {
			return getOpeningSortKey(games[j]) == getOpeningSortKey(g)
		})
		return len(iOpenings) > len(jOpenings) // descending, more common openings first
	})
}

func getTopOpenings(games []types.Game, user string, white bool, openingLength int, openingCount int) []string {
	openings := make([]string, 0, len(games))
	for _, g := range games {
		if white && g.White == user {
			openings = append(openings, strings.Join(g.Moves[:(openingLength*2)], " "))
		}
		if !white && g.Black == user {
			openings = append(openings, strings.Join(g.Moves[:(openingLength*2)], " "))
		}
	}
	counts := make(map[string]int)
	for _, opening := range openings {
		counts[opening]++
	}
	deduplicated := make([]string, 0, len(openings))
	for opening, _ := range counts {
		deduplicated = append(deduplicated, opening)
	}
	sort.Slice(deduplicated, func(i, j int) bool {
		return counts[deduplicated[j]] > counts[deduplicated[i]]
	})
	return deduplicated[:openingCount]
}

func getRatio(games []types.Game, numerator func(g types.Game) bool, denominator func(g types.Game) bool) float64 {
	num := filter(games, numerator)
	denom := filter(games, denominator)
	return float64(len(num)) / float64(len(denom))
}

func getDurationPercentiles(games []types.Game) map[int]int {
	gamesCopy := make([]types.Game, len(games))
	copy(gamesCopy, games)
	sort.Slice(gamesCopy, func(i, j int) bool {
		return len(gamesCopy[i].Moves) < len(gamesCopy[j].Moves)
	})
	getPercentile := func(pc int) int {
		return len(gamesCopy[int(math.Round(float64(pc)/100.0*float64(len(gamesCopy))))].Moves) / 2
	}
	return map[int]int{
		50: getPercentile(50),
		75: getPercentile(75),
		90: getPercentile(90),
	}
}

func filter(gs []types.Game, f func(g types.Game) bool) []types.Game {
	newGs := make([]types.Game, 0, len(gs))
	for _, g := range gs {
		if f(g) {
			newGs = append(newGs, g)
		}
	}
	return newGs
}
