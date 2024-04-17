package util_test

import (
	"testing"

	_is "github.com/matryer/is"

	"github.com/intinig/go-openskill/models"
	"github.com/intinig/go-openskill/rating"
	"github.com/intinig/go-openskill/types"
)

// getTeam returns a team of n players all with the same rating
func getTeam(n int) types.Team {
	var t types.Team
	for i := 0; i < n; i++ {
		t = append(t, rating.New())
	}
	return t
}

func TestTeamRatingAggregatesAllPlayersInATeam(t *testing.T) {
	t.Parallel()
	is := _is.New(t)

	t1 := getTeam(1)
	t2 := getTeam(2)

	model := models.NewPlackettLuce(nil)
	results := model.U.TeamRating([]types.Team{t1, t2}, &types.OpenSkillOptions{})
	is.Equal(results, []types.TeamRating{
		{
			TeamMu:           25.0,
			TeamSigmaSquared: 69.44444444444446,
			Team:             t1,
			Rank:             0,
		},
		{
			TeamMu:           50.0,
			TeamSigmaSquared: 138.8888888888889,
			Team:             t2,
			Rank:             1,
		},
	})
}

func TestTeamRatingAggregatesAllPlayersInATeamWithWeirdRankings(t *testing.T) {
	t.Parallel()

	is := _is.New(t)

	t1 := getTeam(1)
	t2 := getTeam(2)

	model := models.NewPlackettLuce(nil)

	results := model.U.TeamRating([]types.Team{t1, t2}, &types.OpenSkillOptions{
		Rank: []int{1, 0},
	})
	is.Equal(results, []types.TeamRating{
		{
			TeamMu:           25.0,
			TeamSigmaSquared: 69.44444444444446,
			Team:             t1,
			Rank:             1,
		},
		{
			TeamMu:           50.0,
			TeamSigmaSquared: 138.8888888888889,
			Team:             t2,
			Rank:             0,
		},
	})
}

func TestTeamRatingAggregatesFiveVsFiveTeams(t *testing.T) {
	t.Parallel()

	is := _is.New(t)

	t1 := getTeam(5)
	t2 := getTeam(5)

	model := models.NewPlackettLuce(nil)
	results := model.U.TeamRating([]types.Team{t1, t2}, nil)
	is.Equal(results, []types.TeamRating{
		{
			TeamMu:           125.0,
			TeamSigmaSquared: 347.2222222222223,
			Team:             t1,
			Rank:             0,
		},
		{
			TeamMu:           125.0,
			TeamSigmaSquared: 347.2222222222223,
			Team:             t2,
			Rank:             1,
		},
	})
}

func TestUtilCComputations(t *testing.T) {
	t.Parallel()

	is := _is.New(t)

	// I am using NewPlackettLuce here to get out
	// the default values for all the constants
	model := models.NewPlackettLuce(nil)

	t1 := getTeam(1)
	t2 := getTeam(2)

	tr := model.U.TeamRating([]types.Team{t2, t1}, nil)
	c := model.U.C(tr)
	is.Equal(c, 15.590239111558091)
}

func TestUtilCFor5v5(t *testing.T) {
	t.Parallel()

	is := _is.New(t)

	// I am using NewPlackettLuce here to get out
	// the default values for all the constants
	model := models.NewPlackettLuce(nil)

	t1 := getTeam(5)
	t2 := getTeam(5)

	tr := model.U.TeamRating([]types.Team{t2, t1}, nil)
	c := model.U.C(tr)
	is.Equal(c, 27.003086243366084)
}

func TestUtilAComputesAsExpected(t *testing.T) {
	t.Parallel()

	is := _is.New(t)

	t1 := getTeam(1)
	t2 := getTeam(2)

	// I am using NewPlackettLuce here to get out
	// the default values for all the constants
	model := models.NewPlackettLuce(nil)

	tr := model.U.TeamRating([]types.Team{t2, t1}, nil)
	a := model.U.A(tr)

	is.Equal(a, []int{1, 1})
}

