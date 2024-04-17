package rating

import (
	"github.com/intinig/go-openskill/types"
)

// NewWithOptions returns a new Rating
// mu is the mean of the rating distribution
// sigma is the standard deviation of the rating distribution
func NewWithOptions(options *types.OpenSkillOptions) types.Rating {
	mu := options.Mu
	if mu == nil {
		mu = ptr.Float64(25.0)
	}

	z := options.Z
	if z == nil {
		z = ptr.Int(3)
	}

	sigma := options.Sigma
	if sigma == nil {
		sigma = ptr.Float64(*mu / float64(*z))
	}

	return types.Rating{
		Mu:    *mu,
		Sigma: *sigma,
		Z:     *z,
	}
}

// New returns a new Rating with default options
func New() types.Rating {
	return NewWithOptions(&types.OpenSkillOptions{})
}
