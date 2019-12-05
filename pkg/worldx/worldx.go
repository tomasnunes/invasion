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
    Cities map[string]*City  // Maps city name to pointer of respective city
    Aliens map[string]*Alien // Maps alien name to pointer of respective alien
}

func (w *WorldX) String() (wStr string) {
    for _, c := range w.Cities {
        wStr += c.String() + "\n"
    }

    return
}

// Reads map of World X from the provided scanner and populates world with the cities and connections described.
// Panics if errors occur while reading the scanner or creating the cities and connections.
func (w *WorldX) ReadWorldMap(scanner *bufio.Scanner) {
    const defaultDirectionSeparator string = "="

    for scanner.Scan() {
        var newCity *City
        cityDetails := strings.Fields(scanner.Text())

        for i, value := range cityDetails {
            if i == 0 {
                newCity = w.CreateCity(value)
            } else if newCity != nil {
                directionDetails := strings.SplitN(value, defaultDirectionSeparator, 2)
                dir := GetDirection(directionDetails[0])

                // Ignore directions without city name, with empty city name, and invalid directions
                if len(directionDetails) == 2 && len(directionDetails[1]) > 0 && dir.IsValid() {

                    // Only add connection if it doesn't exist yet, duplicated connections are ignored
                    if newCity.connectedCities[dir] == nil {
                        connectedCity := w.CreateCity(directionDetails[1])
                        w.AddConnection(newCity, connectedCity, dir)
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
    if totalAliens, totalCities := len(w.Aliens) + numberAliens, len(w.Cities); numberAliens < 0 {
        log.Panicf("GenerateAliens: the number of aliens to be generated need to be positive.")
    } else if totalAliens > totalCities {
        log.Panicf(
            "GenerateAliens: cannot have more aliens in the world than the number of cities! aliens: %d > cities: %d",
            totalAliens, totalCities)
    } else if numberAliens == 0 {
        return
    }

    emptyCities := make([]string, 0, len(w.Cities))
    for name, c := range w.Cities {
        if c.alien == nil {
            emptyCities = append(emptyCities, name)
        }
    }

    rand.Seed(time.Now().UnixNano())
    for alienIndex := 0; alienIndex < numberAliens; alienIndex++ {
        alienName := strconv.Itoa(alienIndex)
        w.CreateAlien(alienName, emptyCities)
    }
}

// Simulates invasion moving each alien `defaultMaxIterations` times or until it's trapped in an isolated city.
// When two aliens meet in the same city they fight and in the process, both aliens die and the city is destroyed
// severing all its connections. Prints message to the writer for every city destroyed.
func (w *WorldX) RunSimulation(writer *bufio.Writer) {
    const defaultMaxIterations int = 10000

    rand.Seed(time.Now().UnixNano())
    for iteration := 0; iteration < defaultMaxIterations; iteration++ {
        for _, a := range w.Aliens {
            w.moveAlien(a, writer)
        }
    }

    if err := writer.Flush(); err != nil {
        log.Panic(err)
    }
}

// Creates and adds city to the world if it doesn't exist yet, returns pointer to city with requested name.
func (w *WorldX) CreateCity(cityName string) *City {
    if w.Cities == nil {
        w.Cities = make(map[string]*City)
    }

    if c, ok := w.Cities[cityName]; ok {
        return c
    } else {
        newCity := City{
            name:            cityName,
            connectedCities: [MaxDirections]*City{},
            alien:           nil,
        }
        w.Cities[cityName] = &newCity
        return &newCity
    }
}

// Add bidirectional connection between city1 and city2.
// Panics if any of the cities doesn't exist or the direction isn't valid.
func (w *WorldX) AddConnection(city1 *City, city2 *City, dir Direction) {
    if city1 == nil || city2 == nil {
        log.Panicln("AddConnection: city doesn't exist, cannot create connection if either city is <nil>.")
    } else if !dir.IsValid() {
        log.Printf("AddConnection: invalid direction %v.", dir)
    }

    city1.connectedCities[dir] = city2
    city2.connectedCities[dir.GetOpposite()] = city1

    if city1.alien != nil && city1.alien.isTrapped {
        city1.alien.isTrapped = false
    }
    if city2.alien != nil && city2.alien.isTrapped {
        city2.alien.isTrapped = false
    }
}

// Returns pointer to random city without an alien.
// Empty cities slice should contain at least one empty city, otherwise enters an infinite loop.
func (w *WorldX) getRandomEmptyCity(emptyCities []string) (randomEmptyCity *City) {
    totalEmptyCities := len(emptyCities)
    for {
        randomEmptyCity = w.Cities[emptyCities[rand.Intn(totalEmptyCities)]]

        if randomEmptyCity.alien == nil {
            return
        }
    }
}

// If alien doesn't exist, creates it in a random empty city and adds it to the world.
// Returns pointer to alien with requested name.
func (w *WorldX) CreateAlien(alienName string, possibleEmptyCities []string) *Alien {
    if w.Aliens == nil {
        w.Aliens = make(map[string]*Alien)
    }

    if a, ok := w.Aliens[alienName]; ok {
        return a
    } else {
        randomEmptyCity := w.getRandomEmptyCity(possibleEmptyCities)
        newAlien := Alien{
            name:      alienName,
            location:  randomEmptyCity,
            isTrapped: randomEmptyCity.IsIsolated(),
        }
        w.Aliens[alienName] = &newAlien
        randomEmptyCity.alien = &newAlien
        return &newAlien
    }
}

// Moves alien from its current city to a random connected city if he isn't trapped,
// if an alien is already present they fight and the city and both aliens are destroyed.
func (w *WorldX) moveAlien(alien *Alien, writer *bufio.Writer) {
    if alien.isTrapped {
        return
    }

    nextCity := alien.location.getRandomConnection()
    if nextCity == nil {
        alien.isTrapped = true
        return
    } else if nextCity.alien != nil {
        cityName, alien1Name, alien2Name := nextCity.name, alien.name, nextCity.alien.name
        defer func() {
            _, err := fmt.Fprintf(writer, "%s has been destroyed by alien %s and alien %s\n",
                cityName, alien1Name, alien2Name)
            if err != nil {
                log.Panic(err)
            }
        }()

        w.destroyCity(nextCity, alien, nextCity.alien)
    } else {
        alien.location.alien = nil
        alien.location = nextCity
        nextCity.alien = alien
        if nextCity.IsIsolated() {
            alien.isTrapped = true
        }
    }
}

// Removes connections to the city, and destroys the city and both aliens.
func (w *WorldX) destroyCity(city *City, alien1 *Alien, alien2 *Alien) {
    w.deleteAlien(alien1)
    w.deleteAlien(alien2)
    w.deleteCity(city)
}

func (w *WorldX) deleteCity(city *City) {
    if city == nil {
        return
    }

    for dir, connection := range city.connectedCities {
        if connection != nil {
            connection.connectedCities[Direction(dir).GetOpposite()] = nil
        }
    }

    if city.alien != nil {
        city.alien.location = nil
        city.alien = nil
    }
    delete(w.Cities, city.name)
    city = nil
}

func (w *WorldX) deleteAlien(alien *Alien) {
    if alien == nil {
        return
    }

    if alien.location != nil {
        alien.location.alien = nil
        alien.location = nil
    }
    delete(w.Aliens, alien.name)
    alien = nil
}

type Alien struct {
    name      string
    location  *City
    isTrapped bool
}

func (a *Alien) Name() string {
    return a.name
}

func (a *Alien) Location() *City {
    return a.location
}

func (a *Alien) IsTrapped() bool {
    return a.isTrapped
}

type City struct {
    name            string
    connectedCities [MaxDirections]*City
    alien           *Alien
}

func (c *City) Name() string {
    return c.name
}

func (c *City) Connection(dir Direction) *City {
    return c.connectedCities[dir]
}

func (c *City) Alien() *Alien {
    return c.alien
}

func (c *City) String() (cStr string) {
    cStr = c.name
    for dir, connection := range c.connectedCities {
        if connection != nil {
            cStr += " " + Direction(dir).String() + "=" + connection.name
        }
    }
    return
}

func (c *City) IsIsolated() bool {
    for _, connection := range c.connectedCities {
        if connection != nil {
            return false
        }
    }
    return true
}

// Returns pointer to random connected city or nil if city is isolated (i.e. doesn't have connections).
func (c *City) getRandomConnection() (randomCity *City) {
    if c.IsIsolated() {
        return nil
    }

    for {
        if randomCity = c.connectedCities[rand.Intn(int(MaxDirections))]; randomCity != nil {
            return
        }
    }
}

type Direction int

const (
    North Direction  = iota
    South
    East
    West
    MaxDirections
    UnknownDirection = -1
)

func GetDirection(dir string) Direction {
    switch dir {
    case "north":
        return North
    case "south":
        return South
    case "east":
        return East
    case "west":
        return West
    default:
        return UnknownDirection
    }
}

func (d Direction) String() string {
    if d < 0 || d >= MaxDirections {
        return "unknown"
    }
    return [...]string{"north", "south", "east", "west"}[d]
}

func (d Direction) IsValid() bool {
    return d == North || d == South || d == East || d == West
}

func (d Direction) GetOpposite() Direction {
    switch d {
    case North:
        return South
    case South:
        return North
    case East:
        return West
    case West:
        return East
    default:
        return UnknownDirection
    }
}
