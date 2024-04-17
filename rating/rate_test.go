package rating_test

import (
	"testing"

	_is "github.com/matryer/is"

	"github.com/intinig/go-openskill/ptr"
	"github.com/intinig/go-openskill/rating"
	"github.com/intinig/go-openskill/test"
	"github.com/intinig/go-openskill/types"
)

func assertMuAndSigma(is *_is.I, team types.Rating, mu, sigma float64) {
	is.Equal(team.Mu, mu)
	is.Equal(team.Sigma, sigma)
}

func TestRateAcceptsAndRunsAPlacketLuceModelByDefault(t *testing.T) {
	t.Parallel()
	is := _is.New(t)
	teams := rating.Rate([]types.Team{
		{test.Teams["a1"]},
		{test.Teams["b1"]},
		{test.Teams["c1"]},
		{test.Teams["d1"]},
	}, nil)

	is.Equal([]types.Team{
		{types.Rating{Mu: 30.209971908310553, Sigma: 4.764898977359521, Z: 3.0}},
		{types.Rating{Mu: 27.64460833689499, Sigma: 4.882789305097372, Z: 3.0}},
		{types.Rating{Mu: 17.403586731283518, Sigma: 6.100723440599442, Z: 3.0}},
		{types.Rating{Mu: 19.214790707434826, Sigma: 7.8542613981643985, Z: 3.0}},
	}, teams)
}

func TestReversesRank(t *testing.T) {
	t.Parallel()
	is := _is.New(t)
	teams := rating.Rate([]types.Team{
		{rating.New()},
		{rating.New()},
	},
		&types.OpenSkillOptions{
			Rank: []int{2, 1},
		},
	)
	is.Equal(teams, []types.Team{
		{types.Rating{Mu: 22.36476861652635, Sigma: 8.065506316323548, Z: 3.0}},
		{types.Rating{Mu: 27.63523138347365, Sigma: 8.065506316323548, Z: 3.0}},
	})
}

func TestKeepsRank(t *testing.T) {
	t.Parallel()
	is := _is.New(t)
	teams := rating.Rate([]types.Team{
		{rating.New()},
		{rating.New()},
	},
		&types.OpenSkillOptions{
			Rank: []int{1, 2},
		},
	)
	is.Equal(teams, []types.Team{
		{types.Rating{Mu: 27.63523138347365, Sigma: 8.065506316323548, Z: 3.0}},
		{types.Rating{Mu: 22.36476861652635, Sigma: 8.065506316323548, Z: 3.0}},
	})
}

func TestAcceptsAMisorderedRankOrdering(t *testing.T) {
	t.Parallel()
	is := _is.New(t)
	teams := rating.Rate([]types.Team{
		{test.Teams["d1"]},
		{test.Teams["d1"]},
		{test.Teams["d1"]},
		{test.Teams["d1"]},
	},
		&types.OpenSkillOptions{
			Rank: []int{2, 1, 4, 3},
		},
	)
	is.Equal(teams, []types.Team{
		{types.Rating{Mu: 26.552824984374855, Sigma: 8.179213704945203, Z: 3.0}},
		{types.Rating{Mu: 27.795084971874736, Sigma: 8.263160757613477, Z: 3.0}},
		{types.Rating{Mu: 20.96265504062538, Sigma: 8.083731307186588, Z: 3.0}},
		{types.Rating{Mu: 24.68943500312503, Sigma: 8.083731307186588, Z: 3.0}},
	})
}

func TestAcceptsTeamsInRatingOrder(t *testing.T) {
	t.Parallel()
	is := _is.New(t)
	teams := rating.Rate([]types.Team{
		{test.Teams["a1"], test.Teams["d1"]},
		{test.Teams["b1"], test.Teams["e1"]},
		{test.Teams["c1"], test.Teams["f1"]},
	}, &types.OpenSkillOptions{
		Rank: []int{3, 1, 2},
	})

	assertMuAndSigma(is, teams[0][0], 27.857928218465247, 4.743791738484319)
	assertMuAndSigma(is, teams[1][0], 27.99071775460834, 4.901007097140011)
	assertMuAndSigma(is, teams[2][0], 17.60695098907354, 6.140737155130899)
	assertMuAndSigma(is, teams[0][1], 20.979038689398703, 8.129445198549202)
	assertMuAndSigma(is, teams[1][1], 27.341134074194173, 8.231039243636156)
	assertMuAndSigma(is, teams[2][1], 26.679827236407125, 8.148750467726549)
}

func TestAllowsTies(t *testing.T) {
	t.Parallel()
	a := rating.NewWithOptions(
		&types.OpenSkillOptions{
			Mu:    ptr.Float64(10),
			Sigma: ptr.Float64(8),
		})
	b := rating.NewWithOptions(
		&types.OpenSkillOptions{
			Mu:    ptr.Float64(5),
			Sigma: ptr.Float64(10),
		})
	c := rating.NewWithOptions(
		&types.OpenSkillOptions{
			Mu:    ptr.Float64(0),
			Sigma: ptr.Float64(12),
		})
	teams := rating.Rate([]types.Team{{a}, {b}, {c}}, &types.OpenSkillOptions{
		Rank: []int{1, 2, 2}})

	is := _is.New(t)
	assertMuAndSigma(is, teams[0][0], 11.942833056030613, 7.926463661123746)
	assertMuAndSigma(is, teams[1][0], 2.938193791485662, 9.65347573412201)
	assertMuAndSigma(is, teams[2][0], -1.4023734358082325, 11.323360667700934)
}

