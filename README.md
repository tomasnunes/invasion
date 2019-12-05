# Simulate Invasion

This command-line application reads and constructs a World X, generates aliens, allocates them to an empty city,
and simulates an invasion. During the simulation, aliens move between cities at random, when two aliens meet in
the same city they fight and in the process both aliens die and the city is destroyed severing all its connections.
The final state of the world is printed at the end.

### Input FILE Fromat (World Map Description)

```text
<new city name> [north=<connected city name> south=<...> east=<...> west=<...>]
...
```

- City names cannot contain spaces, any other character is allowed.
- No effort is made to verify if a `<new city name>` is empty and therefore invalid.
- Any directional connection is optional, a city doesn't need to connect to other cities in all directions.

### Packages

- `worldx`

    This package contains the `WorldX`, `City`, and `Alien` types and possible interactions with the world to run
    a full simulation, the most relevant exported functions are described below:
    - `ReadWorldMap()` → Reads map of World X from the provided scanner and populates the world with the cities and
    connections described. Panics if errors occur while reading the scanner or creating the cities and connections.
    - `GenerateAliens(numberAliens int)` → Generates aliens one at a time placing them in a random empty city.
    Panics on the tentative to generate more aliens than the number of cities.
    - `RunSimulation(writer *bufio.Writer)` → Simulates invasion moving each alien `defaultMaxIterations` times 
    or until it's trapped in an isolated city. When two aliens meet in the same city they fight and in the process,
    both aliens die and the city is destroyed severing all its connections.
    Prints message to the writer for every city destroyed. 

## Usage

### Install Packages

```shell script
$ go get github.com/tomasnunes/invasion/pkg/worldx
```

### Run Program
#### Locally:
```shell script
$ cd invasion
$ go run cmd/invasion/main.go [-h] [N] [INPUT_FILE] [OUTPUT_INVASION]
```

#### Build:
```shell script
$ cd invasion
$ go build ./cmd/invasion
$ ./invasion [-h] [N] [INPUT_FILE] [OUTPUT_FILE]
```

- -h          → Prints help message
- N           → Number of alien invaders, if none provided defaults to `defaultNumberAliens`
- INPUT_FILE  → Name of the input file with the world map description, if none provided defaults to `defaultInputFile`
- OUTPUT_FILE → Name of the output file to create and print the program information,
if none provided defaults to the `stdout`

#### Tests:
```shell script
$ cd invasion/pkg/worldx
$ go test
```

## Assumptions

- The city names and the alien names are unique.
- The number of aliens is less than the number of cities otherwise a city would be destroyed immediately.
- The roads are bidirectional, an Alien invading the world would not be stopped by a unidirectional road, although
the provided map doesn't need to specify both connections, one is enough to generate the bidirectional connection.
- City names are case-sensitive and cannot contain spaces, any other character is allowed.
- The order in which the cities are printed is irrelevant.
- The connections to each city are printed in the following order `north=<...> south=<...> east=<...> west=<...>`
independently of the order in which they were read.

## Trade-Offs, Optimizations and Possible Changes

- Being this a very simple application I didn't find the need to use a library to create the CLI.
If this application was meant to be extended in the future I would have used either 
[`https://github.com/spf13/cobra`](https://github.com/spf13/cobra) or
[`https://github.com/urfave/cli`](https://github.com/urfave/cli).
- For the connections between cities I decided to use a fixed-size array `[MaxDirections]*City`,
another option would be a map, `map[Direction]*City`, mapping the direction of the connection
to the connected city. Using an array requires the allocation of a fixed amount of memory, even if 
the city doesn't have connections in all directions, but an array is way more efficient than the map.
**TL;DR:** Could change `connectedCities` to `map[Direction]*City` if the gain in memory out-weights
the loss in performance.
- If crypto-level randomness was required I would have used `crypto/rand` instead of `math/rand`.
The latter is enough for the required use cases and much more efficient.
-There were more efficient ways to obtain random empty cities and a random direction to move than trial and
error but it wouldn't be as random, ergo my implementation. 
- Could add concurrency, for example when creating cities and connections adding a `sync.Mutex` to the
`City struct` and verifying if a city with the same name already exists in the world using a `sync.RWMutex`
in the `WorldX struct` allowing multiple concurrent reads.
