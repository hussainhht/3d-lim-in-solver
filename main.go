package main

import (
	"fmt"
	lemin "lemin/handles"
	"log"
	"os"
)


func main() {
	// * take input file name from command line arguments
	args := os.Args[1:]
	if len(args) != 1 {
		fmt.Println("Usage: go run main.go <input_file>")
		return
	}

	// * Read the input file and parse the data
	FileName := args[0]
	farm := lemin.ReadFile(FileName)
	if farm == nil {
		log.Println("Error reading file")
		return
	}

	//* write the input data to the output file
	lemin.WriteFile(lemin.TakeData(FileName))
	// lemin.TakeData(FileName)

	lemin.RunAlgorithm(farm)

}
