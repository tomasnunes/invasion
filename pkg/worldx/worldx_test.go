package worldx

import (
    "bufio"
    "bytes"
    "strconv"
    "strings"
    "testing"

    "github.com/tomasnunes/invasion/pkg/worldx"
)

func TestReadWorldMap(t *testing.T) {
    const inputWorldMap =
`
Zuu south=Zoo
Zoo north=Zuu east=D'Foo
D'Foo east=Baz south=Bar
Bar north=D'Foo south=Guu east=Foo
Guu north=Bar east=Fuu south=P=NP
P=NP north=Guu
Baz west=D'Foo south=Foo east=Qu-ux
Foo north=Baz west=Bar south=Fuu east=Goo
Fuu north=Foo west=Guu
Goo north=Qu-ux west=Foo
Alone
`

    // Should expect city 'Qu-ux' since a connection of 'Goo' mention it.
    // Should read 'Alone' even without any connections.
    var expectedCities = map[string][worldx.MaxDirections]string{
        "Zuu":   {                       worldx.South: "Zoo"},
        "Zoo":   {worldx.North: "Zuu",                         worldx.East: "D'Foo"},
        "D'Foo": {                       worldx.South: "Bar",  worldx.East: "Baz",  worldx.West: "Zoo"},
        "Bar":   {worldx.North: "D'Foo", worldx.South: "Guu",  worldx.East: "Foo"},
        "Guu":   {worldx.North: "Bar",   worldx.South: "P=NP", worldx.East: "Fuu"},
        "P=NP":  {worldx.North: "Guu"},
        "Baz":   {                       worldx.South: "Foo",  worldx.East: "Qu-ux", worldx.West: "D'Foo"},
        "Foo":   {worldx.North: "Baz",   worldx.South: "Fuu",  worldx.East: "Goo",   worldx.West: "Bar"},
        "Fuu":   {worldx.North: "Foo",                                               worldx.West: "Guu"},
        "Goo":   {worldx.North: "Qu-ux",                                             worldx.West: "Foo"},
        "Qu-ux": {                       worldx.South: "Goo",                        worldx.West: "Baz"},
        "Alone": {},
    }

    scannerWithWorldMap := bufio.NewScanner(strings.NewReader(inputWorldMap))

    actualWorld := worldx.WorldX{}
    actualWorld.ReadWorldMap(scannerWithWorldMap)

    if totalActualCities, totalExpectedCities := len(actualWorld.Cities), len(expectedCities);
        totalActualCities != totalExpectedCities {
        t.Errorf("The amount of cities created doesn't equal the expected amount: expected: %d != actual: %d",
            totalExpectedCities, totalActualCities)
    }

    for expectedName, expectedConnections := range expectedCities {
        if actualCity, ok := actualWorld.Cities[expectedName]; !ok {
            t.Errorf("Unable to find expected city: %s", expectedName)
        } else {
            for dir := worldx.Direction(0); dir < worldx.MaxDirections; dir++ {
                actualConnection := actualCity.Connection(dir)
                if actualConnection == nil && expectedConnections[dir] != "" {
                    t.Errorf("Expected connection: %s %s=%s", expectedName, dir.String(), expectedConnections[dir])
                } else if actualConnection != nil && actualConnection.Name() != expectedConnections[dir] {
                    t.Errorf("Wrong connection (if empty name in connection, no connection was expected): expected: %s %s=%s != actual: %s %s=%s",
                        expectedName, dir.String(), expectedConnections[dir],
                        actualCity.Name(), dir.String(), actualConnection.Name())
                }
            }
        }
    }
}

func getGenerateAliensTestWorld(numberCities int) (testWorld worldx.WorldX) {
    for i := 0; i < numberCities; i++ {
        cityName := strconv.Itoa(i)
        testWorld.CreateCity(cityName)
    }

    return
}

func TestGenerateAliens(t *testing.T) {
    const maxAmountAliensToTest = 4

    var generateAliensTests = []struct{
        numberAliens         int // input
        expectedAmountAliens int // expected amount of aliens created
    }{
        {0, 0},
        {1, 1},
        {maxAmountAliensToTest, maxAmountAliensToTest},
    }

    for _, test := range generateAliensTests {
        testWorld := getGenerateAliensTestWorld(maxAmountAliensToTest)
        testWorld.GenerateAliens(test.numberAliens)

        if actualAmountAliens := len(testWorld.Aliens); actualAmountAliens != test.expectedAmountAliens {
            t.Errorf("Wrong amount of aliens generated: expected %d != actual: %d",
                test.expectedAmountAliens, actualAmountAliens)
        }

        for i := 0; i < test.expectedAmountAliens; i++ {
            expectedName := strconv.Itoa(i)
            if actualAlien, ok := testWorld.Aliens[expectedName]; !ok || actualAlien.Name() != expectedName {
                t.Errorf("Expected alien %s was not generated", expectedName)
            } else if actualAlien.Location() == nil {
                t.Error("Alien was not placed in any city")
            } else if actualAlien.Location().Alien() != actualAlien {
                t.Error("Alien pointer in the city not updated")
            } else if !actualAlien.IsTrapped() {
                t.Error("Alien should be trapped when placed in a city without connections")
            }
        }
    }
}

