package rating

import (
	"github.com/montanaflynn/stats"

	"github.com/intinig/go-openskill/types"
)

func Ordinal(r types.Rating) float64 {
	return r.Mu - float64(r.Z)*r.Sigma
}

func TeamOrdinal(t types.TeamRating) float64 {
	teamOrdinals := make([]float64, len(t.Team))
	for i, r := range t.Team {
		teamOrdinals[i] = Ordinal(r)
	}
	// Ignoring this error because we're sure to always being good guys with this
	median, _ := stats.Median(teamOrdinals)

	return median
}
