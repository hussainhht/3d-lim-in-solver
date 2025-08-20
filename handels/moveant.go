//todo: calculate the number of moves of each path

package lemin

//! this is all not used

func PathCost(pathLen, assigned int) int {
	// ? (pathLen - 1) + assigned
	return (pathLen - 1) + assigned
}

func AssignGreedy(totalAnts int, paths []*Path) []int {

	assigned := make([]int, len(paths))

	if totalAnts == 0 || len(paths) == 0 {
		return assigned
	}

	bestPath := 0

	for i := 0; i < totalAnts; i++ {
		cost := PathCost(len(paths[i].Rooms), assigned[i])
		if cost < PathCost(len(paths[bestPath].Rooms), assigned[bestPath]) {
			bestPath = i
		}
	}

	assigned[bestPath] += totalAnts

	return assigned
}
