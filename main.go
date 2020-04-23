package main

import (
	"fmt"
	"os"
)

func main() {
	if err := run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func usage() error {
	fmt.Printf("Usage: \n")
	fmt.Printf("\t%s save - saves the stats to a json file\n", os.Args[0])
	fmt.Printf("\t%s plot - creates a plot from the json file\n", os.Args[0])
	return nil
}

func run() error {
	if len(os.Args) != 2 {
		return usage()
	}

	switch os.Args[1] {
	case "save":
		return save()
	case "plot":
		return createPlot()
	default:
		return usage()
	}
}
