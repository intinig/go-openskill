package test

import (
	"github.com/intinig/go-openskill/ptr"
	"github.com/intinig/go-openskill/rating"
	"github.com/intinig/go-openskill/types"
)

var Teams = map[string]types.Rating{
	"a1": rating.NewWithOptions(&types.OpenSkillOptions{
		Mu:    ptr.Float64(29.182),
		Sigma: ptr.Float64(4.782),
	}),
	"b1": rating.NewWithOptions(&types.OpenSkillOptions{
		Mu:    ptr.Float64(27.174),
		Sigma: ptr.Float64(4.922),
	}),
	"c1": rating.NewWithOptions(&types.OpenSkillOptions{
		Mu:    ptr.Float64(16.672),
		Sigma: ptr.Float64(6.217),
	}),
	"d1": rating.New(),
	"e1": rating.New(),
	"f1": rating.New(),
	"w1": rating.NewWithOptions(&types.OpenSkillOptions{
		Mu:    ptr.Float64(15.0),
		Sigma: ptr.Float64(25 / 3.0),
	}),
	"x1": rating.NewWithOptions(&types.OpenSkillOptions{
		Mu:    ptr.Float64(20.0),
		Sigma: ptr.Float64(25 / 3.0),
	}),
	"y1": rating.NewWithOptions(&types.OpenSkillOptions{
		Mu:    ptr.Float64(25.0),
		Sigma: ptr.Float64(25 / 3.0),
	}),
	"z1": rating.NewWithOptions(&types.OpenSkillOptions{
		Mu:    ptr.Float64(30.0),
		Sigma: ptr.Float64(25 / 3.0),
	}),
}

var PredictWinTeams = map[string]types.Rating{
	"a1": rating.New(),
	"a2": rating.NewWithOptions(&types.OpenSkillOptions{
		Mu:    ptr.Float64(32.444),
		Sigma: ptr.Float64(5.123),
	}),
	"b1": rating.NewWithOptions(&types.OpenSkillOptions{
		Mu:    ptr.Float64(73.381),
		Sigma: ptr.Float64(1.421),
	}),
	"b2": rating.NewWithOptions(&types.OpenSkillOptions{
		Mu:    ptr.Float64(25.188),
		Sigma: ptr.Float64(6.211),
	}),
}

var PredictDrawTeams = map[string]types.Rating{
	"a1": rating.New(),
	"a2": rating.NewWithOptions(&types.OpenSkillOptions{
		Mu:    ptr.Float64(32.444),
		Sigma: ptr.Float64(1.123),
	}),
	"b1": rating.NewWithOptions(&types.OpenSkillOptions{
		Mu:    ptr.Float64(35.881),
		Sigma: ptr.Float64(0.0001),
	}),
	"b2": rating.NewWithOptions(&types.OpenSkillOptions{
		Mu:    ptr.Float64(25.188),
		Sigma: ptr.Float64(1.421),
	}),
}
