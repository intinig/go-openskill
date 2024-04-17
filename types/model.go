package types

type RatingModel interface {
	Rate(teams []Team, options *OpenSkillOptions) []Team
}
