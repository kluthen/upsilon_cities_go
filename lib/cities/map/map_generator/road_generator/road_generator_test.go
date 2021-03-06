package road_generator

import (
	"log"
	_ "net/http/pprof"
	"testing"
	"upsilon_cities_go/lib/cities/city/resource"
	"upsilon_cities_go/lib/cities/city/resource_generator"
	"upsilon_cities_go/lib/cities/map/grid"
	"upsilon_cities_go/lib/cities/map/map_generator/city_generator"
	"upsilon_cities_go/lib/cities/map/map_generator/forest_generator"
	"upsilon_cities_go/lib/cities/map/map_generator/mountain_generator"
	"upsilon_cities_go/lib/cities/map/map_generator/sea_generator"
	"upsilon_cities_go/lib/cities/nodetype"
	"upsilon_cities_go/lib/db"
	"upsilon_cities_go/lib/misc/config/system"
	"upsilon_cities_go/lib/misc/generator"
)

func TestRoadGenerator(t *testing.T) {

	system.LoadConf()
	dbh := db.NewTest()
	db.FlushDatabase(dbh)
	generator.Load()

	dg := city_generator.Create()
	dg.Density.Min = 3
	dg.Density.Max = 3
	gd := new(grid.CompoundedGrid)
	gd.Base = grid.Create(20, nodetype.Plain)
	gd.Delta = grid.Create(20, nodetype.NoGround)

	gd.Base.Insert(dbh)

	for idx := range gd.Base.Nodes {
		gd.Base.Nodes[idx].Activated = []resource.Resource{resource_generator.MustOne("Fer")}
	}

	dg.Generate(gd, dbh)

	//cities are on a layer below, so compact and reset delta layer.
	gd.Base = gd.Compact()
	gd.Delta = grid.Create(20, nodetype.NoGround)

	rg := Create()
	rg.Generate(gd, dbh)
	gd.Base = gd.Compact()

	log.Printf("Delta: \n%s", gd.Delta.String())
	log.Printf("Result: \n%s", gd.Base.String())

	t.Error("Not implemented")
}

func TestRoadGeneratorComplex(t *testing.T) {

	system.LoadConf()
	dbh := db.NewTest()
	db.FlushDatabase(dbh)
	generator.Load()

	dg := city_generator.Create()
	gd := new(grid.CompoundedGrid)
	gd.Base = grid.Create(20, nodetype.Plain)
	gd.Delta = grid.Create(20, nodetype.NoGround)
	gd.Base.Insert(dbh)

	for idx := range gd.Base.Nodes {
		gd.Base.Nodes[idx].Activated = []resource.Resource{resource_generator.MustOne("Fer")}
	}

	mg := mountain_generator.Create()
	mg.Generate(gd, dbh)
	gd.Base = gd.Compact()

	sg := sea_generator.Create()
	sg.Generate(gd, dbh)
	gd.Base = gd.Compact()
	fg := forest_generator.Create()
	fg.Generate(gd, dbh)
	gd.Base = gd.Compact()
	dg.Generate(gd, dbh)

	//cities are on a layer below, so compact and reset delta layer.
	gd.Base = gd.Compact()

	rrg := Create()
	err := rrg.Generate(gd, dbh)
	if err != nil {
		t.Errorf("Failed to generate road map: %s", err)
		log.Printf("Delta: \n%s", gd.Delta.String())
		gd.Base = gd.Compact()
		log.Printf("Result: \n%s", gd.Base.String())
		return
	}
	gd.Base = gd.Compact()

	log.Printf("Delta: \n%s", gd.Delta.String())
	log.Printf("Result: \n%s", gd.Base.String())

	t.Error("Not implemented")
}
