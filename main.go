package main

import (
	"fmt"
	lemin "lemin/handels"
	"log"
	"os"
)

/* !! remove this
// Graph struct with adjacency list
type Graph struct {
	adj map[string][]string
}

// Constructor: creates a new Graph with empty adjacency list
func NewGraph() *Graph {
	return &Graph{adj: make(map[string][]string)}
}

// AddEdge adds an undirected edge between room a and room b
func (g *Graph) AddEdge(a, b string) {
	g.adj[a] = append(g.adj[a], b) // a → b
	g.adj[b] = append(g.adj[b], a) // b → a
}

// PrintGraph prints the adjacency list of the graph
func (g *Graph) PrintGraph() {
	fmt.Println("Graph adjacency list:")
	for room, neighbors := range g.adj {
		fmt.Printf("%s: %v\n", room, neighbors)
	}
}
 ! */

func main() {
	// * take input file name from command line arguments

	// ! display movement (DisplayMovement handles its own output)
	// lemin.DisplayMovement(1, "A")
	// return

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

	lemin.RunAlgorithm(farm)

}
