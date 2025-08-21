// TODO: find  all path we can take from start to end
package lemin

import (
	// "fmt"
	"sort"
)

func FindPaths(antFarm *Antsfarm) []*Path {
	paths := BFS(antFarm)
	findBestPaths(&paths)
	// findBestPaths(&paths)
	// * Sort paths by length (shortest → longest)
	// the first input is the path being sorted, the second input is the function
	// the function works like this [1 2] 3 checks the condition if it's true it swaps them
	// then 1 [2 3]
	sort.Slice(paths, func(i, j int) bool {
		return len(paths[i].Rooms) > len(paths[j].Rooms)
	})

	paths = FilterCheckCompare(paths, antFarm.TotalAnts)

	return paths
}

// findBestPaths selects the best set of non-overlapping paths from all paths
func findBestPaths(paths *[]*Path) {
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
	bestCombo := findBestCombination(overlaps, n)

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
func pathsOverlap(path1, path2 *Path) bool {
	// Create a map for rooms in path1
	rooms := make(map[string]bool)

	// Add all rooms from path1 (except first and last) to the map
	for i := 1; i < len(path1.Rooms)-1; i++ {
		rooms[path1.Rooms[i].Name] = true
	}

	// Check if any room from path2 (except first and last) exists in path1
	for i := 1; i < len(path2.Rooms)-1; i++ {
		if rooms[path2.Rooms[i].Name] {
			return true
		}
	}
	return false
}

// findBestCombination uses backtracking to find the best combination of non-overlapping paths
func findBestCombination(overlaps map[int]map[int]bool, n int) map[int]bool {
	bestCombination := make(map[int]bool)
	currentCombination := make(map[int]bool)

	var tryPath func(pathIndex int) int
	tryPath = func(pathIndex int) int {
		if pathIndex == n {
			// Count current combination size
			count := 0
			for _, included := range currentCombination {
				if included {
					count++
				}
			}

			// Update best combination if current is better
			bestCount := 0
			for _, included := range bestCombination {
				if included {
					bestCount++
				}
			}
			if count > bestCount {
				for k, v := range currentCombination {
					bestCombination[k] = v
				}
			}
			return count
		}

		// Try excluding current path
		currentCombination[pathIndex] = false
		exclude := tryPath(pathIndex + 1)

		// Check if we can include current path
		canInclude := true
		for other := 0; other < pathIndex; other++ {
			if currentCombination[other] && overlaps[pathIndex][other] {
				canInclude = false
				break
			}
		}

		var include int
		if canInclude {
			currentCombination[pathIndex] = true
			include = tryPath(pathIndex + 1)
			if include > exclude {
				return include
			}
		}

		currentCombination[pathIndex] = false
		return exclude
	}

	tryPath(0)
	return bestCombination
}

// FilterCheckCompare compares paths and only keeps the useful ones.
func FilterCheckCompare(paths []*Path, totalAnts int) []*Path {
	if len(paths) == 0 {
		return paths
	}

	shortest := len(paths[0].Rooms)
	var worthwhile []*Path

	for i, p := range paths {

		if (len(p.Rooms) - shortest) <= (totalAnts / (i + 1)) {
			worthwhile = append(worthwhile, p)
		}
	}

	return worthwhile
}

func UpdatePath(antFarm *Antsfarm, path *Path, dist *map[string]int) {
	for i := len(path.Ants) - 1; i >= 0; i-- {
		ant := path.Ants[i]
		finished := ant.Move(antFarm)
		DisplayMovement(ant.ID, ant.Path.Rooms[ant.Position].Name)
		if finished {
			// safely remove
			path.Ants = append(path.Ants[:i], path.Ants[i+1:]...)
		}
	}
	if (*dist)[path.String()] > 0 {

		//assign ants to path if there are still ants left
		
		// if (*dist)[path] <= 0 {
		// 	return
		// }

		// // Only assign ant in start room (Position = 0)
		// newAnt := &Ant{
		// 	ID:       antFarm.NextAnt + 1,
		// 	Path:     path,
		// 	Position: 0, // start room
		// }
		// path.Ants = append([]*Ant{newAnt}, path.Ants...)
		// path.Rooms[0].Ant = newAnt // assign ant to start room
		
		// antFarm.NextAnt++
		// (*dist)[path]--
		// DisplayMovement(newAnt.ID, antFarm.StartRoom.Name)
	}
}
