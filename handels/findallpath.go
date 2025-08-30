package lemin

import (
	"log"
	"sort"
)

// FindPaths returns a pruned, sorted slice of viable paths from start to end.
func FindPaths(antFarm *Antsfarm) []*Path {
	// Get all raw paths via BFS
	paths := BFS(antFarm)
	if len(paths) == 0 {
		log.Fatal("No paths found from start to end")
	}
	// Reduce to best non-overlapping paths
	findBestPaths(antFarm, &paths)
	
	// * Sort paths by length (shortest → longest)
	// the first input is the path being sorted, the second input is the function
	// the function works like this [1 2] 3 checks the condition if it's true it swaps them
	// then 1 [2 3]
	sort.Slice(paths, func(i, j int) bool {
		return len(paths[i].Rooms) > len(paths[j].Rooms)
	})

	// Remove paths that are not useful

	return paths
}

// findBestPaths selects the best set of non-overlapping paths from all paths
func findBestPaths(antFarm *Antsfarm, paths *[]*Path) {
	n := len(*paths)
	if n == 0 {
		return
	}

	// Build overlap map: overlaps[i][j] = true if path i and j share rooms
	overlaps := make(map[int]map[int]bool)
	for i := range *paths {
		overlaps[i] = make(map[int]bool)
		for j := range *paths {
			if i != j && pathsOverlap((*paths)[i], (*paths)[j]) {
				overlaps[i][j] = true
			}
		}
	}

	// Find best combination using backtracking
	bestCombo := findBestCombination(antFarm, *paths, overlaps, n)

	// Collect the paths that are in the best combination
	var bestPaths []*Path
	for i := range *paths {
		if bestCombo[i] {
			bestPaths = append(bestPaths, (*paths)[i])
		}
	}

	// Update the original slice
	*paths = bestPaths
}

// pathsOverlap checks if two paths share any rooms (except start and end)
func pathsOverlap(a, b *Path) bool {
	// Create a map for rooms in path1
	rooms := make(map[string]bool)
	for i := 0; i < len(a.Rooms)-1; i++ {
		rooms[a.Rooms[i].Name] = true
	}
	for i := 0; i < len(b.Rooms)-1; i++ {
		if rooms[b.Rooms[i].Name] {
			return true
		}
	}
	return false
}

// ?find number of rounds
// Relies on two package-level globals that must exist in this package:
// var allPaths []*Path
// var totalAnts int
// Compute which set yields fewer rounds using Antsfarm + paths (no globals).
// Returns -1 if currentSet is better (fewer rounds), 1 if bestSet is better, 0 if tie.
func numberOfRounds(antFarm *Antsfarm, paths []*Path, bestSet map[int]bool, currentSet map[int]bool) int {
	const INF = int(^uint(0) >> 1)

	roundsFor := func(sel map[int]bool) int {
		if len(sel) == 0 || antFarm == nil || antFarm.TotalAnts <= 0 {
			return INF
		}

		// B_i = len(Rooms) - 1 (base travel time for an empty path)
		base := make([]int, 0, len(sel))
		for i, ok := range sel {
			if !ok || i < 0 || i >= len(paths) || paths[i] == nil {
				continue
			}
			b := len(paths[i].Rooms) - 1
			base = append(base, b)
		}
		if len(base) == 0 {
			return INF
		}

		// Binary search bounds: low at shortest base, high at (longest base + total ants)
		minB, maxB := base[0], base[0]
		for _, b := range base[1:] {
			if b < minB {
				minB = b
			}
			if b > maxB {
				maxB = b
			}
		}
		low, high := minB, maxB+antFarm.TotalAnts

		// Check if T rounds can deliver all ants across the selected paths
		canDeliver := func(T int) bool {
			sum := 0
			for _, b := range base {
				c := T - b + 1 // contribution of this path
				if c > 0 {
					sum += c
					if sum >= antFarm.TotalAnts {
						return true
					}
				}
			}
			return false
		}

		for low < high {
			mid := (low + high) / 2
			if canDeliver(mid) {
				high = mid
			} else {
				low = mid + 1
			}

			// safety net: prevent infinite loop
			if high-low <= 1 {
				if canDeliver(low) {
					return low
				}
				return high
			}
		}

		return low
	}

	bestT := roundsFor(bestSet)
	currT := roundsFor(currentSet)

	if currT < bestT {
		return -1
	}
	if currT > bestT {
		return 1
	}
	return 0
}

// findBestCombination uses backtracking to find the best combination of non-overlapping paths
// It picks the set with the least number of rounds, not just the most paths
func findBestCombination(antFarm *Antsfarm, paths []*Path, overlaps map[int]map[int]bool, n int) map[int]bool {

	bestCombination := make(map[int]bool)
	currentCombination := make(map[int]bool)

	var tryPath func(pathIndex int)
	tryPath = func(pathIndex int) {
		if pathIndex == n {
			// Compare with best using numberOfRounds
			cmp := numberOfRounds(antFarm, paths, bestCombination, currentCombination)
			if cmp < 0 {
				// currentCombination is better
				bestCombination = make(map[int]bool)
				for k, v := range currentCombination {
					bestCombination[k] = v
				}
			}
			return
		}

		// Try excluding current path
		currentCombination[pathIndex] = false
		tryPath(pathIndex + 1)

		// Try including current path if no conflict
		canInclude := true
		for other := 0; other < pathIndex; other++ {
			if currentCombination[other] && overlaps[pathIndex][other] {
				canInclude = false
				break
			}
		}
		if canInclude {
			currentCombination[pathIndex] = true
			tryPath(pathIndex + 1)
			currentCombination[pathIndex] = false
		}
	}

	tryPath(0)
	return bestCombination
}

//? FilterCheckCompare compares paths and only keeps the useful ones.

func UpdatePath(antFarm *Antsfarm, path *Path, dist *map[string]int) {
	for i := len(path.Ants) - 1; i >= 0; i-- {
		ant := path.Ants[i]
		finished := ant.Move(antFarm)
		DisplayMovement(ant.ID, ant.Path.Rooms[ant.Position].Name)
		if finished {
			path.Ants = append(path.Ants[:i], path.Ants[i+1:]...)
		}
	}
	_ = dist
}
