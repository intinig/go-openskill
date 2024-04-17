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
	// Z is the number of standard deviations a rating can deviate from the mean
	// before it is considered to be outside the "normal" range of skill.
	// The default value is 3.29, which covers approximately 99.7% of the normal
	// distribution.
	Z *int
	// Mu is the mean of the rating distribution.
	// The default value is 25.0.
	Mu *float64
	// Sigma is the standard deviation of the rating distribution.
	// The default value is 25.0 / 3.29 = 7.59.
	Sigma *float64
	// Epsilon is the minimum change in rating required for a player to gain or
	// lose rating. The default value is 0.000001.
	Epsilon *float64
	// Beta is the "precision" of the rating distribution.
	// The default value is 0.5.
	Beta *float64
	// Model is the rating model to use.
	// The default value is Glicko2.
	Model RatingModel
	// Rank is the rank of each team, where 0 is the highest rank. The default value
	// is [0, 1, ...]. This is only supported when passed through the Rate function
	Rank []int
	// Score is the score of each team. It is optional and used only if Rank
	// is not specified.
	Score []int
	// Weight is the weight of each player on a team
	Weight [][]float64
	// Tau is the "precision" of the rating distribution.
	// The default value is 0.5.
	Tau *float64
	// PreventSigmaIncrease is a flag that prevents sigma from increasing.
	// The default value is false.
	PreventSigmaIncrease bool
	// Gamma is a function that returns the dynamic factor for a given rating.
	Gamma func(TeamRating) float64
}

// Rating represents a rating.
type Rating struct {
	Mu    float64
	Sigma float64
	Z     int
}
