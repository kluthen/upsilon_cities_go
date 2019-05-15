package producer

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"
	"upsilon_cities_go/lib/cities/item"
	"upsilon_cities_go/lib/cities/storage"
	"upsilon_cities_go/lib/cities/tools"
)

type requirement struct {
	Ressource string
	Type      bool // tell whether require is by type or by name ;)
	Quality   tools.IntRange
	Quantity  int
}

//Producer tell what it produce, within which criteria
type Producer struct {
	ID              int
	Name            string
	ProductName     string
	ProductType     string
	UpgradePoint    int
	BigUpgradePoint int
	Quality         tools.IntRange
	Quantity        tools.IntRange
	BasePrice       int
	Requirements    []requirement
	Delay           int // in cycles
	Level           int // mostly informative, as levels will be applied directly to ranges, requirements and delay
	CurrentXP       int
	NextLevel       int
	Advanced        bool
}

//Production active production stuff ;)
type Production struct {
	ProducerID  int
	StartTime   time.Time
	EndTime     time.Time
	Production  item.Item
	Reservation int64 // storage space reservation ticket

}

//Produce create a new item based on template
func (prod *Producer) produce() (res item.Item) {
	res.Name = prod.ProductName
	res.Type = prod.ProductType
	res.Quality = prod.Quality.Roll()
	res.Quantity = prod.Quantity.Roll()
	res.BasePrice = prod.BasePrice
	return
}

func (rq requirement) String() string {
	return fmt.Sprintf("%d x %s Q[%d-%d]", rq.Quantity, rq.Ressource, rq.Quality.Min, rq.Quality.Max)
}

//Leveling all leveling related action
func (prod *Producer) Leveling(point int) {
	prod.CurrentXP += point
	for prod.CurrentXP >= prod.NextLevel {
		prod.CurrentXP = prod.CurrentXP - prod.NextLevel
		prod.Level++
		prod.UpgradePoint++
		if prod.Level%5 == 0 {
			prod.BigUpgradePoint++
		}
		prod.NextLevel = GetNextLevel(prod.Level)
	}
}

//CanUpgrade Producer can make a simple Upgrade
func (prod *Producer) CanUpgrade() bool {
	return prod.UpgradePoint > 0
}

//CanBigUpgrade Producer can make a big upgrade
func (prod *Producer) CanBigUpgrade() bool {
	return prod.BigUpgradePoint > 0
}

//GetNextLevel return next level needed xp
func GetNextLevel(acLevel int) int {
	return 10 * acLevel
}

//CanProduceShort tell whether it's able to produce item
func CanProduceShort(store *storage.Storage, prod *Producer) (producable bool, err error) {
	// producable immediately ?

	count := 0
	missing := make([]string, 0)
	found := make(map[string]int)
	var missitem string
	for _, v := range prod.Requirements {
		found[v.Ressource] = 0
		count += v.Quantity

		for _, foundling := range store.All(storage.ByTypeOrNameNQuality(v.Ressource, v.Type, v.Quality)) {
			found[v.Ressource] += foundling.Quantity
		}

		if found[v.Ressource] < v.Quantity {
			missitem = fmt.Sprintf("%s need %d have %d", v.String(), v.Quantity, found[v.Ressource])
			missing = append(missing, missitem)
		}
	}

	if len(missing) > 0 {
		return false, fmt.Errorf("not enough ressources: %s", strings.Join(missing, ", "))
	}

	if store.Spaceleft()+count < prod.Quantity.Min {
		return false, fmt.Errorf("not enough space available: potentially got: %d required %d", (store.Spaceleft() + count), prod.Quantity.Min)
	}
	return true, nil
}

