package rating_test

import (
	"testing"

	_is "github.com/matryer/is"

	"github.com/intinig/go-openskill/ptr"
	"github.com/intinig/go-openskill/rating"
	"github.com/intinig/go-openskill/test"
	"github.com/intinig/go-openskill/types"
)

func TestPredictsWinOutcomeForTwoTeams(t *testing.T) {
	t.Parallel()
	is := _is.New(t)
	teams := []types.Team{
		{
			test.PredictWinTeams["a1"],
			test.PredictWinTeams["a2"],
		},
		{
			test.PredictWinTeams["b1"],
			test.PredictWinTeams["b2"],
		},
	}

	probs := rating.PredictWin(teams, nil)
	is.Equal(probs, []float64{0.34641823958165474, 0.6535817604183453})
}

func TestPredictsWinIgnoresRankings(t *testing.T) {
	t.Parallel()
	is := _is.New(t)
	teams := []types.Team{
		{
			test.PredictWinTeams["a1"],
		},
		{
			test.PredictWinTeams["b1"],
		},
		{
			test.PredictWinTeams["b2"],
		},
	}

	p1 := rating.PredictWin(teams, &types.OpenSkillOptions{
		Rank: []int{2, 1, 3},
	})

	p2 := rating.PredictWin(teams, &types.OpenSkillOptions{
		Rank: []int{1, 2, 3},
	})

	is.Equal(p1, p2)
}

func TestPredictsWinOutcomeForMultipleAsymmetricTeams(t *testing.T) {
	t.Parallel()
	is := _is.New(t)
	teams := []types.Team{
		{
			test.PredictWinTeams["a1"],
			test.PredictWinTeams["a2"],
		},
		{
			test.PredictWinTeams["b1"],
			test.PredictWinTeams["b2"],
		},
		{
			test.PredictWinTeams["a2"],
		},
		{
			test.PredictWinTeams["b2"],
		},
	}

	probs := rating.PredictWin(teams, nil)
	is.Equal(probs, []float64{0.26135159416422216, 0.4111743094338915, 0.17509059831123944, 0.15238349809064683})
}

func TestPredictsWinWith3PlayerNewbieFFA(t *testing.T) {
	t.Parallel()
	is := _is.New(t)
	teams := []types.Team{
		{
			test.PredictWinTeams["a1"],
		},
		{
			test.PredictWinTeams["a1"],
		},
		{
			test.PredictWinTeams["a1"],
		},
	}

	probs := rating.PredictWin(teams, nil)
	is.Equal(probs, []float64{0.3333333333333333, 0.3333333333333333, 0.3333333333333333})
}

func TestPredictsWinWith4PlayerNewbieFFA(t *testing.T) {
	t.Parallel()
	is := _is.New(t)
	teams := []types.Team{
		{
			test.PredictWinTeams["a1"],
		},
		{
			test.PredictWinTeams["a1"],
		},
		{
			test.PredictWinTeams["a1"],
		},
		{
			test.PredictWinTeams["a1"],
		},
	}

	probs := rating.PredictWin(teams, nil)
	is.Equal(probs, []float64{0.25, 0.25, 0.25, 0.25})
}

func TestPredictsWinWith5PlayerNewbieFFA(t *testing.T) {
	t.Parallel()
	is := _is.New(t)
	teams := []types.Team{
		{
			test.PredictWinTeams["a1"],
		},
		{
			test.PredictWinTeams["a1"],
		},
		{
			test.PredictWinTeams["a1"],
		},
		{
			test.PredictWinTeams["a1"],
		},
		{
			test.PredictWinTeams["a1"],
		},
	}

	probs := rating.PredictWin(teams, nil)
	is.Equal(probs, []float64{0.2, 0.2, 0.2, 0.2, 0.2})
}