func TestAllowsTiesWithReorder(t *testing.T) {
	t.Parallel()
	a := rating.NewWithOptions(
		&types.OpenSkillOptions{
			Mu:    ptr.Float64(10),
			Sigma: ptr.Float64(8),
		})
	b := rating.NewWithOptions(
		&types.OpenSkillOptions{
			Mu:    ptr.Float64(5),
			Sigma: ptr.Float64(10),
		})
	c := rating.NewWithOptions(
		&types.OpenSkillOptions{
			Mu:    ptr.Float64(0),
			Sigma: ptr.Float64(12),
		})

	t1 := []types.Team{{b}, {c}, {a}}
	t2 := []types.Team{{a}, {b}, {c}}

	res1 := rating.Rate(t1, &types.OpenSkillOptions{
		Rank: []int{2, 2, 1}})

	res2 := rating.Rate(t2, &types.OpenSkillOptions{
		Rank: []int{1, 2, 2}})

	is := _is.New(t)
	is.Equal(res1[0], res2[1])
	is.Equal(res1[1], res2[2])
	is.Equal(res1[2], res2[0])
}

func TestFourWayTieWithNewbies(t *testing.T) {
	t.Parallel()
	is := _is.New(t)
	r := rating.New()
	teams := rating.Rate([]types.Team{
		{r}, {r}, {r}, {r},
	}, &types.OpenSkillOptions{
		Rank: []int{1, 1, 1, 1},
	})
	is.Equal(teams, []types.Team{
		{types.Rating{Mu: 25.0, Sigma: 8.263160757613477, Z: 3.0}},
		{types.Rating{Mu: 25.0, Sigma: 8.263160757613477, Z: 3.0}},
		{types.Rating{Mu: 25.0, Sigma: 8.263160757613477, Z: 3.0}},
		{types.Rating{Mu: 25.0, Sigma: 8.263160757613477, Z: 3.0}},
	})
}

func TestFixesOrdersOfTies(t *testing.T) {
	t.Parallel()
	teams := []types.Team{
		{test.Teams["w1"]}, {test.Teams["x1"]}, {test.Teams["y1"]}, {test.Teams["z1"]},
	}

	teams = rating.Rate(teams, &types.OpenSkillOptions{
		Rank: []int{2, 4, 2, 1},
	})

	is := _is.New(t)
	assertMuAndSigma(is, teams[0][0], 15.340046366255285, 8.21273604193863)
	assertMuAndSigma(is, teams[1][0], 18.007807436399276, 8.188629384105589)
	assertMuAndSigma(is, teams[2][0], 24.25804790911316, 8.166514496319483)
	assertMuAndSigma(is, teams[3][0], 32.39409828823228, 8.247276990243211)
}

func TestRunsAModelWithTiesForFirst(t *testing.T) {
	t.Parallel()
	t.Skip("ThurstoneMostellerFull not implemented yet")
}

func TestAcceptsAScoreInsteadOfARank(t *testing.T) {
	t.Parallel()
	teams := []types.Team{
		{test.Teams["e1"]}, {test.Teams["e1"]}, {test.Teams["e1"]},
	}

	teams = rating.Rate(teams, &types.OpenSkillOptions{
		Score: []int{1, 1, 1},
	})

	is := _is.New(t)
	assertMuAndSigma(is, teams[0][0], 25, 8.204837030780652)
	assertMuAndSigma(is, teams[1][0], 25, 8.204837030780652)
	assertMuAndSigma(is, teams[2][0], 25, 8.204837030780652)
}

func TestAcceptsWeightsForPartialPlay(t *testing.T) {
	t.Parallel()
	t.Skip("Needs a model that cares about partial scores")
}

func TestAcceptsATauTerm(t *testing.T) {
	t.Parallel()
	r := rating.NewWithOptions(
		&types.OpenSkillOptions{
			Mu:    ptr.Float64(25),
			Sigma: ptr.Float64(3),
		})

	teams := []types.Team{{r}, {r}}

	teams = rating.Rate(teams, &types.OpenSkillOptions{
		Tau: ptr.Float64(0.3),
	})

	is := _is.New(t)
	assertMuAndSigma(is, teams[0][0], 25.624880438870754, 2.9879993738476953)
	assertMuAndSigma(is, teams[1][0], 24.375119561129246, 2.9879993738476953)
}

func TestPreventsSigmaFromRising(t *testing.T) {
	t.Parallel()
	r := rating.NewWithOptions(
		&types.OpenSkillOptions{
			Mu:    ptr.Float64(40),
			Sigma: ptr.Float64(3),
		})

	s := rating.NewWithOptions(
		&types.OpenSkillOptions{
			Mu:    ptr.Float64(-20),
			Sigma: ptr.Float64(3),
		})

	teams := []types.Team{{r}, {s}}
	teams = rating.Rate(teams, &types.OpenSkillOptions{
		Tau:                  ptr.Float64(0.3),
		PreventSigmaIncrease: true,
	})

	is := _is.New(t)
	assertMuAndSigma(is, teams[0][0], 40.00032667136128, 3)
	assertMuAndSigma(is, teams[1][0], -20.000326671361275, 3)
}
