package invade

import (
    "bufio"
    "fmt"
    "log"
    "os"
    "strings"

    "github.com/tomasnunes/invasion/pkg/worldx"
)

// Reads map of World X from file, returns new world with the cities and connections described
func ReadWorldMap(file *os.File) (world worldx.WorldX) {
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {

		var newCity *worldx.City
		cityDetails := strings.Fields(scanner.Text())
		for i, value := range cityDetails {

			if i == 0 {
			    newCity = world.CreateCity(value)
			} else if newCity != nil {
				directionDetails := strings.SplitN(value, "=", 2)
				direction := worldx.Direction(directionDetails[0])

				// Ignore directions without city name, with empty city name, and invalid directions
				if len(directionDetails) == 2 && len(directionDetails[1]) > 0 && direction.IsValid() {

				    // Only add connection if it doesn't exist yet, duplicated connections are ignored
				    if _, ok := newCity.ConnectedCities[direction]; !ok {
                        connectedCity := world.CreateCity(directionDetails[1])
                        world.CreateConnection(newCity, connectedCity, direction)
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

// Reads world map from file, generates aliens and runs the simulation of the invasion
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
	world.RunSimulation()

	fmt.Print(world.String())
}
