package rating_test

import (
	"math/rand"
	"testing"

	_is "github.com/matryer/is"

	"github.com/intinig/go-openskill/models"
	"github.com/intinig/go-openskill/rating"
	"github.com/intinig/go-openskill/types"
)

func TestOrdinalConvertsARatingIntoAnOrdinal(t *testing.T) {
	t.Parallel()
	is := _is.New(t)
	r := rating.NewWithOptions(&types.OpenSkillOptions{
		Mu:    ptr.Float64(5.0),
		Sigma: ptr.Float64(2.0),
	})

	is.Equal(rating.Ordinal(r), -1.0)
}

func TestOrdinalRespectsCustomZ(t *testing.T) {
	t.Parallel()
	is := _is.New(t)
	r := rating.NewWithOptions(&types.OpenSkillOptions{
		Mu:    ptr.Float64(24.0),
		Sigma: ptr.Float64(6.0),
		Z:     ptr.Int(2),
	})

	is.Equal(rating.Ordinal(r), 12.0)
}

func TestTeamOrdinalWorksAsExpectedMaybe(t *testing.T) {
	t.Parallel()
	is := _is.New(t)
	r := rating.New()

	model := models.NewPlackettLuce(nil)

	results := model.U.TeamRating([]types.Team{{r, r, r, r}}, &types.OpenSkillOptions{})

	is.Equal(rating.TeamOrdinal(results[0]), 0.0)
}

func TestTeamOrdinalWorksWith1Person(t *testing.T) {
	t.Parallel()

	is := _is.New(t)
	model := models.NewPlackettLuce(nil)
	for i := 0; i < 10; i++ {
		r := rating.NewWithOptions(&types.OpenSkillOptions{
			Mu: ptr.Float64(rand.Float64() * 100),
		})
		results := model.U.TeamRating([]types.Team{{r}}, nil)

		is.Equal(rating.TeamOrdinal(results[0]), rating.Ordinal(r))
	}
}

func TestTeamOrdinalForNegativeMu(t *testing.T) {
	t.Parallel()
	is := _is.New(t)
	r := rating.NewWithOptions(&types.OpenSkillOptions{
		Mu:    ptr.Float64(-55.0),
		Sigma: ptr.Float64(1.0),
	})
	model := models.NewPlackettLuce(nil)
	results := model.U.TeamRating([]types.Team{{r, r}}, nil)
	is.Equal(rating.Ordinal(r), rating.TeamOrdinal(results[0]))
}
