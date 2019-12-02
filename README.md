# Simulate Invasion

This command-line application reads and constructs a World X, generates aliens, allocates them to an empty city,
and simulates an invasion. During the simulation, aliens move between cities at random, when two aliens meet in
the same city they fight and in the process both aliens die and the city is destroyed severing all its connections.
The final state of the world is printed at the end.

### Input FILE Fromat (World Map Description)

    <new city name> [north=<connected city name> south=<...> east=<...> west=<...>]
    ...

- City names cannot contain spaces, any other character is allowed.
- No effort is made to verify if a `<new city name>` is empty.
- Any directional connection is optional, a city doesn't need to connect to other cities in all directions.

### Packages

- `invade`

    This package depends on the `worldx` package and provides the logic to create a new World X from a file and
    run the invasion with a specified number of aliens.

- `worldx`

    This package contains the `WorldX`, `Alien`, `City`, and `Direction` types and possible interactions with 
    this types, like creating cities, generating aliens and running the simulation.

## Usage

### Install Packages

    $ go get github.com/tomasnunes/invasion/pkg/invade
    $ go get github.com/tomasnunes/invasion/pkg/worldx

### Run Program

    $ ./invasion [-h] [N] [FILE]

-h → Prints help message

N → Number of alien invaders, if none provided defaults to `defaultNumberAliens`

FILE → Name of the input file with the world map description, if none provided defaults to `defaultInputFile`

## Assumptions

- The city names and the alien names are unique.
- The number of aliens is less than the number of cities otherwise a city would be destroyed immediately.
- The roads are bidirectional, an Alien invading the world would not be stopped by a unidirectional road, although
the provided map doesn't need to specify both connections, one is enough to generate the bidirectional connection.
- City names cannot contain spaces, any other character is allowed.

## Trade-Offs, Optimizations and Possible Changes

- Being this a very simple application I didn't find the need to use a library to create the CLI.
If this application was meant to be extended in the future I would have used either 
[`https://github.com/spf13/cobra`](https://github.com/spf13/cobra) or
[`https://github.com/urfave/cli`](https://github.com/urfave/cli).
- For the connections between cities I decided to use maps, mapping the direction of the connection
to the connected city, I believe it's more readable and intuitive, but it comes at an efficiency cost, 
using a fixed-sized array and mapping the directions to a number between 0-3 would be much more efficient
although every city would have an array with size 4 even if the city doesn't have connections in all directions.
**TL;DR**: Change `ConnectedCities` to `[4]*City` instead of `map[Direction]*City` if the gain in time efficiency 
out-weights the loss in readability.
- If crypto-level randomness was required I would have used `crypto/rand` instead of `math/rand`.
The latter is enough for the required use cases and much more efficient.
