package lemin

import (
	"fmt"
)

func RunAlgorithm() {
	// antfarm := ReadFile()
	// totalAnts := antfarm.antsnum
	// doneAnts := 0
	// nextAntOut := 0
	// paths := findPaths(antfarm)
	// for doneAnts < totalAnts {
	// }
}

func findPaths(antFarm *Antsfarm) [][]*Room {
	startRoom := antFarm.startroom
	var paths [][]*Room
	for _, room := range startRoom.links {
		path := BFS(room, antFarm.endroom)
		if path != nil {
			paths = append(paths, path)
		}
	}
	paths = sortPaths(paths)
	return paths
}

func sortPaths(paths [][]*Room) [][]*Room {
	if len(paths) <= 1 {
		return paths
	}
	// Merge sort
	left := sortPaths(paths[:len(paths)/2])
	right := sortPaths(paths[len(paths)/2:])
	var result [][]*Room
	i := 0
	j := 0

	for i < len(left) && j < len(right) {
		if len(left[i]) <= len(right[j]) {
			result = append(result, left[i])
			i++
		} else {
			result = append(result, right[j])
			j++
		}
	}
	for i < len(left) {
		result = append(result, left[i])
		i++
	}
	for j < len(right) {
		result = append(result, right[j])
		j++
	}
	return result
}

func BFS(currentRoom *Room, endRoom *Room) []*Room {
	visited := make(map[string]bool)
	parent := make(map[string]*Room)
	queue := []*Room{currentRoom}

	visited[currentRoom.name] = true
	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]

		if node.name == endRoom.name {
			break
		}

		for _, neighbor := range node.links {
			if !visited[neighbor.name] {
				visited[neighbor.name] = true
				parent[neighbor.name] = node
				queue = append(queue, neighbor)
			}
		}
	}

	if !visited[endRoom.name] {
		return nil
	}

	path := []*Room{}
	current := endRoom

	for current != nil {
		path = append([]*Room{current}, path...)
		current = parent[current.name]
	}
	return path
}

