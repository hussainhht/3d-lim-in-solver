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

func BFS(antFarm *Antsfarm) []*Path {
	var results []*Path

	// Queue of paths; each path is a slice of Room pointers
	queue := [][]*Room{{antFarm.StartRoom}}

	for len(queue) > 0 {
		// Dequeue first path
		path := queue[0]
		queue = queue[1:]

		lastRoom := path[len(path)-1]

		if lastRoom == antFarm.EndRoom  {
			// Found a complete path
			results = append(results, &Path{Rooms: path[1:]})
			continue
		}

		// Explore neighbors
		for _, neighbor := range lastRoom.Links {
			if !roomInPath(path, neighbor) {
				newPath := append([]*Room{}, path...) // copy current path
				newPath = append(newPath, neighbor)
				queue = append(queue, newPath)
			}
		}
	}

	return results
}

// Helper: check if a room is already in the path (to avoid cycles)
func roomInPath(path []*Room, room *Room) bool {
	for _, r := range path {
		if r == room {
			return true
		}
	}
	return false
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
