package lemin

import (
	"fmt"
)

// todo run algo

func RunAlgorithm(antFarm *Antsfarm) {
	paths := FindPaths(antFarm)
	fmt.Println("Found paths:")

	//! this is not the final printing this is only for debugging
	for _, path := range paths {
		fmt.Println(path.String())
	}

	// Precompute ant distribution
	dist := computeAntDistribution(paths, antFarm.TotalAnts)
	fmt.Println("Ant distribution:", dist)
	for antFarm.DoneAnts < antFarm.TotalAnts {
		resetVisited(paths)
		for _, path := range paths { // move ants in paths forward
			UpdatePath(antFarm, path)
		}
		AssignAntsToPath(antFarm, paths, dist)
		fmt.Println()
	}
}

func resetVisited(paths []*Path) {
	for _, path := range paths {
		for _, room := range path.Rooms {
			room.Visited = false
		}
	}
}

// ? why this unused
// func isPathAvailable(path *Path) bool {
// 	if path.Rooms[0].Ant != nil || path.Rooms[0].Visited {
// 		return false
// 	}
// 	return true
// }

// TODO: Implement pathfinding algorithm

func BFS(currentRoom *Room, endRoom *Room, occupied *map[string]bool, antFarm *Antsfarm) []*Room {
	visited := make(map[string]bool)
	parent := make(map[string]*Room)
	queue := []*Room{currentRoom}

	visited[currentRoom.Name] = true
	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]

		if node.Name == endRoom.Name {
			break
		}

		for _, neighbor := range node.Links {
			if !visited[neighbor.Name] && !(*occupied)[neighbor.Name] {
				visited[neighbor.Name] = true
				parent[neighbor.Name] = node
				queue = append(queue, neighbor)
			}
		}
	}

	if !visited[endRoom.Name] {
		return nil
	}

	path := []*Room{}
	current := endRoom

	for current != nil {
		path = append([]*Room{current}, path...)
		current = parent[current.Name]
	}

	for _, room := range path {
		if room.Name == antFarm.EndRoom.Name {
			continue
		}

		(*occupied)[room.Name] = true
	}

	// [] braket
	// () paranthesis
	// {} curly braket

	return path
}

func computeAntDistribution(paths []*Path, totalAnts int) map[*Path]int {
	// Map: path → number of ants assigned
	dist := make(map[*Path]int)
	if len(paths) == 0 {
		return dist
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
		dist[path] = 0
	}

	for remaining > 0 {
		// Find path with minimum (length + assigned ants)
		minTurns := int(^uint(0) >> 1) // Max int
		var bestPath *Path
		for _, path := range paths {
			turns := len(path.Rooms) + dist[path]
			if turns < minTurns {
				minTurns = turns
				bestPath = path
			}
		}
		dist[bestPath]++
		remaining--
	}

	return dist
}
