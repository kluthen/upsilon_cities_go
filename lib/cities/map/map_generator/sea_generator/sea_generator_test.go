package sea_generator

import (
	"testing"
	"upsilon_cities_go/lib/cities/map/grid"
	"upsilon_cities_go/lib/cities/nodetype"
	"upsilon_cities_go/lib/db"
	"upsilon_cities_go/lib/misc/config/system"
)

func TestMountainGenerator(t *testing.T) {

	system.LoadConf()
	dbh := db.NewTest()
	db.FlushDatabase(dbh)

	sg := Create()
	gd := new(grid.CompoundedGrid)
	gd.Base = grid.Create(20, nodetype.Plain)
	gd.Delta = grid.Create(20, nodetype.NoGround)

	sg.Generate(gd, dbh)

}
