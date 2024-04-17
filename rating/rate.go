package rating

import (
	"math"
	"sort"

	"github.com/intinig/go-openskill/models"
	"github.com/intinig/go-openskill/types"
	"github.com/intinig/go-openskill/unwind"
)

// Rate takes an array of ratings and returns a new array of ratings based on their performance
func Rate(teams []types.Team, options *types.OpenSkillOptions) []types.Team {
	if options == nil {
		options = &types.OpenSkillOptions{}
	}

	// Defaults to Plackett-Luce model
	// TODO: implement all models
	if options.Model == nil {
		options.Model = models.NewPlackettLuce(options)
	}

	// Save for later
	orig := teams

	// if tau is provided, use additive dynamics factor to prevent sigma from dropping too low
	if options.Tau != nil {
		t2 := *options.Tau * *options.Tau
		var newTeams = make([]types.Team, len(teams))
		for i, team := range teams {
			newTeams[i] = make(types.Team, len(team))
			for j, rating := range team {
				newTeams[i][j] = types.Rating{
					Mu:    rating.Mu,
					Sigma: math.Sqrt(rating.Sigma*rating.Sigma + t2),
					Z:     rating.Z,
				}
			}
		}
		teams = newTeams
	}

	// Prime rank with default values
	rank := make([]int, len(teams))
	for i := range rank {
		rank[i] = i
	}

	if options.Rank != nil {
		// if options.Rank is provided, use it instead
		rank = options.Rank
	} else if options.Score != nil {
		// if options.Score is provided, use it to calculate rank
		for i := range options.Score {
			rank[i] = -options.Score[i]
		}
	}

	// Unwind teams and rank
	teams, tenet := unwind.Teams(teams, rank)

	// We re-sort the rank now that teams have been sorted
	sort.Ints(rank)
	options.Rank = rank
	// Now we apply the new calculations
	newRatings := options.Model.Rate(teams, options)

	// Reverse the unwinding
	teams, _ = unwind.Teams(newRatings, tenet)

	if options.Tau != nil && options.PreventSigmaIncrease {
		for i, team := range teams {
			for j, rating := range team {
				teams[i][j] = types.Rating{
					Mu:    rating.Mu,
					Sigma: math.Min(rating.Sigma, orig[i][j].Sigma),
					Z:     rating.Z,
				}
			}
		}
	}

	return teams
}
