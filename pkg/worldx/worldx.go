package worldx

import (
    "fmt"
    "log"
    "math/rand"
    "strconv"
    "time"
)

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
        return "unknown"
    }
}

type Alien struct {
    Name      string
    Location  *City
    IsTrapped bool
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

// Returns pointer to random connected city or nil if city is isolated (i.e. doesn't have connections)
func (c *City) getRandomConnectedCity() (randomCity *City) {
    maxConnections := len(c.ConnectedCities)
    if maxConnections > 0 {
        move, randomConnection := 0, rand.Intn(maxConnections)
        for _, randomCity = range c.ConnectedCities {
            if move == randomConnection {
                return
            }
            move++
        }
        return
    }

    return nil
}

type WorldX struct {
    Cities map[string]*City  // Maps city name to pointer of respective city
    Aliens map[string]*Alien // Maps alien name to pointer of respective alien
}

func (w *WorldX) String() (wStr string) {
    for _, city := range w.Cities {
        wStr += city.String() + "\n"
    }

    return
}

// Creates and adds city to the world if it doesn't exist yet, returns pointer to city with requested name
func (w *WorldX) CreateCity(cityName string) *City {
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

// Add bidirectional connection between city1 and city2, panics if any of the cities doesn't exist
func (w *WorldX) CreateConnection(city1 *City, city2 *City, dir Direction) {
    if city1 == nil || city2 == nil {
        log.Panicln("CreateConnection: city doesn't exist, cannot create connection if city is nil.")
    }

    city1.ConnectedCities[dir] = city2
    city2.ConnectedCities[dir.GetOpposite()] = city1
}

// Returns pointer to random city without an alien.
// Empty cities slice should contain at least one empty city, otherwise enters an infinite loop
func (w *WorldX) getRandomEmptyCity(emptyCities []string) (randomEmptyCity *City) {
    totalEmptyCities := len(emptyCities)
    for {
        randomEmptyCity = w.Cities[emptyCities[rand.Intn(totalEmptyCities)]]

        if randomEmptyCity.Alien == nil {
            return
        }
    }
}

// Generates aliens one at a time placing them in a random empty city.
// Fails on the tentative to generate more aliens than the number of cities.
func (w *WorldX) GenerateAliens(numberAliens int) {
    if totalAliens, totalCities := len(w.Aliens) + numberAliens, len(w.Cities); totalAliens > totalCities {
        log.Panicf(
            "GenerateAliens: cannot have more aliens in the world than the number of cities! Aliens: %d > Cities: %d",
            totalAliens, totalCities)
    } else if numberAliens <= 0 {
        log.Panicf("GenerateAliens: the number of aliens to be generated need to be positive.")
    }

    emptyCities := make([]string, 0, len(w.Cities))
    for cityName, city := range w.Cities {
        if city.Alien == nil {
            emptyCities = append(emptyCities, cityName)
        }
    }

    rand.Seed(time.Now().UnixNano())
    for alienIndex := 0; alienIndex < numberAliens; alienIndex++ {
        alienName := strconv.Itoa(alienIndex)
        w.CreateAlien(alienName, emptyCities)
    }
}

// If alien doesn't exist, creates it in a random empty city and adds it to the world
// Returns pointer to alien with requested name
func (w *WorldX) CreateAlien(alienName string, possibleEmptyCities []string) *Alien {
    if w.Aliens == nil {
        w.Aliens = make(map[string]*Alien)
    }

    if alien, ok := w.Aliens[alienName]; ok {
        return alien
    } else {
        randomEmptyCity := w.getRandomEmptyCity(possibleEmptyCities)
        isTrapped := len(randomEmptyCity.ConnectedCities) == 0
        newAlien := Alien{alienName, randomEmptyCity, isTrapped}
        w.Aliens[alienName] = &newAlien
        randomEmptyCity.Alien = &newAlien
        return &newAlien
    }
}

// Simulates invasion moving each alien maxIterations times and destroying city if aliens collide.
func (w *WorldX) RunSimulation() {
    const maxIterations int = 10000

    rand.Seed(time.Now().UnixNano())
    for iteration := 0; iteration < maxIterations; iteration++ {
        for _, alien := range w.Aliens {
            w.MoveAlien(alien)
        }
    }
}

// Moves alien from its current city to a random connected city
// If an alien is already present they fight and both the city and aliens are destroyed
func (w *WorldX) MoveAlien(alien *Alien) {
    if alien.IsTrapped {
        return
    }

    nextCity := alien.Location.getRandomConnectedCity()
    if nextCity == nil {
        alien.IsTrapped = true
        return
    } else if nextCity.Alien != nil {
        w.destroyCity(nextCity, alien, nextCity.Alien)
    } else {
        alien.Location.Alien = nil
        alien.Location = nextCity
        nextCity.Alien = alien
        if len(nextCity.ConnectedCities) == 0 {
            alien.IsTrapped = true
        }
    }
}

// Removes connections to the city, and destroys the city and both aliens
func (w *WorldX) destroyCity(city *City, alien1 *Alien, alien2 *Alien) {
    defer fmt.Printf("%s has been destroyed by alien %s and alien %s\n", city.Name, alien1.Name, alien2.Name)
    w.deleteAlien(alien1)
    w.deleteAlien(alien2)
    w.deleteCity(city)
}

func (w *WorldX) deleteCity(city *City) {
    for dir, connectedCity := range city.ConnectedCities {
        delete(connectedCity.ConnectedCities, dir.GetOpposite())
    }

    if city.Alien != nil {
        city.Alien.Location = nil
        city.Alien = nil
    }
    delete(w.Cities, city.Name)
    city = nil
}

func (w *WorldX) deleteAlien(alien *Alien) {
    if alien.Location != nil {
        alien.Location.Alien = nil
        alien.Location = nil
    }
    delete(w.Aliens, alien.Name)
    alien = nil
}
