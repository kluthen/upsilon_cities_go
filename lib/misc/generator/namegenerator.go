package generator

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"upsilon_cities_go/lib/misc/config/system"
)

var nameList map[string]*WordPart

//WordType use to generate a name depending of surrounding
type WordType struct {
	Sea      []string
	Mountain []string
	Forest   []string
	Neutral  []string
	Special  []string
}

//WordPart use to generate a full name
type WordPart struct {
	Body       WordType
	Prefix     WordType
	Suffix     WordType
	BodySuffix WordType
}

// CreateSampleFile does what it says
func CreateSampleFile() {

	word := new(WordPart)
	word.Body.Sea = []string{"_", "_", "_"}
	word.Body.Mountain = []string{"_", "_", "_"}
	word.Body.Forest = []string{"_", "_", "_"}
	word.Body.Neutral = []string{"_", "_", "_"}
	word.Body.Special = []string{"_", "_", "_"}

	word.Prefix.Sea = []string{"_", "_", "_"}
	word.Prefix.Mountain = []string{"_", "_", "_"}
	word.Prefix.Forest = []string{"_", "_", "_"}
	word.Prefix.Neutral = []string{"_", "_", "_"}
	word.Prefix.Special = []string{"_", "_", "_"}

	word.Suffix.Sea = []string{"_", "_", "_"}
	word.Suffix.Mountain = []string{"_", "_", "_"}
	word.Suffix.Forest = []string{"_", "_", "_"}
	word.Suffix.Neutral = []string{"_", "_", "_"}
	word.Suffix.Special = []string{"_", "_", "_"}

	word.BodySuffix.Sea = []string{"_", "_", "_"}
	word.BodySuffix.Mountain = []string{"_", "_", "_"}
	word.BodySuffix.Forest = []string{"_", "_", "_"}
	word.BodySuffix.Neutral = []string{"_", "_", "_"}
	word.BodySuffix.Special = []string{"_", "_", "_"}

	bytes, _ := json.MarshalIndent(word, "", "\t")
	ioutil.WriteFile(fmt.Sprintf("%s/%s", system.Get("data_names", "data/names"), "sample.json.sample"), bytes, 0644)
}

//Load prepare the whole list for later use ;)
func Load() {
	nameList = make(map[string]*WordPart)

	filepath.Walk(system.MakePath(system.Get("data_names", "data/names")), func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatalf("Name Generator: prevent panic by handling failure accessing a path %q: %v\n", system.MakePath(system.Get("data_names", "data/names")), err)
			return err
		}

		if strings.HasSuffix(info.Name(), ".json") {
			f, ferr := os.Open(path)
			if ferr != nil {
				log.Fatalln("Name Generator: No Name data file present")
			}

			nameJSON, ferr := ioutil.ReadAll(f)
			if ferr != nil {
				log.Fatalln("Name Generator: Data file found but unable to read it all.")
			}

			f.Close()

			names := new(WordPart)
			json.Unmarshal(nameJSON, &names)

			filename := strings.TrimSuffix(info.Name(), filepath.Ext(info.Name()))
			log.Printf("Name Generator: Adding %s to name generator", filename)
			nameList[filename] = names
		}

		return nil
	})

	log.Printf("Name Generator: Loaded %d file(s)", len(nameList))

}

//CityName Generate a new city name
func CityName() string {

	boolPre := rand.Int31n(100) <= 5
	boolSuf := rand.Int31n(100) <= 5

	bodyList := nameList["city"].Body.Neutral
	prefixList := nameList["city"].Prefix.Neutral
	suffixList := nameList["city"].Suffix.Neutral

	name := bodyList[rand.Intn((len(bodyList) - 1))]

	if boolPre {
		prefix := prefixList[rand.Intn((len(prefixList) - 1))]
		name = fmt.Sprintf("%s-%s", prefix, name)
	}

	if boolSuf {
		name = fmt.Sprintf("%s-%s", name, suffixList[rand.Intn(len(suffixList)-1)])
	}

	return name
}

//RegionName Generate a new region name
func RegionName(regionType string) string {

	bodyList := nameList["region"].Body.Neutral
	name := bodyList[rand.Intn((len(bodyList) - 1))]

	return name
}

//CorpName Generate a new corporation name
func CorpName(exeptionList []string) string {

	bodyList := nameList["corp"].Body.Neutral

	name := "GoldenCorp"
	isAvailable := true

	for i := 0; i < 15; i++ {
		name = bodyList[rand.Intn((len(bodyList) - 1))]
		for _, b := range exeptionList {
			if b == name {
				isAvailable = false
			}
		}

		if isAvailable {
			break
		}
	}

	return name
}
