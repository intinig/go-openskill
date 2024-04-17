package models

import (
	"math"

	"github.com/intinig/go-openskill/ptr"

	"github.com/intinig/go-openskill/types"
	"github.com/intinig/go-openskill/util"
)

type PlackettLuceOptions struct {
	Epsilon *float64
	Beta    *float64
	Mu      *float64
	Sigma   *float64
	Z       *int
}

type PlackettLuce struct {
	Epsilon        float64
	Beta           float64
	BetaSquared    float64
	TwoBetaSquared float64
	Mu             float64
	Sigma          float64
	Z              int
	U              *util.Util
}

// NewPlackettLuce returns a new PlackettLuce model
func NewPlackettLuce(options *types.OpenSkillOptions) *PlackettLuce {
	if options == nil {
		options = &types.OpenSkillOptions{}
	}

	epsilon := options.Epsilon
	if epsilon == nil {
		epsilon = ptr.Float64(0.0001)
	}

	z := options.Z
	if z == nil {
		z = ptr.Int(3)
	}

	mu := options.Mu
	if mu == nil {
		mu = ptr.Float64(25.0)
	}

	sigma := options.Sigma
	if sigma == nil {
		sigma = ptr.Float64(*mu / float64(*z))
	}

	beta := options.Beta
	if beta == nil {
		beta = ptr.Float64(*sigma / 2.0)
	}

	bSquared := *beta * *beta
	u := util.NewWithOptions(&util.Options{
		BetaSquared: ptr.Float64(bSquared),
	})

	return &PlackettLuce{
		Epsilon:        *epsilon,
		Z:              *z,
		Mu:             *mu,
		Sigma:          *sigma,
		Beta:           *beta,
		BetaSquared:    bSquared,
		TwoBetaSquared: 2 * (*beta * *beta),
		U:              u,
	}
}

// Rate rates a set of teams
func (p *PlackettLuce) Rate(teams []types.Team, options *types.OpenSkillOptions) []types.Team {
	// Initialize options
	if options == nil {
		options = &types.OpenSkillOptions{}
	}

	// We will return this at the end, it contains the same teams, with same
	// players, but with new ratings since the only thing that we really have here
	// is ratings, we will just return new ratings and what is important is the
	// order.
	returning := make([]types.Team, len(teams))

	// Create a teamRatings struct for each team
	teamRatings := p.U.TeamRating(teams, options)

	// Draw probability (?)
	c := p.U.C(teamRatings)

	// sumQ
	sumQ := p.U.SumQ(teamRatings, c)

	// Draws per team
	a := p.U.A(teamRatings)

	// Main loop, we iterate across all teamRatings
	for i, teamRating := range teamRatings {

		iMuOverCe := math.Exp(teamRating.TeamMu / c)

		var lowerRanks []types.TeamRating
		for _, teamRating2 := range teamRatings {
			if teamRating2.Rank <= teamRating.Rank {
				lowerRanks = append(lowerRanks, teamRating2)
			}
		}

		omega, delta := 0.0, 0.0
		for q := range lowerRanks {
			quotient := iMuOverCe / sumQ[q]
			if i == q {
				omega += (1.0 - quotient) / float64(a[q])
			} else {
				omega -= quotient / float64(a[q])
			}
			delta += (quotient * (1 - quotient)) / float64(a[q])
		}

		gamma := options.Gamma
		if gamma == nil {
			gamma = func(tr types.TeamRating) float64 {
				return math.Sqrt(tr.TeamSigmaSquared) / c
			}
		}

		iGamma := gamma(teamRating)
		iOmega := omega * (teamRating.TeamSigmaSquared / c)
		iDelta := iGamma * delta * (teamRating.TeamSigmaSquared / (c * c))

		returningTeam := make(types.Team, len(teamRating.Team))

		for j, rating := range teamRating.Team {
			returningTeam[j] = types.Rating{
				Mu: rating.Mu + (rating.Sigma*rating.Sigma/teamRating.TeamSigmaSquared)*iOmega,
				Sigma: rating.Sigma * math.Sqrt(
					math.Max(
						1-((rating.Sigma*rating.Sigma)/teamRating.TeamSigmaSquared)*iDelta,
						p.Epsilon,
					),
				),
				Z: rating.Z,
			}
		}

		returning[i] = returningTeam
	}

	return returning
}
