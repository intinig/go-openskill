# GO-OPENSKILL
[![CircleCI](https://dl.circleci.com/status-badge/img/gh/intinig/go-openskill/tree/main.svg?style=svg)](https://dl.circleci.com/status-badge/redirect/gh/intinig/go-openskill/tree/main)

Go implementation of Weng-Lin Rating, as described at https://www.csie.ntu.edu.tw/~cjlin/papers/online_ranking/online_journal.pdf

## Speed

_Note: The following paragraph is copied verbatim from the javascript implementation and is not yet verified for this implementation._

Up to 20x faster than TrueSkill!

| Model                            | Speed (higher is better) | Variance |         Samples |
|----------------------------------|-------------------------:|---------:|----------------:|
| Openskill/bradleyTerryFull       |           62,643 ops/sec |   ±1.09% | 91 runs sampled |
| Openskill/bradleyTerryPart       |           40,152 ops/sec |   ±0.73% | 91 runs sampled |
| Openskill/thurstoneMostellerFull |           59,336 ops/sec |   ±0.74% | 93 runs sampled |
| Openskill/thurstoneMostellerPart |           38,666 ops/sec |   ±1.21% | 92 runs sampled |
| Openskill/plackettLuce           |           23,492 ops/sec |   ±0.26% | 91 runs sampled |
| TrueSkill                        |            2,962 ops/sec |   ±3.23% | 82 runs sampled |

See [this post](https://philihp.com/2020/openskill.html) for more.

## Installation

`go get github.com/intinig/go-openskill`

## Usage

Ratings are kept as an object which represent a gaussian curve, with properties where `mu` represents the _mean_, and `sigma` represents the spread or standard deviation. Create these with:

```go
package main

import (
	"github.com/intinig/go-openskill/rating"
	"github.com/intinig/go-openskill/types"
)

func main() {
	a1 := rating.New() // mu = 25, sigma = 8.333
	a2 := rating.NewWithOptions(&types.OpenSkillOptions{
		Mu: ptr.Float64(32.444),
		Sigma: ptr.Float64(5.123),
	}) // mu = 32.444, sigma = 5.123

	b1 := rating.NewWithOptions(&types.OpenSkillOptions{
		Mu: ptr.Float64(43.381),
		Sigma: ptr.Float64(2.421),
	}) // mu = 43.381, sigma = 2.421
	b2 := rating.NewWithOptions(&types.OpenSkillOptions{
		Mu: ptr.Float64(25.188),
		Sigma: ptr.Float64(6.211),
	}) // mu = 25.188, sigma = 6.211

}
```

If `a1` and `a2` are on a team, and wins against a team of `b1` and `b2`, send this into `Rate`

```go
package main

import (
	"github.com/intinig/go-openskill/rating"
	"github.com/intinig/go-openskill/types"
)

func main() {
	rating.Rate([]types.Team{
		{
			a1,
			a2,
		},
		{
			b1,
			b2
		},
		&types.OpenSkillOptions{},
	}) // []types.Team{{mu: 28.67..., sigma: 8.07...}, ...
}
```

Teams can be asymmetric, too! For example, a game like [Axis and Allies](https://en.wikipedia.org/wiki/Axis_%26_Allies) can be 3 vs 2, and this can be modeled here.

### Ranking

When displaying a rating, or sorting a list of ratings, you can use `ordinal`

```go
package main

import (
	"github.com/intinig/go-openskill/rating"
	"github.com/intinig/go-openskill/types"
)

func main() {
	a1 := rating.NewWithOptions(&types.OpenSkillOptions{
		Mu: ptr.Float64(43.07),
		Sigma: ptr.Float64(2.42),
	}) 
	ratings.Ordinal(a1) // 35.81
}
```

By default, this returns `mu - 3*sigma`, showing a rating for which there's a [99.7%](https://en.wikipedia.org/wiki/68–95–99.7_rule) likelihood the player's true rating is higher, so with early games, a player's ordinal rating will usually go up and could go up even if that player loses.

### Artificial Ranking

If your teams are listed in one order but your ranking is in a different order, for convenience you can specify a `Ranks` option, such as

```go
package main

import (
	"github.com/intinig/go-openskill/rating"
	"github.com/intinig/go-openskill/types"
)

func main() {
	a1 := rating.New()
	rating.Rate([]types.Team{{a1}, {a1}, {a1}, {a1}}, &types.OpenSkillOptions{
		Rank: []int{4, 1, 3, 2},
    }) // []types.Team{Mu: 20.963..., Sigma: 8.084...} 🐌, { Mu: 27.795, Sigma: 8.263 }🥇
}
```

It's assumed that the lower ranks are better (wins), while higher ranks are worse (losses). You can provide a `Score` instead, where lower is worse and higher is better. These can just be raw scores from the game, if you want.

Ties should have either equivalent rank or score.

```go
package main

import (
	"github.com/intinig/go-openskill/rating"
	"github.com/intinig/go-openskill/types"
)

func main() {
	a1 := rating.New()
	rating.Rate([]types.Team{{a1}, {a1}, {a1}, {a1}}, &types.OpenSkillOptions{
		Score: []int{44, 16, 23, 21},
    }) // []types.Team{Mu: 20.963..., Sigma: 8.084...} 🐌, { Mu: 27.795, Sigma: 8.263 }🥇
}
```

### Predicting Winners

For a given match of any number of teams, using `PredictWin` you can find a relative
odds that each of those teams will win.

```go
package main

import (
	"github.com/intinig/go-openskill/rating"
	"github.com/intinig/go-openskill/types"
)

func main() {
	a1 := rating.New()
	a2 := rating.NewWithOptions(&types.OpenSkillOptions{
		Mu: ptr.Float64(33.564),
		Sigma: ptr.Float64(1.123),
	}) 
    
	rating.PredictWin([]types.Team{{a1}, {a2}}) // [ 0.45110899943132493, 0.5488910005686751 ]  they add up to 1
}
```

### Predicting Draws

Also, for a given match, using `PredictDraw` you can get the relative chance that these
teams will draw. The number returned here should be treated as relative to other matches, but in reality the odds of an actual legal draw will be impacted by some meta-function based on the rules of the game.

```go
package main

import (
	"github.com/intinig/go-openskill/rating"
	"github.com/intinig/go-openskill/types"
)

func main() {
	a1 := rating.New()
	a2 := rating.NewWithOptions(&types.OpenSkillOptions{
		Mu: ptr.Float64(33.564),
		Sigma: ptr.Float64(1.123),
	}) 
    
	rating.PredictDraw([]types.Team{{a1}, {a2}}) // 0.09025530533015186
}
```

This can be used in a similar way that you might use _quality_ in TrueSkill if you were optimizing a matchmaking system, or optimizing a tournament tree structure for exciting finals and semi-finals such as in the NCAA.

### Alternative Models

By default, we use a Plackett-Luce model, which is probably good enough for most cases. When speed is an issue, the library runs faster with other models

```go
package main

import (
	"github.com/intinig/go-openskill/rating"
	"github.com/intinig/go-openskill/types"
)

func main() {
	a1 := rating.New()
	a2 := rating.NewWithOptions(&types.OpenSkillOptions{
		Mu: ptr.Float64(33.564),
		Sigma: ptr.Float64(1.123),
	})

	rating.Rate([]types.Team{{a1}, {a2}}, &types.OpenSkillOptions{
		Model: types.ModelBradleyTerryFull,
    }) 
}
```

- Bradley-Terry rating models follow a logistic distribution over a player's skill, similar to Glicko.
- Thurstone-Mosteller rating models follow a gaussian distribution, similar to TrueSkill. Gaussian CDF/PDF functions differ in implementation from system to system (they're all just chebyshev approximations anyway). The accuracy of this model isn't usually as great either, but tuning this with an alternative gamma function can improve the accuracy if you really want to get into it.
- Full pairing should have more accurate ratings over partial pairing, however in high _k_ games (like a 100+ person marathon race), Bradley-Terry and Thurstone-Mosteller models need to do a calculation of joint probability which involves is a _k_-1 dimensional integration, which is computationally expensive. Use partial pairing in this case, where players only change based on their neighbors.
- Plackett-Luce (**default**) is a generalized Bradley-Terry model for _k_ &GreaterEqual; 3 teams. It scales best.

## Implementations in other languages

- Python https://github.com/OpenDebates/openskill.py
- Kotlin https://github.com/brezinajn/openskill.kt
- Elixir https://github.com/philihp/openskill.ex
- Lua https://github.com/bstummer/openskill.lua
- Google Sheets https://docs.google.com/spreadsheets/d/12TA1ZG_qpBi4kDTclaOGB4sd5uJK8w-0My6puMd2-CY/edit?usp=sharing
- Google Apps Script https://github.com/haya14busa/gas-openskill

## TODO

Current status of porting the project through porting the test suite from the original project.
- [x] src/__tests__/util/rankings.test.ts (6.331 s)
- [ ] src/models/__tests__/index.test.ts (6.933 s)
- [ ] src/models/__tests__/thurstone-mosteller-part-series.test.ts (6.941 s)
- [ ] src/models/__tests__/thurstone-mosteller-full-series.test.ts (6.942 s)
- [x] src/__tests__/rating.test.ts (6.953 s)
- [x] src/models/__tests__/plackett-luce-series.test.ts (6.953 s)
- [ ] src/models/__tests__/bradley-terry-part-series.test.ts (6.953 s)
- [ ] src/models/__tests__/bradley-terry-part.test.ts (6.948 s)
- [x] src/__tests__/rate.test.ts (6.961 s)
- [ ] src/models/__tests__/bradley-terry-full.test.ts (6.962 s)
- [ ] src/models/__tests__/bradley-terry-full-series.test.ts (6.961 s)
- [ ] src/models/__tests__/thurstone-mosteller-full.test.ts (6.963 s)
- [x] src/__tests__/predict-win.test.ts (6.968 s)
- [ ] src/models/__tests__/thurstone-mosteller-part.test.ts (6.97 s)
- [x] src/models/__tests__/plackett-luce.test.ts (6.984 s)
- [ ] src/__tests__/predict-draw.test.ts
- [ ] src/__tests__/util/ladder-pairs.test.ts
- [x] src/__tests__/ordinal.test.ts
- [x] src/__tests__/util/util-c.test.ts
- [x] src/__tests__/util/util-a.test.ts
- [ ] src/__tests__/util/score.test.ts
- [x] src/__tests__/util/team-rating.test.ts
- [x] src/__tests__/util/util-sum-q.test.ts
