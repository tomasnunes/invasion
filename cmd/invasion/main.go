package main

import (
    "bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/tomasnunes/invasion/pkg/worldx"
)

const (
	defaultNumberAliens int = 10
	defaultInputFile string = "test/world_map"
)

func printUsage() {
	fmt.Printf(
		"invasion - This program reads and constructs world X, simulates an alien invasion and prints the final state of the world.\n" +
		"\n" +
		"Usage:\n" +
		"%s [-h] [N] [FILE]\n" +
		"\n" +
		"Flags:\n" +
		"-h\tPrints this message.\n" +
		"\n" +
		"Args:\n" +
		"N\tNumber of alien invaders, if none provided defaults to %d.\n" +
		"FILE\tName of the input file with the world map description, if none provided defaults to '%s'.\n",
		os.Args[0], defaultNumberAliens, defaultInputFile)
}

func main() {
	var helpFlag bool
	flag.BoolVar(&helpFlag, "h", false, "Prints usage message.")
	flag.Parse()
	if helpFlag {
		printUsage()
		return
	}

	args := os.Args
	var numberAliens int
	if len(args) > 1 {
		var err error
		if numberAliens, err = strconv.Atoi(args[1]); err != nil {
			printUsage()
			log.Fatal(err)
		}
	} else {
		numberAliens = defaultNumberAliens
	}

	var filename string
	if len(args) > 2 {
		filename = args[2]
	} else {
		filename = defaultInputFile
	}

	if info, err := os.Stat(filename); err != nil {
		printUsage()
		log.Fatal(err)
	} else if info.IsDir() {
		printUsage()
		log.Fatalf("%s: is a directory, should be a file with the description of the world map.", filename)
	}

	Invade(filename, numberAliens)
}

// Creates world map with the description on the file, generates aliens, runs the simulation of the invasion,
// and prints the final state of the world
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

    world := worldx.WorldX{}

    scanner := bufio.NewScanner(file)
    world.ReadWorldMap(scanner)
    world.GenerateAliens(numberAliens)
    world.RunSimulation()

    fmt.Print(world.String())
}
