package models_test

import (
	"testing"

	_is "github.com/matryer/is"

	"github.com/intinig/go-openskill/models"
	"github.com/intinig/go-openskill/ptr"
	"github.com/intinig/go-openskill/rating"
	"github.com/intinig/go-openskill/types"
)

func TestPlackettLuceInitializationEpsilon(t *testing.T) {
	t.Parallel()
	is := _is.New(t)
	model := models.NewPlackettLuce(&types.OpenSkillOptions{})
	is.Equal(model.Epsilon, 0.0001)

	model = models.NewPlackettLuce(&types.OpenSkillOptions{
		Epsilon: ptr.Float64(0.00002),
	})
	is.Equal(model.Epsilon, 0.00002)
}

func TestPlackettLuceInitializationZ(t *testing.T) {
	t.Parallel()
	is := _is.New(t)
	model := models.NewPlackettLuce(&types.OpenSkillOptions{})
	is.Equal(model.Z, 3)

	model = models.NewPlackettLuce(&types.OpenSkillOptions{
		Z: ptr.Int(5),
	})
	is.Equal(model.Z, 5)
}

func TestPlackettLuceInitializationMu(t *testing.T) {
	t.Parallel()
	is := _is.New(t)
	model := models.NewPlackettLuce(&types.OpenSkillOptions{})
	is.Equal(model.Mu, 25.0)

	model = models.NewPlackettLuce(&types.OpenSkillOptions{
		Mu: ptr.Float64(100.0),
	})
	is.Equal(model.Mu, 100.0)
}

func TestPlackettLuceInitializationSigma(t *testing.T) {
	t.Parallel()
	is := _is.New(t)
	model := models.NewPlackettLuce(&types.OpenSkillOptions{})
	is.Equal(model.Sigma, 25.0/3.0)

	model = models.NewPlackettLuce(&types.OpenSkillOptions{
		Sigma: ptr.Float64(100.0),
	})
	is.Equal(model.Sigma, 100.0)
}

func TestNewPlackettLuceInitializationBeta(t *testing.T) {
	t.Parallel()
	is := _is.New(t)
	model := models.NewPlackettLuce(&types.OpenSkillOptions{})
	is.Equal(model.Beta, model.Sigma/2.0)

	model = models.NewPlackettLuce(&types.OpenSkillOptions{
		Beta: ptr.Float64(11.11),
	})
	is.Equal(model.Beta, 11.11)
}

func TestPlackettLuceSoloGameDoesNotChangeRating(t *testing.T) {
	t.Parallel()
	is := _is.New(t)

	r := rating.New()
	model := models.NewPlackettLuce(nil)
	teams := model.Rate([]types.Team{{r}}, nil)
	is.Equal(teams, []types.Team{{r}})
}

func Test2PlayersFreeForAll(t *testing.T) {
	t.Parallel()
	is := _is.New(t)

	model := models.NewPlackettLuce(nil)
	teams := model.Rate([]types.Team{
		{rating.New()},
		{rating.New()},
	}, nil)
	is.Equal(teams, []types.Team{
		{types.Rating{Mu: 27.63523138347365, Sigma: 8.065506316323548, Z: 3}},
		{types.Rating{Mu: 22.36476861652635, Sigma: 8.065506316323548, Z: 3}},
	})
}

func Test3PlayersFreeForAll(t *testing.T) {
	t.Parallel()
	is := _is.New(t)

	model := models.NewPlackettLuce(nil)
	teams := model.Rate([]types.Team{
		{rating.New()},
		{rating.New()},
		{rating.New()},
	}, nil)
	is.Equal(teams, []types.Team{
		{types.Rating{Mu: 27.868876552746237, Sigma: 8.204837030780652, Z: 3}},
		{types.Rating{Mu: 25.717219138186557, Sigma: 8.057829747583874, Z: 3}},
		{types.Rating{Mu: 21.413904309067206, Sigma: 8.057829747583874, Z: 3}},
	})
}

func Test4PlayersFreeForAll(t *testing.T) {
	t.Parallel()
	is := _is.New(t)

	model := models.NewPlackettLuce(nil)
	teams := model.Rate([]types.Team{
		{rating.New()},
		{rating.New()},
		{rating.New()},
		{rating.New()},
	}, nil)
	is.Equal(teams, []types.Team{
		{types.Rating{Mu: 27.795084971874736, Sigma: 8.263160757613477, Z: 3}},
		{types.Rating{Mu: 26.552824984374855, Sigma: 8.179213704945203, Z: 3}},
		{types.Rating{Mu: 24.68943500312503, Sigma: 8.083731307186588, Z: 3}},
		{types.Rating{Mu: 20.96265504062538, Sigma: 8.083731307186588, Z: 3}},
	})
}

