package main

import (
	"fmt"
	lemin "lemin/handels"
)

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

func main() {
    // Create a new graph
    farm := lemin.ReadFile()
    if farm == nil {
        fmt.Println("Error reading file")
        return
    }
    fmt.Println(farm.String())

   
    
}