func RunStressTests() {
	// Create rooms
	rooms := make(map[string]*Room)
	for _, name := range []string{
		"A", "B", "C", "D", "E", "F", "G", "H", "I", "J",
		"K", "L", "M", "N", "O", "P", "Q", "R", "S", "T",
		"U", "V", "W", "X", "Y", "Z",
	} {
		rooms[name] = &Room{name: name}
	}

	// Main graph
	rooms["A"].links = []*Room{rooms["B"], rooms["C"]}
	rooms["B"].links = []*Room{rooms["A"], rooms["D"], rooms["E"]}
	rooms["C"].links = []*Room{rooms["A"], rooms["E"]}
	rooms["D"].links = []*Room{rooms["B"], rooms["F"]}
	rooms["E"].links = []*Room{rooms["B"], rooms["C"], rooms["F"], rooms["H"]}
	rooms["F"].links = []*Room{rooms["D"], rooms["E"], rooms["G"]}
	rooms["G"].links = []*Room{rooms["F"], rooms["I"]}
	rooms["H"].links = []*Room{rooms["E"], rooms["I"]}
	rooms["I"].links = []*Room{rooms["G"], rooms["H"], rooms["J"]}
	rooms["J"].links = []*Room{rooms["I"]}

	// Disconnected cluster
	rooms["K"].links = []*Room{rooms["L"]}
	rooms["L"].links = []*Room{rooms["K"]}
	rooms["M"].links = []*Room{rooms["N"], rooms["O"]}
	rooms["N"].links = []*Room{rooms["M"], rooms["O"]}
	rooms["O"].links = []*Room{rooms["M"], rooms["N"]}

	// Additional small chains for variety
	rooms["P"].links = []*Room{rooms["Q"]}
	rooms["Q"].links = []*Room{rooms["P"], rooms["R"]}
	rooms["R"].links = []*Room{rooms["Q"], rooms["S"]}
	rooms["S"].links = []*Room{rooms["R"], rooms["T"]}
	rooms["T"].links = []*Room{rooms["S"]}

	// Define tests
	tests := []struct {
		start *Room
		end   *Room
		name  string
	}{
		// Main cluster paths
		{rooms["A"], rooms["D"], "Check 1: A -> D"},
		{rooms["A"], rooms["E"], "Check 2: A -> E"},
		{rooms["A"], rooms["G"], "Check 3: A -> G"},
		{rooms["A"], rooms["J"], "Check 4: A -> J"},
		{rooms["F"], rooms["A"], "Check 5: F -> A"},
		{rooms["B"], rooms["B"], "Check 6: B -> B"},

		// Disconnected / isolated nodes
		{rooms["A"], rooms["K"], "Check 7: A -> K"},
		{rooms["L"], rooms["A"], "Check 8: L -> A"},
		{rooms["K"], rooms["L"], "Check 9: K -> L"},
		{rooms["P"], rooms["T"], "Check 10: P -> T"},
		{rooms["T"], rooms["P"], "Check 11: T -> P"},
		{rooms["U"], rooms["V"], "Check 12: U -> V"}, // both isolated
		{rooms["V"], rooms["U"], "Check 13: V -> U"}, // both isolated

		// Cycles / multiple paths
		{rooms["E"], rooms["G"], "Check 14: E -> G"},
		{rooms["I"], rooms["E"], "Check 15: I -> E"},
		{rooms["H"], rooms["G"], "Check 16: H -> G"},
		{rooms["M"], rooms["O"], "Check 17: M -> O"},
		{rooms["O"], rooms["M"], "Check 18: O -> M"},

		// Deeper chains / alternative paths
		{rooms["A"], rooms["H"], "Check 19: A -> H"},
		{rooms["A"], rooms["F"], "Check 20: A -> F"},
		{rooms["B"], rooms["G"], "Check 21: B -> G"},
		{rooms["C"], rooms["I"], "Check 22: C -> I"},
		{rooms["D"], rooms["J"], "Check 23: D -> J"},
		{rooms["F"], rooms["H"], "Check 24: F -> H"},
		{rooms["P"], rooms["R"], "Check 25: P -> R"},
		{rooms["Q"], rooms["T"], "Check 26: Q -> T"},

		// Start = end edge cases
		{rooms["M"], rooms["M"], "Check 27: M -> M"},
		{rooms["K"], rooms["K"], "Check 28: K -> K"},
		{rooms["A"], rooms["A"], "Check 29: A -> A"},
		{rooms["P"], rooms["P"], "Check 30: P -> P"},

		// Mix of main and disconnected clusters
		{rooms["A"], rooms["M"], "Check 31: A -> M"},
		{rooms["M"], rooms["A"], "Check 32: M -> A"},
		{rooms["F"], rooms["O"], "Check 33: F -> O"},
		{rooms["O"], rooms["F"], "Check 34: O -> F"},
		{rooms["J"], rooms["H"], "Check 35: J -> H"},
		{rooms["T"], rooms["P"], "Check 36: T -> P"},

		// Some longer chains
		{rooms["P"], rooms["S"], "Check 37: P -> S"},
		{rooms["Q"], rooms["S"], "Check 38: Q -> S"},
		{rooms["R"], rooms["T"], "Check 39: R -> T"},
		{rooms["S"], rooms["P"], "Check 40: S -> P"},

		// Random deep paths
		{rooms["A"], rooms["I"], "Check 41: A -> I"},
		{rooms["B"], rooms["H"], "Check 42: B -> H"},
		{rooms["C"], rooms["G"], "Check 43: C -> G"},
		{rooms["D"], rooms["H"], "Check 44: D -> H"},
		{rooms["E"], rooms["J"], "Check 45: E -> J"},
		{rooms["F"], rooms["J"], "Check 46: F -> J"},
		{rooms["G"], rooms["H"], "Check 47: G -> H"},
		{rooms["H"], rooms["J"], "Check 48: H -> J"},
		{rooms["I"], rooms["A"], "Check 49: I -> A"},
		{rooms["J"], rooms["B"], "Check 50: J -> B"},
	}

	// Run tests
	for _, t := range tests {
		path := BFS(t.start, t.end)
		fmt.Print(t.name + ": ")
		if path == nil {
			fmt.Println("nil (no path)")
		} else {
			for _, room := range path {
				fmt.Print(room.name, " ")
			}
			fmt.Println()
		}
	}
}
