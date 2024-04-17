package unwind

import (
	"sort"

	"github.com/intinig/go-openskill/types"
)

type Zip struct {
	Team   types.Team
	Rank   int
	DeRank int
}

type Zips []Zip

func (z Zips) Len() int {
	return len(z)
}

func (z Zips) Less(i, j int) bool {
	return z[i].Rank < z[j].Rank
}

func (z Zips) Swap(i, j int) {
	z[i], z[j] = z[j], z[i]
}

// Teams unwinds a set of teams and their ranks and returns both the unwound
// teams and the tenet to reverse the operation
func Teams(src []types.Team, rank []int) ([]types.Team, []int) {
	zips := make(Zips, len(src))
	for i := range src {
		zips[i] = Zip{
			Team:   src[i],
			Rank:   rank[i],
			DeRank: i,
		}
	}

	sort.Sort(zips)

	dest := make([]types.Team, len(src))
	tenet := make([]int, len(rank))

	for i, z := range zips {
		dest[i] = z.Team
		tenet[i] = z.DeRank
	}

	return dest, tenet
}
