package util

import (
	"math"
	"sort"

	"github.com/intinig/go-openskill/ptr"
	"github.com/intinig/go-openskill/types"
)

// Options is a struct for the options of the util constructor
type Options struct {
	BetaSquared *float64
}

// Util is a struct for the util functions
type Util struct {
	BetaSquared float64
}

// NewWithOptions returns a new Util with custom options
func NewWithOptions(options *Options) *Util {
	betaSquared := options.BetaSquared
	if betaSquared == nil {
		betaSquared = ptr.Float64(0.0)
	}

	return &Util{
		BetaSquared: *betaSquared,
	}
}

// TeamRating aggregates a rating for all teams and returns the reduced data
// structure
func (u *Util) TeamRating(teams []types.Team, options *types.OpenSkillOptions) []types.TeamRating {
	if options == nil {
		options = &types.OpenSkillOptions{}
	}

	teamRankings := u.Rankings(teams, options.Rank)
	teamRatings := make([]types.TeamRating, len(teams))
	for i, team := range teams {
		tMu := 0.0
		tSigmaSquare := 0.0
		for _, rating := range team {
			tMu += rating.Mu
			tSigmaSquare += rating.Sigma * rating.Sigma
		}

		teamRatings[i] = types.TeamRating{
			TeamMu:           tMu,
			TeamSigmaSquared: tSigmaSquare,
			Team:             team,
			Rank:             teamRankings[i],
		}
	}
	return teamRatings
}

// C is a constant used in the calculation of the draw probability
func (u *Util) C(teamRatings []types.TeamRating) float64 {
	var sum float64
	for _, teamRating := range teamRatings {
		sum += teamRating.TeamSigmaSquared + u.BetaSquared
	}

	return math.Sqrt(sum)
}

// A returns an array with the number of draws per each team
func (u *Util) A(teamRatings []types.TeamRating) []int {
	returning := make([]int, len(teamRatings))
	for i, teamRating := range teamRatings {
		for _, teamRating2 := range teamRatings {
			if teamRating.Rank == teamRating2.Rank {
				returning[i]++
			}
		}
	}
	return returning
}

// SumQ returns the sum of the Q function (iMu/c) for each team
func (u *Util) SumQ(teamRatings []types.TeamRating, c float64) []float64 {
	returning := make([]float64, len(teamRatings))
	for i, teamRating := range teamRatings {
		var lowerRanked []types.TeamRating
		for _, teamRating2 := range teamRatings {
			if teamRating2.Rank >= teamRating.Rank {
				lowerRanked = append(lowerRanked, teamRating2)
			}
		}
		for _, tr := range lowerRanked {
			returning[i] += math.Exp(tr.TeamMu / c)
		}
	}
	return returning
}

// Rankings returns the normalized rankings of the teams
func (u *Util) Rankings(teams []types.Team, ranks []int) []int {
	if ranks == nil {
		ranks = make([]int, len(teams))
		for i := range ranks {
			ranks[i] = i
		}
		return ranks
	}

	origMap := make(map[int][]int)
	for i, rank := range ranks {
		origMap[rank] = append(origMap[rank], i)
	}

	uniques := make([]int, 0, len(origMap))
	for k := range origMap {
		uniques = append(uniques, k)
	}

	sort.Ints(uniques)

	returning := make([]int, len(teams))
	delta := 0
	for i, rank := range uniques {
		for _, index := range origMap[rank] {
			returning[index] = i + delta
		}
		// offset the next rank by the number of teams with the same rank
		delta += len(origMap[rank]) - 1
	}

	return returning
}
