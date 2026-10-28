package lemin

import (
	"fmt"
	"slices"
)

// todo run algo

// RunAlgorithm simulates the colony and prints one line per turn to stdout.
func RunAlgorithm(antFarm *Antsfarm) {
	for _, turn := range Simulate(antFarm).Turns {
		for _, m := range turn {
			DisplayMovement(m.AntID, m.Room)
		}
		fmt.Println()
	}
}

// Simulate finds the paths, distributes the ants and plays the colony turn by
// turn until every ant has reached the end room. It prints nothing.
func Simulate(antFarm *Antsfarm) *Result {

	paths := FindPaths(antFarm)

	dist := computeAntDistribution(paths, antFarm.TotalAnts)

	result := &Result{Paths: paths}
	for antFarm.DoneAnts < antFarm.TotalAnts {
		resetVisited(paths)
		for _, path := range paths { // move ants in paths forward
			UpdatePath(antFarm, path)
		}
		AssignAntsToPath(antFarm, paths, dist)
		result.Turns = append(result.Turns, antFarm.flushTurn())
	}
	return result
}

// * reset visited status of rooms in paths
func resetVisited(paths []*Path) {
	for _, path := range paths {
		for _, room := range path.Rooms {
			room.Visited = false
		}
	}
}

// * FindPaths uses BFS to find all paths from start to end in the ant farm.
func BFS(antFarm *Antsfarm) []*Path {

	var results []*Path
	start := antFarm.StartRoom

	// Queue of paths; each path is a slice of Room pointers
	queue := [][]*Room{}

	// Initialize queue with start room's neighbors
	for _, link := range start.Links {
		queue = append(queue, []*Room{link})
	}

	for len(queue) > 0 {
		if antFarm.budget != nil {
			antFarm.budget.enumeration(len(queue), len(results))
		}
		// Dequeue first path
		path := queue[0]
		queue = queue[1:]

		lastRoom := path[len(path)-1]

		// Check if we reached the end room
		if lastRoom == antFarm.EndRoom {
			// Found a complete path, store without start room
			results = append(results, &Path{Rooms: path})
			continue
		}

		// Explore neighbors
		for _, neighbor := range lastRoom.Links {
			// Skip rooms already in path OR start room
			if !roomInPath(path, neighbor) && neighbor != start {
				newPath := append([]*Room{}, path...) // copy current path
				newPath = append(newPath, neighbor)
				queue = append(queue, newPath)
			}
		}
	}

	if antFarm.budget != nil {
		antFarm.budget.enumeration(0, len(results))
	}
	return results
}

// Helper: check if a room is already in the path (to avoid cycles)
func roomInPath(path []*Room, room *Room) bool {
	return slices.Contains(path, room)
}

func computeAntDistribution(paths []*Path, totalAnts int) *map[string]int {
	// Map: path → number of ants assigned
	dist := make(map[string]int)
	if len(paths) == 0 {
		return &dist
	}

	//! not important
	// sort.Slice(paths,
	// 	func(i, j int) bool {
	// 		return len(paths[i].Rooms) < len(paths[j].Rooms)
	// 	},
	// )

	// Calculate distribution using "turns equalization"
	remaining := totalAnts
	for _, path := range paths {
		dist[path.String()] = 0
	}

	// Find path with minimum (length + assigned ants)
	for remaining > 0 {
		// Find path with minimum (length + assigned ants)

		minTurns := int(^uint(0) >> 1) // Max int [minTurns = 9223372036854775807]
		var bestPath *Path
		for _, path := range paths {
			turns := len(path.Rooms) + dist[path.String()]
			if turns < minTurns {
				minTurns = turns
				bestPath = path
			}
		}
		dist[bestPath.String()]++
		remaining--
	}

	return &dist
}
