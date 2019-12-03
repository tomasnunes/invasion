package worldx

import (
    "bufio"
    "fmt"
    "log"
    "math/rand"
    "strconv"
    "strings"
    "time"
)

type WorldX struct {
    cities map[string]*city  // Maps city name to pointer of respective city
    aliens map[string]*alien // Maps alien name to pointer of respective alien
}

func (w *WorldX) String() (wStr string) {
    for _, c := range w.cities {
        wStr += c.String() + "\n"
    }

    return
}

// Reads map of World X from the provided scanner and populates world with the cities and connections described.
// Panics if errors occur while reading the scanner or creating the cities and connections.
func (w *WorldX) ReadWorldMap(scanner *bufio.Scanner) {
    const defaultDirectionSeparator string = "="

    for scanner.Scan() {
        var newCity *city
        cityDetails := strings.Fields(scanner.Text())

        for i, value := range cityDetails {
            if i == 0 {
                newCity = w.createCity(value)
            } else if newCity != nil {
                directionDetails := strings.SplitN(value, defaultDirectionSeparator, 2)
                dir := getDirection(directionDetails[0])

                // Ignore directions without city name, with empty city name, and invalid directions
                if len(directionDetails) == 2 && len(directionDetails[1]) > 0 && dir.isValid() {

                    // Only add connection if it doesn't exist yet, duplicated connections are ignored
                    if newCity.connectedCities[dir] == nil {
                        connectedCity := w.createCity(directionDetails[1])
                        w.createConnection(newCity, connectedCity, dir)
                    }
                }
            }
        }
    }

    if err := scanner.Err(); err != nil {
        log.Panic(err)
    }
}

// Generates aliens one at a time placing them in a random empty city.
// Panics on the tentative to generate more aliens than the number of cities.
func (w *WorldX) GenerateAliens(numberAliens int) {
    if totalAliens, totalCities := len(w.aliens) + numberAliens, len(w.cities); numberAliens <= 0 {
        log.Panicf("GenerateAliens: the number of aliens to be generated need to be positive.")
    } else if totalAliens > totalCities {
        log.Panicf(
            "GenerateAliens: cannot have more aliens in the world than the number of cities! aliens: %d > cities: %d",
            totalAliens, totalCities)
    }

    emptyCities := make([]string, 0, len(w.cities))
    for name, c := range w.cities {
        if c.alien == nil {
            emptyCities = append(emptyCities, name)
        }
    }

    rand.Seed(time.Now().UnixNano())
    for alienIndex := 0; alienIndex < numberAliens; alienIndex++ {
        alienName := strconv.Itoa(alienIndex)
        w.createAlien(alienName, emptyCities)
    }
}

// Simulates invasion moving each alien defaultMaxIterations times and destroying city if aliens collide.
func (w *WorldX) RunSimulation() {
    const defaultMaxIterations int = 10000

    rand.Seed(time.Now().UnixNano())
    for iteration := 0; iteration < defaultMaxIterations; iteration++ {
        for _, a := range w.aliens {
            w.moveAlien(a)
        }
    }
}

// Creates and adds city to the world if it doesn't exist yet, returns pointer to city with requested name
func (w *WorldX) createCity(cityName string) *city {
    if w.cities == nil {
        w.cities = make(map[string]*city)
    }

    if c, ok := w.cities[cityName]; ok {
        return c
    } else {
        newCity := city{
            name:            cityName,
            connectedCities: [maxDirections]*city{},
            alien:           nil,
        }
        w.cities[cityName] = &newCity
        return &newCity
    }
}

// Add bidirectional connection between city1 and city2, panics if any of the cities doesn't exist
func (w *WorldX) createConnection(city1 *city, city2 *city, dir direction) {
    if city1 == nil || city2 == nil {
        log.Panicln("createConnection: city doesn't exist, cannot create connection if city is nil.")
    }

    city1.connectedCities[dir] = city2
    city2.connectedCities[dir.getOpposite()] = city1
}

