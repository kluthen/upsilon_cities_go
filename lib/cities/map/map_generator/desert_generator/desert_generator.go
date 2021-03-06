package desert_generator

import (
	"log"
	"math"
	"upsilon_cities_go/lib/cities/map/grid"
	"upsilon_cities_go/lib/cities/map/map_generator/map_level"
	"upsilon_cities_go/lib/cities/node"
	"upsilon_cities_go/lib/cities/nodetype"
	"upsilon_cities_go/lib/cities/tools"
	"upsilon_cities_go/lib/db"
)

//DesertGenerator generate desert ahah
type DesertGenerator struct {
	Width     tools.IntRange
	Range     tools.IntRange
	Disparity int
}

//Create a new desert generator with randomized conf
func Create() (mg DesertGenerator) {
	mg.Width = tools.MakeIntRange(3, tools.RandInt(3, 5))
	mg.Range = tools.MakeIntRange(3, tools.RandInt(10, 20))
	mg.Disparity = 1
	return
}

//Level of the sub generator see Generator Level
func (mg DesertGenerator) Level() map_level.GeneratorLevel {
	return map_level.Landscape
}

//Generate Will apply generator to provided grid
func (mg DesertGenerator) Generate(gd *grid.CompoundedGrid, dbh *db.Handler) error {

	width := mg.Width.Roll()
	rg := mg.Range.Roll()

	pt := tools.MakeIntRange(0, gd.Base.Size-1)

	test := 0
	// test 3 times to get the right place for a nice mountain, failure ? don't care ... :)
	for test < 3 {
		nd := node.NP(pt.Roll(), pt.Roll())
		log.Printf("DesertGenerator: Base %d set to %s", test+1, nd.String())
		if !gd.IsFilled(nd, nodetype.Ground) {
			targets := node.PointsAtDistance(nd, rg, gd.Base.Size)
			lentarget := len(targets)
			log.Printf("DesertGenerator: Found %d potential targets", lentarget)
			for i := 0; i < lentarget; i++ {
				target := targets[tools.RandInt(0, lentarget-1)]
				log.Printf("DesertGenerator: Trying with target %s", target.String())
				if !gd.IsFilled(target, nodetype.Ground) {

					dist := math.Sqrt(math.Pow(float64(target.X-nd.X), 2) + math.Pow(float64(target.Y-nd.Y), 2))

					// unit vector = { X/V(X²+Y²), Y/V(X²+Y²) }
					unitX := float64(target.X-nd.X) / dist
					unitY := float64(target.Y-nd.Y) / dist
					log.Printf("DesertGenerator: dist: %f, UnitX: %f UnitY %f", dist, unitX, unitY)

					for idx := width - mg.Disparity; idx < (rg - (width - mg.Disparity)); idx = idx + width + mg.Disparity {
						center := node.NP(int(unitX*float64(idx)), int(unitY*float64(idx)))
						center.X = center.X + nd.X
						center.Y = center.Y + nd.Y
						log.Printf("DesertGenerator: Adding circle of mountains at: %s", center.String())

						for _, nd := range node.PointsWithinInCircle(center, width, gd.Base.Size) {
							gd.SetPGT(nd.X, nd.Y, nodetype.Desert)
						}
					}

					log.Printf("DesertGenerator: Successfully added mountain width: %d, range %d, from %s, to %s", width, rg, nd.String(), target.String())
					return nil
				}
			}
		} else {
			log.Printf("DesertGenerator: Already filled, trying something else")
		}
		test++
	}

	log.Printf("DesertGenerator: Failed to add mountain width: %d, range %d ", width, rg)
	// tried three times to add a mountain range, but couldn't ... that's okay.
	return nil
}

//Name of the generator
func (mg DesertGenerator) Name() string {
	return "DesertGenerator"
}