func TestUtilAComputesAsExpectedForReversedTeamOrder(t *testing.T) {
	t.Parallel()

	is := _is.New(t)

	t1 := getTeam(1)
	t2 := getTeam(2)

	// I am using NewPlackettLuce here to get out
	// the default values for all the constants
	model := models.NewPlackettLuce(nil)

	tr := model.U.TeamRating([]types.Team{t1, t2}, &types.OpenSkillOptions{
		Rank: []int{2, 1},
	})
	a := model.U.A(tr)

	is.Equal(a, []int{1, 1})
}

func TestUtilACountsOneTeamPerRannk(t *testing.T) {
	t.Parallel()

	is := _is.New(t)

	t1 := getTeam(1)
	t2 := getTeam(2)
	t3 := getTeam(2)
	t4 := getTeam(1)

	// I am using NewPlackettLuce here to get out
	// the default values for all the constants
	model := models.NewPlackettLuce(nil)
	tr := model.U.TeamRating([]types.Team{t1, t2, t3, t4}, nil)
	a := model.U.A(tr)

	is.Equal(a, []int{1, 1, 1, 1})
}

func TestUtilACountsHowManyTeamsShareARank(t *testing.T) {
	t.Parallel()

	is := _is.New(t)

	t1 := getTeam(1)
	t2 := getTeam(2)
	t3 := getTeam(2)
	t4 := getTeam(1)

	// I am using NewPlackettLuce here to get out
	// the default values for all the constants
	model := models.NewPlackettLuce(nil)

	tr := model.U.TeamRating([]types.Team{t1, t2, t3, t4}, &types.OpenSkillOptions{
		Rank: []int{1, 1, 1, 4},
	})

	a := model.U.A(tr)

	is.Equal(a, []int{3, 3, 3, 1})
}

func TestUtilSumQComputesAsExpected(t *testing.T) {
	t.Parallel()

	is := _is.New(t)

	t1 := getTeam(1)
	t2 := getTeam(2)

	// I am using NewPlackettLuce here to get out
	// the default values for all the constants
	model := models.NewPlackettLuce(nil)

	tr := model.U.TeamRating([]types.Team{t1, t2}, nil)
	c := model.U.C(tr)
	q := model.U.SumQ(tr, c)

	is.Equal(q, []float64{29.67892702634643, 24.70819334370875})
}

func TestUtilSumQComputesAsExpectedWithUnreversedRanks(t *testing.T) {
	t.Parallel()

	is := _is.New(t)

	t1 := getTeam(1)
	t2 := getTeam(1)

	// I am using NewPlackettLuce here to get out
	// the default values for all the constants
	model := models.NewPlackettLuce(nil)

	tr := model.U.TeamRating([]types.Team{t1, t2}, nil)
	c := model.U.C(tr)
	q := model.U.SumQ(tr, c)

	is.Equal(q, []float64{13.336621888349223, 6.668310944174611})
}

func TestUtilSumQComputesAsExpectedWithReversedRank(t *testing.T) {
	t.Parallel()

	is := _is.New(t)

	t1 := getTeam(1)
	t2 := getTeam(1)

	// I am using NewPlackettLuce here to get out
	// the default values for all the constants
	model := models.NewPlackettLuce(&types.OpenSkillOptions{Rank: []int{2, 1}})

	tr := model.U.TeamRating([]types.Team{t1, t2}, &types.OpenSkillOptions{
		Rank: []int{2, 1},
	})
	c := model.U.C(tr)
	q := model.U.SumQ(tr, c)

	is.Equal(q, []float64{6.668310944174611, 13.336621888349223})
}