func TestGenerateAliensWithMoreAliensThanCities(t *testing.T) {
    const numberAliens = 5
    testWorld := getGenerateAliensTestWorld(numberAliens - 1)

    defer func() {
        if r := recover(); r == nil {
            t.Errorf("Expected panic when trying to generate more aliens (%d) than cities (%d).",
                numberAliens, len(testWorld.Cities))
        }
    }()

    testWorld.GenerateAliens(numberAliens)
}

func getRunSimulationTestWorld(numberCities int, numberAliens int) (testWorld worldx.WorldX) {
    cityNames := make([]string, numberCities)
    for i := 0; i < numberCities; i++ {
        cityName := strconv.Itoa(i)
        cityNames[i] = cityName
        testWorld.CreateCity(cityName)
    }

    for i := 0; i < numberAliens; i++ {
        alienName := strconv.Itoa(i)
        testWorld.CreateAlien(alienName, cityNames)
    }

    return
}

func TestRunSimulationWithTrappedAliens(t *testing.T) {
    const numberCities = 2
    const numberAliens = 2

    testWorld := getRunSimulationTestWorld(numberCities, numberAliens)

    buf := new(bytes.Buffer)
    writer := bufio.NewWriter(buf)

    testWorld.RunSimulation(writer)

    if len(testWorld.Cities) != numberCities || len(testWorld.Aliens) != numberAliens {
        t.Error("When all aliens are trapped neither the amount of cities or aliens should change")
    }
}

func TestRunSimulationWhenAliensFight(t *testing.T) {
    const numberCities = 2
    const numberAliens = 2

    testWorld := getRunSimulationTestWorld(numberCities, numberAliens)
    testWorld.AddConnection(testWorld.Cities["0"], testWorld.Cities["1"], worldx.North)

    buf := new(bytes.Buffer)
    writer := bufio.NewWriter(buf)

    testWorld.RunSimulation(writer)

    destroyMessage := buf.String()
    if destroyMessage != "0 has been destroyed by alien 0 and alien 1\n" &&
        destroyMessage != "0 has been destroyed by alien 1 and alien 0\n" &&
        destroyMessage != "1 has been destroyed by alien 0 and alien 1\n" &&
        destroyMessage != "1 has been destroyed by alien 1 and alien 0\n" {
        t.Errorf("Wrong message when destroying city:%s", destroyMessage)
    }

    if len(testWorld.Cities) != 1 || len(testWorld.Aliens) != 0 {
        t.Error("When aliens meet and fight one city should be destroyed and both aliens die")
    }

    var remainingCity *worldx.City
    var remainingCityName string
    for name, city := range testWorld.Cities {
        remainingCityName = name
        remainingCity = city
    }

    if remainingCityName != "0" && remainingCityName != "1" {
        t.Errorf("Name of the cities shouldn't change during simulation, remaining city name: %s", remainingCityName)
    } else if remainingCity == nil {
        t.Error("Only one city should be destroyed when two aliens fight")
    } else if !remainingCity.IsIsolated() {
        t.Error("Connections should be removed when a city is destroyed")
    } else if remainingCity.Alien() != nil {
        t.Error("City field with pointer to alien should be set to <nil> when alien is destroyed")
    }
}

func TestDirection(t *testing.T) {
    var directionTests = []struct{
        dir              worldx.Direction // input
        expectedValid    bool             // expected IsValid() result
        expectedString   string           // expected String() result
        expectedOpposite worldx.Direction // expected GetOpposite() result
    }{
        {worldx.North,            true,  "north",   worldx.South},
        {worldx.South,            true , "south",   worldx.North},
        {worldx.East,             true , "east",    worldx.West},
        {worldx.West,             true , "west",    worldx.East},
        {worldx.UnknownDirection, false, "unknown", worldx.UnknownDirection},
        {worldx.MaxDirections,    false, "unknown", worldx.UnknownDirection},
        {worldx.Direction(-10),   false, "unknown", worldx.UnknownDirection},
    }

    for _, test := range directionTests {
        actualValid := test.dir.IsValid()
        if actualValid != test.expectedValid {
            t.Errorf("%v.IsValid(): expected %v, actual %v", test.dir, test.expectedValid, actualValid)
        }

        actualString := test.dir.String()
        if actualString != test.expectedString {
            t.Errorf("%v.String(): expected %s, actual %s", test.dir, test.expectedString, actualString)
        }

        actualOpposite := test.dir.GetOpposite()
        if actualOpposite != test.expectedOpposite {
            t.Errorf("%v.GetOpposite(): expected %v, actual %v", test.dir, test.expectedOpposite, actualOpposite)
        }
    }
}
