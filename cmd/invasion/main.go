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
		"%s [-h] [N] [INPUT_FILE] [OUTPUT_FILE]\n" +
		"\n" +
		"Flags:\n" +
		"-h\t\tPrints this message.\n" +
		"\n" +
		"Args:\n" +
		"N\t\tNumber of alien invaders, if none provided defaults to %d.\n" +
		"INPUT_FILE\tName of the input file with the world map description, if none provided defaults to '%s'.\n" +
        "OUTPUT_FILE\tName of the output file to create and print program information, if none provided defaults to the stdout.\n",
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
	totalArgs := len(args)
	var numberAliens int
	if totalArgs > 1 {
		var err error
		if numberAliens, err = strconv.Atoi(args[1]); err != nil {
			defer printUsage()
			log.Panic(err)
		}
	} else {
		numberAliens = defaultNumberAliens
	}

	var filename string
	if totalArgs > 2 {
		filename = args[2]
	} else {
		filename = defaultInputFile
	}

	var writer *bufio.Writer
	if totalArgs > 3 {
        if outputFile, err := os.Create(args[3]); err != nil {
            defer printUsage()
            log.Panic(err)
        } else {
            writer = bufio.NewWriter(outputFile)
        }
    } else {
        writer = bufio.NewWriter(os.Stdout)
    }

	if info, err := os.Stat(filename); err != nil {
		defer printUsage()
		log.Panic(err)
	} else if info.IsDir() {
		defer printUsage()
		log.Panicf("%s is a directory, should be a file with the description of the world map.", filename)
	}

	Invade(filename, numberAliens, writer)
}

// Creates world map with the description on the file, generates aliens, runs the simulation of the invasion,
// and prints the final state of the world.
func Invade(filename string, numberAliens int, writer *bufio.Writer) {
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
    world.RunSimulation(writer)

    if _, err := fmt.Fprint(writer, world.String()); err != nil {
        log.Panic(err)
    }
    if err := writer.Flush(); err != nil {
        log.Panic(err)
    }
}