func TestUtilSumQComputesFiveVsFive(t *testing.T) {
	t.Parallel()

	is := _is.New(t)

	t1 := getTeam(5)
	t2 := getTeam(5)

	// I am using NewPlackettLuce here to get out
	// the default values for all the constants
	model := models.NewPlackettLuce(nil)

	tr := model.U.TeamRating([]types.Team{t1, t2}, nil)
	c := model.U.C(tr)
	q := model.U.SumQ(tr, c)

	is.Equal(q, []float64{204.8437881059862, 102.4218940529931})
}

func TestUtilRankingsRanksAZeroElementArray(t *testing.T) {
	t.Parallel()

	is := _is.New(t)

	model := models.NewPlackettLuce(nil)
	ranks := model.U.Rankings([]types.Team{}, nil)
	is.Equal(ranks, []int{})
}

func TestUtilRankingsRanksASingleElementArray(t *testing.T) {
	t.Parallel()

	is := _is.New(t)

	t1 := getTeam(1)

	model := models.NewPlackettLuce(nil)
	ranks := model.U.Rankings([]types.Team{t1}, nil)
	is.Equal(ranks, []int{0})
}

func TestUtilRankingsRanksGivenANilRanksArray(t *testing.T) {
	t.Parallel()

	is := _is.New(t)

	t1 := getTeam(1)
	t2 := getTeam(1)
	t3 := getTeam(1)
	t4 := getTeam(1)

	model := models.NewPlackettLuce(nil)
	ranks := model.U.Rankings([]types.Team{t1, t2, t3, t4}, nil)
	is.Equal(ranks, []int{0, 1, 2, 3})
}

func TestUtilRankingsRanksGivenIncrementals(t *testing.T) {
	t.Parallel()

	is := _is.New(t)

	t1 := getTeam(1)
	t2 := getTeam(1)
	t3 := getTeam(1)
	t4 := getTeam(1)

	model := models.NewPlackettLuce(nil)
	ranks := model.U.Rankings([]types.Team{t1, t2, t3, t4}, []int{1, 2, 3, 4})
	is.Equal(ranks, []int{0, 1, 2, 3})
}

func TestUtilRankingsRanksWitTiesInStart(t *testing.T) {
	t.Parallel()

	is := _is.New(t)

	t1 := getTeam(1)
	t2 := getTeam(1)
	t3 := getTeam(1)
	t4 := getTeam(1)

	model := models.NewPlackettLuce(nil)
	ranks := model.U.Rankings([]types.Team{t1, t2, t3, t4}, []int{1, 1, 3, 4})
	is.Equal(ranks, []int{0, 0, 2, 3})
}

func TestUtilRankingsWithTiesInEnd(t *testing.T) {
	t.Parallel()

	is := _is.New(t)

	t1 := getTeam(1)
	t2 := getTeam(1)
	t3 := getTeam(1)
	t4 := getTeam(1)

	model := models.NewPlackettLuce(nil)
	ranks := model.U.Rankings([]types.Team{t1, t2, t3, t4}, []int{1, 2, 3, 3})
	is.Equal(ranks, []int{0, 1, 2, 2})
}

func TestUtilRankingsWithTiesInTheMiddle(t *testing.T) {
	t.Parallel()

	is := _is.New(t)

	t1 := getTeam(1)
	t2 := getTeam(1)
	t3 := getTeam(1)
	t4 := getTeam(1)

	model := models.NewPlackettLuce(nil)
	ranks := model.U.Rankings([]types.Team{t1, t2, t3, t4}, []int{1, 2, 2, 4})
	is.Equal(ranks, []int{0, 1, 1, 3})
}

func TestUtilRankingsWithSparseScores(t *testing.T) {
	t.Parallel()

	is := _is.New(t)

	t1 := getTeam(1)
	t2 := getTeam(1)
	t3 := getTeam(1)
	t4 := getTeam(1)
	t5 := getTeam(1)

	model := models.NewPlackettLuce(nil)
	ranks := model.U.Rankings([]types.Team{t1, t2, t3, t4, t5}, []int{14, 32, 47, 47, 48})
	is.Equal(ranks, []int{0, 1, 2, 2, 4})
}
