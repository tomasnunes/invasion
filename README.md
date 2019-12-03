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

    This package contains the `WorldX` type and possible interactions with the world to run a full simulation,
    the exported functions are described bellow:
    - `ReadWorldMap()` → Reads map of World X from the provided scanner and populates world with the cities and
    connections described. Panics if errors occur while reading the scanner or creating the cities and connections.

    - `GenerateAliens(numberAliens int)` → Generates aliens one at a time placing them in a random empty city.
    Panics on the tentative to generate more aliens than the number of cities.

    - `RunSimulation()` → Simulates invasion moving each alien `defaultMaxIterations` times and destroying city if
    aliens collide.

## Usage

### Install Packages

```shell script
$ go get github.com/tomasnunes/invasion/pkg/worldx
```

### Run Program

```shell script
$ ./invasion [-h] [N] [FILE]
```

- -h → Prints help message
- N → Number of alien invaders, if none provided defaults to `defaultNumberAliens`
- FILE → Name of the input file with the world map description, if none provided defaults to `defaultInputFile`

## Assumptions

- The city names and the alien names are unique.
- The number of aliens is less than the number of cities otherwise a city would be destroyed immediately.
- The roads are bidirectional, an Alien invading the world would not be stopped by a unidirectional road, although
the provided map doesn't need to specify both connections, one is enough to generate the bidirectional connection.
- City names cannot contain spaces, any other character is allowed.
- The connections to each city are printed in order `north=<...> south=<...> east=<...> west=<...>`,
independently of the order in which they were read.

## Trade-Offs, Optimizations and Possible Changes

- Being this a very simple application I didn't find the need to use a library to create the CLI.
If this application was meant to be extended in the future I would have used either 
[`https://github.com/spf13/cobra`](https://github.com/spf13/cobra) or
[`https://github.com/urfave/cli`](https://github.com/urfave/cli).
- For the connections between cities I decided to use a fixed-size array `[maxDirections]*city`,
another option would be a map, `map[direction]*city`, mapping the direction of the connection
to the connected city. Using an array requires the allocation of a fixed amount of memory, even if 
the city doesn't have connections in all directions, but an array is way more efficient than the map.
**TL;DR**: Could change `connectedCities` to `map[direction]*city` if the gain in memory out-weights
the loss in efficiency.
- If crypto-level randomness was required I would have used `crypto/rand` instead of `math/rand`.
The latter is enough for the required use cases and much more efficient.
