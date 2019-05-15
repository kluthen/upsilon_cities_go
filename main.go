package main

import (
	"log"
	"math/rand"
	"time"
	"upsilon_cities_go/lib/cities/city/producer"
	"upsilon_cities_go/lib/cities/city_manager"
	"upsilon_cities_go/lib/cities/grid_manager"
	"upsilon_cities_go/lib/cities/tools"
	"upsilon_cities_go/lib/db"
	"upsilon_cities_go/lib/misc/generator"
	"upsilon_cities_go/web"
	"upsilon_cities_go/web/templates"
	"upsilon_cities_go/web/tools"
)

func main() {
	rand.Seed(time.Now().Unix())

	tools.InitCycle()
	// ensure that in memory storage is fine.
	city_manager.InitManager()
	grid_manager.InitManager()

	generator.CreateSampleFile()
	generator.Init()

	producer.CreateSampleFile()
	producer.Load()

	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile)

	handler := db.New()
	db.CheckVersion(handler)
	handler.Close()

	router := web.RouterSetup()
	tools.SetRouter(router)
	templates.LoadTemplates()
	web.ListenAndServe(router)

}
