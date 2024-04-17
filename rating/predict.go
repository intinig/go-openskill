package rating

import (
	"math"

	"gonum.org/v1/gonum/stat/distuv"

	"github.com/intinig/go-openskill/types"
	"github.com/intinig/go-openskill/util"
)

// getBetas returns the beta squared value, it's a helper function
func getBetas(options *types.OpenSkillOptions) (float64, float64) {
	if options == nil {
		options = &types.OpenSkillOptions{}
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

	return *beta, *beta * *beta
}

// PredictWin returns the probability of each team winning
func PredictWin(teams []types.Team, options *types.OpenSkillOptions) []float64 {
	// Initialize util, used for teamRatings
	_, betaSquared := getBetas(options)
	u := util.NewWithOptions(&util.Options{
		BetaSquared: ptr.Float64(betaSquared),
	})

	// This is used at the end to normalize the results
	n := float64(len(teams))
	denom := (n * (n - 1)) / 2

	// Initialize the stats helpers. This is used to calculate the CDF
	// of the normal distribution
	normal := distuv.Normal{
		Mu:    0,
		Sigma: 1,
	}

	// Pre-calculate the team ratings
	teamRatings := u.TeamRating(teams, options)

	// Initialize the results
	returning := make([]float64, len(teams))

	// Outer loop, iterate over each team
	for i, outerTeamRating := range teamRatings {
		// Initialize the inner loop and run it
		var innerTeamRatings []types.TeamRating
		for j, innerTeamRating := range teamRatings {
			if i != j {
				innerTeamRatings = append(innerTeamRatings, innerTeamRating)
			}
		}

		var prediction float64
		for _, innerTeamRating := range innerTeamRatings {
			muA := outerTeamRating.TeamMu
			muB := innerTeamRating.TeamMu
			betaSQ := u.BetaSquared
			sigmaSqA := outerTeamRating.TeamSigmaSquared
			sigmaSqB := innerTeamRating.TeamSigmaSquared

			prediction += normal.CDF((muA - muB) / math.Sqrt(n*betaSQ+sigmaSqA*sigmaSqA+sigmaSqB*sigmaSqB))
		}

		returning[i] = prediction / denom
	}

	return returning
}

// PredictDraw returns the probability of each team drawing
func PredictDraw(teams []types.Team, options *types.OpenSkillOptions) float64 {
	if len(teams) == 1 {
		return 1.0
	}

	// Initialize util, used for teamRatings
	beta, betaSquared := getBetas(options)
	u := util.NewWithOptions(&util.Options{
		BetaSquared: ptr.Float64(betaSquared),
	})

	// This is used at the end to normalize the results
	n := float64(len(teams))
	denom := n * (n - 1)
	if n <= 2 {
		denom /= 2
	}

	// Initialize the stats helpers. This is used to calculate the CDF
	// of the normal distribution
	normal := distuv.Normal{
		Mu:    0,
		Sigma: 1,
	}

	// Pre-calculate the team ratings
	teamRatings := u.TeamRating(teams, options)

	flattenedLength := float64(len(flattenTeams(teams)))
	drawMargin := math.Sqrt(flattenedLength) * beta * normal.Quantile((1+1/n)/2)

	// Initialize the results
	returning := make([]float64, len(teams))
	for i, outerTeamRating := range teamRatings {
		var innerTeamRatings []types.TeamRating
		for j, innerTeamRating := range teamRatings {
			if i != j {
				innerTeamRatings = append(innerTeamRatings, innerTeamRating)
			}
		}

		var prediction float64
		for _, innerTeamRating := range innerTeamRatings {
			muA := outerTeamRating.TeamMu
			muB := innerTeamRating.TeamMu
			betaSQ := u.BetaSquared
			sigmaSqA := outerTeamRating.TeamSigmaSquared
			sigmaSqB := innerTeamRating.TeamSigmaSquared
			sigmaBar := math.Sqrt(n*betaSQ + sigmaSqA*sigmaSqA + sigmaSqB*sigmaSqB)
			prediction += normal.CDF((drawMargin-muA+muB)/sigmaBar) -
				normal.CDF((muA-muB-drawMargin)/sigmaBar)
		}

		returning[i] = prediction
	}

	sumReturning := 0.0
	for _, v := range returning {
		sumReturning += v
	}
	return math.Abs(sumReturning) / denom
}

// PredictRank returns the probability of each team ranking
func PredictRank(teams []types.Team, _ *types.OpenSkillOptions) ([]int64, []float64) {
	// It's interesting to note that here we return 0 to the probability of being first
	// when you're alone because the convention is that when you're alone you always draw
	if len(teams) == 1 {
		return []int64{1}, []float64{0.0}
	}

	return []int64{}, []float64{}
}

// flattenTeams takes a slice of teams and flattens it into a slice of ratings
func flattenTeams(teams []types.Team) []types.Rating {
	var ratings []types.Rating
	for _, team := range teams {
		for _, rating := range team {
			ratings = append(ratings, rating)
		}
	}

	return ratings
}