func Test5PlayersFreeForAll(t *testing.T) {
	t.Parallel()
	is := _is.New(t)

	model := models.NewPlackettLuce(nil)
	teams := model.Rate([]types.Team{
		{rating.New()},
		{rating.New()},
		{rating.New()},
		{rating.New()},
		{rating.New()},
	}, nil)
	is.Equal(teams, []types.Team{
		{types.Rating{Mu: 27.666666666666668, Sigma: 8.290556877154474, Z: 3}},
		{types.Rating{Mu: 26.833333333333332, Sigma: 8.240145629781066, Z: 3}},
		{types.Rating{Mu: 25.72222222222222, Sigma: 8.179996679645559, Z: 3}},
		{types.Rating{Mu: 24.055555555555557, Sigma: 8.111796013701358, Z: 3}},
		{types.Rating{Mu: 20.72222222222222, Sigma: 8.111796013701358, Z: 3}},
	})
}

func Test3TeamsWithDifferentPlayersNumbers(t *testing.T) {
	t.Parallel()
	is := _is.New(t)

	model := models.NewPlackettLuce(nil)
	teams := model.Rate([]types.Team{
		{rating.New(), rating.New(), rating.New()},
		{rating.New()},
		{rating.New(), rating.New()},
	}, nil)
	is.Equal(teams, []types.Team{
		{
			types.Rating{Mu: 25.939870821784513, Sigma: 8.247641552260456, Z: 3},
			types.Rating{Mu: 25.939870821784513, Sigma: 8.247641552260456, Z: 3},
			types.Rating{Mu: 25.939870821784513, Sigma: 8.247641552260456, Z: 3},
		},
		{types.Rating{Mu: 27.21366020491262, Sigma: 8.274321317985242, Z: 3}},
		{
			types.Rating{Mu: 21.84646897330287, Sigma: 8.213058173195341, Z: 3},
			types.Rating{Mu: 21.84646897330287, Sigma: 8.213058173195341, Z: 3},
		},
	})
}

func TestSeries(t *testing.T) {
	t.Parallel()
	is := _is.New(t)

	model := models.NewPlackettLuce(nil)

	p00 := rating.New()
	p10 := rating.New()
	p20 := rating.New()
	p30 := rating.New()
	p40 := rating.New()

	m1 := rating.Rate([]types.Team{
		{p00}, {p10}, {p20}, {p30}, {p40},
	}, &types.OpenSkillOptions{
		Model: model,
		Score: []int{9, 7, 7, 5, 5},
	})

	p01 := m1[0][0]
	p11 := m1[1][0]
	p21 := m1[2][0]
	p31 := m1[3][0]
	p41 := m1[4][0]

	p02 := p01
	p32 := p31

	m2 := rating.Rate([]types.Team{
		{p41}, {p21}, {p11},
	}, &types.OpenSkillOptions{
		Model: model,
		Score: []int{9, 5, 5},
	})

	p42 := m2[0][0]
	p22 := m2[1][0]
	p12 := m2[2][0]

	p43 := p42

	m3 := rating.Rate([]types.Team{
		{p32}, {p12}, {p22}, {p02},
	}, &types.OpenSkillOptions{
		Model: model,
		Score: []int{9, 9, 7, 7},
	})

	p33 := m3[0][0]
	p13 := m3[1][0]
	p23 := m3[2][0]
	p03 := m3[3][0]

	is.Equal(p03.Mu, 26.3537611030628)
	is.Equal(p03.Sigma, 8.111027060497456)
	is.Equal(p13.Mu, 24.618479788611904)
	is.Equal(p13.Sigma, 7.905335509558602)
	is.Equal(p23.Mu, 23.065819512381218)
	is.Equal(p23.Sigma, 7.822005594848687)
	is.Equal(p33.Mu, 24.47633240286842)
	is.Equal(p33.Sigma, 8.106111471280611)
	is.Equal(p43.Mu, 26.385499684561076)
	is.Equal(p43.Sigma, 8.05409080928062)
}
