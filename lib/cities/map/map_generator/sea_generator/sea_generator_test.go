package sea_generator

import (
	"testing"
	"upsilon_cities_go/lib/cities/map/grid"
	"upsilon_cities_go/lib/cities/nodetype"
)

func TestMountainGenerator(t *testing.T) {
	sg := Create()
	gd := new(grid.CompoundedGrid)
	gd.Base = grid.Create(20, nodetype.Plain)
	gd.Delta = grid.Create(20, nodetype.None)

	sg.Generate(gd)

}