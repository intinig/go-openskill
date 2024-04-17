package rating_test

import (
	"testing"

	_is "github.com/matryer/is"

	"github.com/intinig/go-openskill/ptr"
	"github.com/intinig/go-openskill/rating"
	"github.com/intinig/go-openskill/types"
)

func TestRatingDefaultsMuTo25(t *testing.T) {
	t.Parallel()
	is := _is.New(t)
	r := rating.New()
	is.Equal(r.Mu, 25.0)
}

func TestRatingDefaultsSigmaTo8dot333(t *testing.T) {
	t.Parallel()
	is := _is.New(t)
	r := rating.New()
	is.Equal(r.Sigma, 25.0/3)
}

func TestNewWithMu(t *testing.T) {
	t.Parallel()
	is := _is.New(t)
	mu := 100.0
	r := rating.NewWithOptions(&types.OpenSkillOptions{
		Mu: ptr.Float64(mu),
	})
	is.Equal(r.Mu, mu)
	is.Equal(r.Sigma, mu/3)
}

func TestNewWithSigma(t *testing.T) {
	t.Parallel()
	is := _is.New(t)
	sigma := 6.283185
	r := rating.NewWithOptions(&types.OpenSkillOptions{
		Sigma: ptr.Float64(sigma),
	})
	is.Equal(r.Sigma, sigma)
}

func TestNewWithMuEqualToZero(t *testing.T) {
	t.Parallel()
	is := _is.New(t)
	mu := 0.0
	r := rating.NewWithOptions(&types.OpenSkillOptions{
		Mu: ptr.Float64(mu),
	})
	is.Equal(r.Mu, mu)
	is.Equal(r.Sigma, mu/3)
}

func TestNewWithSigmaEqualToZero(t *testing.T) {
	t.Parallel()
	is := _is.New(t)
	sigma := 0.0
	r := rating.NewWithOptions(&types.OpenSkillOptions{
		Sigma: ptr.Float64(sigma),
	})
	is.Equal(r.Mu, 25.0)
	is.Equal(r.Sigma, sigma)
}

func TestNewWithCustomZ(t *testing.T) {
	t.Parallel()
	is := _is.New(t)
	r := rating.NewWithOptions(&types.OpenSkillOptions{
		Z: ptr.Int(4),
	})
	is.Equal(r.Mu, 25.0)
	is.Equal(r.Sigma, 25.0/float64(4))
}