//CanProduce tell whether it's able to produce item, if it can produce it relayabely or if not, tell why.
func CanProduce(store *storage.Storage, prod *Producer, ressourcesGenerators map[int]*Producer) (producable bool, nb int, recurrent bool, err error) {
	// producable immediately ?

	count := 0
	missing := make([]string, 0)
	found := make(map[string]int)
	for _, v := range prod.Requirements {
		found[v.Ressource] = 0
		count += v.Quantity

		for _, foundling := range store.All(storage.ByTypeOrNameNQuality(v.Ressource, v.Type, v.Quality)) {
			found[v.Ressource] += foundling.Quantity
		}

		if found[v.Ressource] < v.Quantity {
			missing = append(missing, v.String())
		}
	}

	if len(missing) > 0 {
		return false, 0, false, fmt.Errorf("not enough ressources: %s", strings.Join(missing, ", "))
	}

	if store.Spaceleft()-count < prod.Quantity.Min {
		return false, 0, false, fmt.Errorf("not enough space available: potentially got: %d required %d", (store.Spaceleft() - count), prod.Quantity.Min)
	}

	producable = true
	// we can at least produce one this.

	available := make(map[string]int)
	// check if we can produce more than one

	for _, v := range prod.Requirements {
		available[v.Ressource] = found[v.Ressource] / v.Quantity
	}

	nb = 0

	for _, v := range available {
		nb = tools.Max(nb, v)
	}

	// check if we can produce indefinitely ( due to ressource generator => production exceed requirements )

	found = make(map[string]int)
	available = make(map[string]int)
	for _, v := range prod.Requirements {
		for _, gen := range ressourcesGenerators {
			if gen.ProductType == v.Ressource {
				if gen.Quality.Min >= v.Quality.Min {
					found[v.Ressource] += gen.Quantity.Min
				}
			}
		}
	}

	for _, v := range prod.Requirements {
		if _, has := found[v.Ressource]; has {
			available[v.Ressource] = found[v.Ressource] / v.Quantity
		} else {
			return producable, nb, false, nil
		}
	}

	return producable, nb, true, nil
}

//DeductProducFromStorage attempt to remove necessary items from store to start producer.
func deductProducFromStorage(store *storage.Storage, prod *Producer) error {
	found := make(map[int64]int)
	for _, v := range prod.Requirements {
		target := v.Quantity

		for _, foundling := range store.All(storage.ByTypeOrNameNQuality(v.Ressource, v.Type, v.Quality)) {
			used := tools.Min(target, foundling.Quantity)
			found[foundling.ID] = used
			target -= used

			if target <= 0 {
				break
			}
		}

		if target > 0 {
			return fmt.Errorf("Unable to fit requirement %s", v.String())
		}
	}

	for k, v := range found {
		store.Remove(k, v)
	}
	return nil
}

//Product Kicks in Producer and instantiate a Production, if able.
func Product(store *storage.Storage, prod *Producer, startDate time.Time) (*Production, error) {
	producable, _ := CanProduceShort(store, prod)

	// reserve place for to be coming products...

	if !producable {
		return nil, errors.New("unable to use this Producer")
	}

	production := new(Production)

	production.StartTime = tools.RoundTime(startDate)
	production.EndTime = tools.AddCycles(production.StartTime, prod.Delay)
	production.ProducerID = prod.ID
	production.Production = prod.produce()

	err := deductProducFromStorage(store, prod)

	if err != nil {
		return nil, err
	}

	// from here we already know that space left >= Min production quantity.
	production.Production.Quantity = tools.Min(production.Production.Quantity, store.Spaceleft())
	production.Reservation, err = store.Reserve(production.Production.Quantity)

	if err != nil {
		return nil, err
	}

	return production, nil
}

//IsFinished tell whether production is finished or not ;)
func (prtion *Production) IsFinished(now time.Time) (finished bool) {
	finished = now.After(prtion.EndTime) || now.Equal(prtion.EndTime)
	log.Printf("Producer: %s : End: %s, Now %s, Finished ? %v", prtion.Production.Name, prtion.EndTime.Format(time.RFC3339), now.Format(time.RFC3339), finished)
	return
}

//ProductionCompleted Update store
func ProductionCompleted(store *storage.Storage, prtion *Production, nextUpdate time.Time) error {
	if prtion.IsFinished(nextUpdate) {
		return store.Claim(prtion.Reservation, prtion.Production)
	}
	return errors.New("unable to complete production (not finished)")
}