// Returns pointer to random city without an alien.
// Empty cities slice should contain at least one empty city, otherwise enters an infinite loop
func (w *WorldX) getRandomEmptyCity(emptyCities []string) (randomEmptyCity *city) {
    totalEmptyCities := len(emptyCities)
    for {
        randomEmptyCity = w.cities[emptyCities[rand.Intn(totalEmptyCities)]]

        if randomEmptyCity.alien == nil {
            return
        }
    }
}

// If alien doesn't exist, creates it in a random empty city and adds it to the world
// Returns pointer to alien with requested name
func (w *WorldX) createAlien(alienName string, possibleEmptyCities []string) *alien {
    if w.aliens == nil {
        w.aliens = make(map[string]*alien)
    }

    if a, ok := w.aliens[alienName]; ok {
        return a
    } else {
        randomEmptyCity := w.getRandomEmptyCity(possibleEmptyCities)
        newAlien := alien{
            name:      alienName,
            location:  randomEmptyCity,
            isTrapped: randomEmptyCity.isIsolated(),
        }
        w.aliens[alienName] = &newAlien
        randomEmptyCity.alien = &newAlien
        return &newAlien
    }
}

// Moves alien from its current city to a random connected city,
// if an alien is already present they fight and both the city and aliens are destroyed
func (w *WorldX) moveAlien(alien *alien) {
    if alien.isTrapped {
        return
    }

    nextCity := alien.location.getRandomConnection()
    if nextCity == nil {
        alien.isTrapped = true
        return
    } else if nextCity.alien != nil {
        w.destroyCity(nextCity, alien, nextCity.alien)
    } else {
        alien.location.alien = nil
        alien.location = nextCity
        nextCity.alien = alien
        if nextCity.isIsolated() {
            alien.isTrapped = true
        }
    }
}

// Removes connections to the city, and destroys the city and both aliens
func (w *WorldX) destroyCity(city *city, alien1 *alien, alien2 *alien) {
    defer fmt.Printf("%s has been destroyed by alien %s and alien %s\n", city.name, alien1.name, alien2.name)
    w.deleteAlien(alien1)
    w.deleteAlien(alien2)
    w.deleteCity(city)
}

func (w *WorldX) deleteCity(city *city) {
    for dir, connection := range city.connectedCities {
        if connection != nil {
            connection.connectedCities[direction(dir).getOpposite()] = nil
        }
    }

    if city.alien != nil {
        city.alien.location = nil
        city.alien = nil
    }
    delete(w.cities, city.name)
    city = nil
}

func (w *WorldX) deleteAlien(alien *alien) {
    if alien.location != nil {
        alien.location.alien = nil
        alien.location = nil
    }
    delete(w.aliens, alien.name)
    alien = nil
}

type alien struct {
    name      string
    location  *city
    isTrapped bool
}

type city struct {
    name            string
    connectedCities [maxDirections]*city
    alien           *alien
}

func (c *city) String() (cStr string) {
    cStr = c.name
    for dir, connection := range c.connectedCities {
        if connection != nil {
            cStr += " " + direction(dir).String() + "=" + connection.name
        }
    }
    return
}

func (c *city) isIsolated() bool {
    for _, connection := range c.connectedCities {
        if connection != nil {
            return false
        }
    }
    return true
}

// Returns pointer to random connected city or nil if city is isolated (i.e. doesn't have connections)
func (c *city) getRandomConnection() (randomCity *city) {
    if c.isIsolated() {
        return nil
    }

    for {
        if randomCity = c.connectedCities[rand.Intn(int(maxDirections))]; randomCity != nil {
            return
        }
    }
}

type direction int

const (
    north direction = iota
    south
    east
    west
    maxDirections
    unknown         = -1
)

func getDirection(dir string) direction {
    switch dir {
    case "north":
        return north
    case "south":
        return south
    case "east":
        return east
    case "west":
        return west
    default:
        return unknown
    }
}

func (d direction) String() string {
    if d < 0 || d >= maxDirections {
        return "unknown"
    }
    return [...]string{"north", "south", "east", "west"}[d]
}

func (d direction) isValid() bool {
    return d == north || d == south || d == east || d == west
}

func (d direction) getOpposite() direction {
    switch d {
    case north:
        return south
    case south:
        return north
    case east:
        return west
    case west:
        return east
    default:
        return unknown
    }
}