func TestPredictsWinWith5PlayerNewbieFFAWithImpostor(t *testing.T) {
	t.Parallel()
	is := _is.New(t)
	teams := []types.Team{
		{
			test.PredictWinTeams["a1"],
		},
		{
			test.PredictWinTeams["a1"],
		},
		{
			test.PredictWinTeams["a1"],
		},
		{
			test.PredictWinTeams["a2"],
		},
		{
			test.PredictWinTeams["a1"],
		},
	}

	probs := rating.PredictWin(teams, nil)
	is.Equal(probs, []float64{0.19603741652263795, 0.19603741652263795, 0.19603741652263795, 0.2158503339094482, 0.19603741652263795})
}

func TestPredictDraw100PercentForSolitare(t *testing.T) {
	t.Parallel()
	is := _is.New(t)
	teams := []types.Team{
		{
			test.PredictWinTeams["a1"],
		},
	}

	probs := rating.PredictDraw(teams, nil)
	is.Equal(probs, 1.0)
}

func TestPredictDraw100PercentForSelfVsSelf(t *testing.T) {
	t.Parallel()
	is := _is.New(t)
	teams := []types.Team{
		{
			test.PredictDrawTeams["b1"],
		},
		{
			test.PredictDrawTeams["b1"],
		},
	}

	probs := rating.PredictDraw(teams, nil)
	is.Equal(probs, 1.0)
}

func TestPredictDrawForTwoTeams(t *testing.T) {
	t.Parallel()
	is := _is.New(t)
	teams := []types.Team{
		{
			test.PredictDrawTeams["a1"],
			test.PredictDrawTeams["a2"],
		},
		{
			test.PredictDrawTeams["b1"],
			test.PredictDrawTeams["b2"],
		},
	}

	probs := rating.PredictDraw(teams, nil)
	is.Equal(probs, 0.1260703143635969)
}

func TestPredictDrawForThreeAsymmetricTeams(t *testing.T) {
	t.Parallel()
	is := _is.New(t)
	teams := []types.Team{
		{
			test.PredictDrawTeams["a1"],
			test.PredictDrawTeams["a2"],
		},
		{
			test.PredictDrawTeams["b1"],
			test.PredictDrawTeams["b2"],
		},
		{
			test.PredictDrawTeams["a1"],
		},
		{
			test.PredictDrawTeams["a2"],
		},
		{
			test.PredictDrawTeams["b1"],
		},
	}

	probs := rating.PredictDraw(teams, nil)
	is.Equal(probs, 0.04322122887507539)
}

func TestPredictRank(t *testing.T) {
	t.Parallel()
	t.Skip("This test is not working as expected. Need to investigate.")
	is := _is.New(t)

	a1 := rating.NewWithOptions(&types.OpenSkillOptions{
		Mu:    ptr.Float64(34.0),
		Sigma: ptr.Float64(0.25),
	})
	a2 := rating.NewWithOptions(&types.OpenSkillOptions{
		Mu:    ptr.Float64(32.0),
		Sigma: ptr.Float64(0.25),
	})
	a3 := rating.NewWithOptions(&types.OpenSkillOptions{
		Mu:    ptr.Float64(34.0),
		Sigma: ptr.Float64(0.25),
	})

	b1 := rating.NewWithOptions(&types.OpenSkillOptions{
		Mu:    ptr.Float64(24.0),
		Sigma: ptr.Float64(0.5),
	})
	b2 := rating.NewWithOptions(&types.OpenSkillOptions{
		Mu:    ptr.Float64(22.0),
		Sigma: ptr.Float64(0.5),
	})
	b3 := rating.NewWithOptions(&types.OpenSkillOptions{
		Mu:    ptr.Float64(20.0),
		Sigma: ptr.Float64(0.5),
	})
	team1 := []types.Rating{a1, b1}
	team2 := []types.Rating{a2, b2}
	team3 := []types.Rating{a3, b3}

	_, ranks := rating.PredictRank([]types.Team{team1, team2, team3}, nil)
	total := 0.0
	for _, rank := range ranks {
		total += rank
	}
	draw := rating.PredictDraw([]types.Team{team1, team2, team3}, nil)
	is.Equal(total+draw, float64(1))
}
