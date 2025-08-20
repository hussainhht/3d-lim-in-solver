package lemin

import (
	"fmt"
	"sort"
)

// todo run algo

func RunAlgorithm(antFarm *Antsfarm) {
    paths := FindPaths(antFarm)
    fmt.Println("Found paths:")
    for _, path := range paths {
        fmt.Println(path.String())
    }

    // Precompute ant distribution
    dist := computeAntDistribution(paths, antFarm.totalAnts)
	fmt.Println("Ant distribution:", dist)
    for antFarm.doneAnts < antFarm.totalAnts {
        resetVisited(paths)
        for _, path := range paths { // move ants in paths forward
            UpdatePath(antFarm, path)
        }
        assignAntsToPath(antFarm, paths, dist)
        fmt.Println()
    }
}


func resetVisited(paths []*Path) {
	for _, path := range paths {
		for _, room := range path.rooms {
			room.visited = false
		}
	}
}

func assignAntsToPath(antFarm *Antsfarm, paths []*Path, dist map[*Path]int) {
    for _, path := range paths {
        if dist[path] <= 0 {
            continue // quota used
        }

        firstRoom := path.rooms[0]
        // Only place an ant if room is free (except start)
        if firstRoom.ant != nil && firstRoom != antFarm.startRoom {
            continue
        }

        newAnt := &Ant{
            id:       antFarm.nextAnt,
            path:     path,
            position: 0,
        }
        path.ants = append([]*Ant{newAnt}, path.ants...)
        firstRoom.ant = newAnt
        antFarm.nextAnt++
        dist[path]--
        DisplayMovement(newAnt.id, firstRoom.name)

        if antFarm.nextAnt >= antFarm.totalAnts {
            break
        }
    }
}




func isPathAvailable(path *Path) bool {
	if path.rooms[0].ant != nil || path.rooms[0].visited {
		return false
	}
	return true
}

func DisplayMovement(antID int, roomName string) {
	fmt.Printf("L%v->%v ", antID, roomName)
}

// TODO : do
func MoveAnt() {

}

// TODO: Implement pathfinding algorithm

func BFS(currentRoom *Room, endRoom *Room, occupied *map[string]bool, antFarm *Antsfarm) []*Room {
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
			if !visited[neighbor.name] && !(*occupied)[neighbor.name] {
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

	for _, room := range path {
		if room.name == antFarm.endRoom.name {
			continue
		}

		(*occupied)[room.name] = true
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

    // Sort paths by length (shortest → longest)
    sort.Slice(paths, func(i, j int) bool {
        return len(paths[i].rooms) < len(paths[j].rooms)
    })

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
            turns := len(path.rooms) + dist[path]
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
