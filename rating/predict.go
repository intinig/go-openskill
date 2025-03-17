package rating

import (
	"math"
	"sort"

	"gonum.org/v1/gonum/stat/distuv"

	"github.com/intinig/go-openskill/ptr"
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
func PredictRank(teams []types.Team, options *types.OpenSkillOptions) ([]int64, []float64) {
	// If there is only one team, it will always be ranked first
	if len(teams) == 1 {
		return []int64{1}, []float64{0.0}
	}

	// Initialize util, used for teamRatings
	beta, betaSquared := getBetas(options)
	u := util.NewWithOptions(&util.Options{
		BetaSquared: ptr.Float64(betaSquared),
	})

	// Pre-calculate the team ratings
	teamRatings := u.TeamRating(teams, options)

	n := len(teams)
	winProbabilities := make([]float64, n)

	// Calculate win probabilities for each team
	for i, teamI := range teamRatings {
		teamWinProbability := 0.0
		for j, teamJ := range teamRatings {
			if i != j {
				teamWinProbability += phiMajor(
					(teamI.TeamMu - teamJ.TeamMu) /
						math.Sqrt(2*beta*beta+teamI.TeamSigmaSquared+teamJ.TeamSigmaSquared),
				)
			}
		}
		winProbabilities[i] = teamWinProbability / float64(n-1)
	}

	// Normalize probabilities
	totalProbability := 0.0
	for _, p := range winProbabilities {
		totalProbability += p
	}
	normalizedProbabilities := make([]float64, n)
	for i, p := range winProbabilities {
		normalizedProbabilities[i] = p / totalProbability
	}

	// Sort teams by probabilities in descending order
	type teamProbability struct {
		index       int
		probability float64
	}
	sortedTeams := make([]teamProbability, n)
	for i, p := range normalizedProbabilities {
		sortedTeams[i] = teamProbability{i, p}
	}
	sort.Slice(sortedTeams, func(i, j int) bool {
		return sortedTeams[i].probability > sortedTeams[j].probability
	})

	// Assign ranks
	ranks := make([]int64, n)
	currentRank := int64(1)
	for i, team := range sortedTeams {
		if i > 0 && team.probability < sortedTeams[i-1].probability {
			currentRank = int64(i + 1)
		}
		ranks[team.index] = currentRank
	}

	return ranks, normalizedProbabilities
}

// phiMajor is a helper function to calculate the CDF of the standard normal distribution
func phiMajor(x float64) float64 {
	return 0.5 * (1 + math.Erf(x/math.Sqrt2))
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
