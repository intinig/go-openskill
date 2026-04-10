package types

// Team represents a team of players
type Team []Rating

// TeamRating is an aggregated data structure used by RatingModels
type TeamRating struct {
	TeamMu           float64
	TeamSigmaSquared float64
	Team             Team
	Rank             int
}

type OpenSkillOptions struct {
	// Z is the number of standard deviations used in ordinal calculation.
	// The default value is 3, which covers approximately 99.7% of the normal
	// distribution.
	Z *int
	// Mu is the mean of the rating distribution.
	// The default value is 25.0.
	Mu *float64
	// Sigma is the standard deviation of the rating distribution.
	// The default value is Mu / Z (25.0 / 3 ≈ 8.333).
	Sigma *float64
	// Epsilon is the minimum value for the variance factor, preventing it from
	// becoming too small or negative. The default value is 0.0001.
	Epsilon *float64
	// Beta is the distance (in terms of sigma) that guarantees about an 80%
	// chance of winning. The default value is Sigma / 2.
	Beta *float64
	// Model is the rating model to use.
	// The default value is PlackettLuce.
	Model RatingModel
	// Rank is the rank of each team, where 0 is the highest rank. The default value
	// is [0, 1, ...]. This is only supported when passed through the Rate function.
	Rank []int
	// Score is the score of each team. It is optional and used only if Rank
	// is not specified. Higher scores are better.
	Score []int
	// Weight is the weight of each player on a team.
	Weight [][]float64
	// Tau is the additive dynamics parameter that prevents sigma from
	// getting too small, increasing rating change volatility.
	// The default value is nil (disabled).
	Tau *float64
	// PreventSigmaIncrease is a flag that prevents sigma from increasing.
	// The default value is false.
	PreventSigmaIncrease bool
	// Gamma is a function that returns the dynamic factor for a given team rating.
	// The default computes sqrt(TeamSigmaSquared) / c.
	Gamma func(TeamRating) float64
}

// Rating represents a rating.
type Rating struct {
	Mu    float64
	Sigma float64
	Z     int
}
