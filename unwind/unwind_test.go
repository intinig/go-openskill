package unwind_test

import (
	"testing"

	_is "github.com/matryer/is"

	"github.com/intinig/go-openskill/rating"
	"github.com/intinig/go-openskill/test"
	"github.com/intinig/go-openskill/types"
	"github.com/intinig/go-openskill/unwind"
)

//goland:noinspection GoPreferNilSlice
func TestUnwindTeamsAcceptsZeroItems(t *testing.T) {
	t.Parallel()
	is := _is.New(t)
	src := []types.Team{}
	rank := []int{}
	teams, tenet := unwind.Teams(src, rank)
	is.Equal(teams, []types.Team{})
	is.Equal(tenet, []int{})
}

func TestUnwindTeamsAcceptsOneItem(t *testing.T) {
	t.Parallel()
	is := _is.New(t)
	src := []types.Team{
		{
			rating.New(),
		},
	}
	rank := []int{0}
	teams, tenet := unwind.Teams(src, rank)
	is.Equal(teams, src)
	is.Equal(tenet, rank)
}

func TestUnwindTeamsAcceptsTwoItems(t *testing.T) {
	t.Parallel()
	is := _is.New(t)

	src := []types.Team{
		{
			test.Teams["a1"],
		},
		{
			test.Teams["b1"],
		},
	}
	rank := []int{1, 0}
	teams, tenet := unwind.Teams(src, rank)
	is.Equal(teams, []types.Team{
		{
			test.Teams["b1"],
		},
		{
			test.Teams["a1"],
		},
	})
	is.Equal(tenet, []int{1, 0})
}

func TestUnwindTeamsAcceptsThreeItems(t *testing.T) {
	t.Parallel()
	is := _is.New(t)

	src := []types.Team{
		{
			test.Teams["a1"],
		},
		{
			test.Teams["b1"],
		},
		{
			test.Teams["c1"],
		},
	}
	rank := []int{1, 2, 0}
	teams, tenet := unwind.Teams(src, rank)
	is.Equal(teams, []types.Team{
		{
			test.Teams["c1"],
		},
		{
			test.Teams["a1"],
		},
		{
			test.Teams["b1"],
		},
	})
	is.Equal(tenet, []int{2, 0, 1})
}

func TestUnwindTeamsAcceptsFourTeams(t *testing.T) {
	t.Parallel()
	is := _is.New(t)

	src := []types.Team{
		{
			test.Teams["a1"],
		},
		{
			test.Teams["b1"],
		},
		{
			test.Teams["c1"],
		},
		{
			test.Teams["d1"],
		},
	}
	rank := []int{1, 2, 3, 0}
	teams, tenet := unwind.Teams(src, rank)
	is.Equal(teams, []types.Team{
		{
			test.Teams["d1"],
		},
		{
			test.Teams["a1"],
		},
		{
			test.Teams["b1"],
		},
		{
			test.Teams["c1"],
		},
	})
	is.Equal(tenet, []int{3, 0, 1, 2})
}

func TestUnwindTeamsAcceptsDuplicateIndexes(t *testing.T) {
	t.Parallel()
	is := _is.New(t)
	src := []types.Team{
		{
			test.Teams["a1"],
		},
		{
			test.Teams["b1"],
		},
		{
			test.Teams["c1"],
		},
	}
	rank := []int{1, 2, 1}
	teams, tenet := unwind.Teams(src, rank)
	is.Equal(teams, []types.Team{
		{
			test.Teams["a1"],
		},
		{
			test.Teams["c1"],
		},
		{
			test.Teams["b1"],
		},
	})

	teams, _ = unwind.Teams(teams, tenet)
	is.Equal(teams, src)
}
