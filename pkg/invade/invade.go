package invade

import (
    "bufio"
    "fmt"
    "log"
    "math/rand"
    "os"
    "strconv"
    "strings"
    "sync"
    "time"
)

type WorldX struct {
	Cities map[string]*City  // Maps city name to pointer of respective city
	Aliens map[string]*Alien // Maps alien name to pointer of respective alien
	mux    sync.RWMutex
}

func (w *WorldX) String() (wStr string) {
	w.mux.RLock()
	defer w.mux.RUnlock()

	for _, city := range w.Cities {
		wStr += city.String() + "\n"
	}

	return
}

// If city doesn't exist yet, creates and adds it to the world
func (w *WorldX) AddCity(cityName string) {
	w.mux.Lock()
	w.getCity(cityName)
	w.mux.Unlock()
}

// Adds connection between city1 and city2, if any of the cities doesn't exist, it's created and added to the world
func (w *WorldX) AddCityConnection(cityName1 string, cityName2 string, dir Direction) {
	w.mux.Lock()

	city1, city2 := w.getCity(cityName1), w.getCity(cityName2)
	city1.ConnectedCities[dir] = city2
	city2.ConnectedCities[dir.GetOpposite()] = city1

	w.mux.Unlock()
}

func (w *WorldX) getCity(cityName string) *City {
	if w.Cities == nil {
		w.Cities = make(map[string]*City)
	}

	if city, ok := w.Cities[cityName]; ok {
		return city
	} else {
		newCity := City{cityName, map[Direction]*City{}, nil}
		w.Cities[cityName] = &newCity
		return &newCity
	}
}

// Generates aliens and places them in a random empty city.
// Fails on the tentative to generate more aliens than the number of cities.
func (w *WorldX) GenerateAliens(numberAliens int) {
    w.mux.Lock()
    defer w.mux.Unlock()

    if totalAliens, totalCities := len(w.Aliens) + numberAliens, len(w.Cities); totalAliens > totalCities {
        log.Panicf(
            "GenerateAliens: cannot have more aliens in the world than the number of cities! Aliens: %d > Cities: %d",
            totalAliens, totalCities)
    } else if numberAliens <= 0 {
        log.Panicf("GenerateAliens: the number of aliens to be generated need to be positive.")
    }

    if w.Aliens == nil {
        w.Aliens = make(map[string]*Alien)
    }

    // Get slice with the names of empty cities
    emptyCities := make([]string, 0, len(w.Cities))
    for cityName, city := range w.Cities {
        if city.Alien == nil {
            emptyCities = append(emptyCities, cityName)
        }
    }
    totalEmptyCities := len(emptyCities)

    rand.Seed(time.Now().UnixNano())
    for alienIndex := 0; alienIndex < numberAliens; alienIndex++ {
        randomEmptyCity := emptyCities[rand.Intn(totalEmptyCities)]
        for w.Cities[randomEmptyCity].Alien != nil {
            randomEmptyCity = emptyCities[rand.Intn(totalEmptyCities)]
        }

        alienName := strconv.Itoa(alienIndex)
        newAlien := Alien{alienName, w.Cities[randomEmptyCity]}
        w.Aliens[alienName] = &newAlien
        w.Cities[randomEmptyCity].Alien = &newAlien

        fmt.Println("Alien", w.Aliens[alienName].Name, "placed in city", w.Aliens[alienName].Location.Name)
    }
}

func (w *WorldX) RunSimulation() {
    panic("RunSimulation is not yet implemented!")
}

type Direction string

const (
	North Direction = "north"
	South           = "south"
	East            = "east"
	West            = "west"
)

func (d Direction) String() string {
	return string(d)
}

func (d Direction) IsValid() bool {
	return d == North || d == South || d == East || d == West
}

func (d Direction) GetOpposite() Direction {
	switch d {
	case North :
		return South
	case South :
		return North
	case East :
		return West
	case West :
		return East
	default :
		return ""
	}
}

type City struct {
	Name            string
	ConnectedCities map[Direction]*City
	Alien           *Alien
}

func (c *City) String() (cStr string) {
	cStr = c.Name
	for direction, city := range c.ConnectedCities {
		cStr += " " + direction.String() + "=" + city.Name
	}

	return
}

type Alien struct {
	Name     string
	Location *City
}

func ReadWorldMap(file *os.File) (world WorldX) {
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {

		var newCity string
		cityDetails := strings.Fields(scanner.Text())
		for i, value := range cityDetails {

			if i == 0 {
			    newCity = value
                world.AddCity(newCity)
			} else {
				directionDetails := strings.SplitN(value, "=", 2)
				direction := Direction(directionDetails[0])

				// Ignore directions without city name, with empty city name, and invalid directions
				if len(directionDetails) == 2 && len(directionDetails[1]) > 0 && direction.IsValid() {

				    // Only add connection if it doesn't exist yet, duplicated connections are ignored
				    if _, ok := world.Cities[newCity].ConnectedCities[direction]; !ok {
                        connectedCity := directionDetails[1]
                        world.AddCityConnection(newCity, connectedCity, direction)
                    }
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Panic(err)
	}

	return world
}

func Invade(filename string, numberAliens int) {
	file, err := os.Open(filename)
	if err != nil {
		log.Panic(err)
	}
	defer func() {
		if err = file.Close(); err != nil {
			log.Panic(err)
		}
	}()

	world := ReadWorldMap(file)

	world.GenerateAliens(numberAliens)
	//world.RunSimulation()

	fmt.Print(world.String())
}
