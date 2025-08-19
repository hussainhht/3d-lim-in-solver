package lemin

import (
	"fmt"
)

// todo run algo

func RunAlgorithm() {
	antfarm := ReadFile()
	totalAnts := antfarm.antsnum
	doneAnts := 0
	// nextAntOut := 0
	paths := FindPaths(antfarm)
	for doneAnts < totalAnts {
		for _, path := range paths {
			updatePath(path)
		}
		fmt.Println()
	}
}

func updatePath(path *Path) {
	for _, ant := range path.ants {
		if ant.position >= len(ant.path.rooms) {
			continue
		}
		ant.Move()
		DisplayMovement(ant.id, ant.path.rooms[ant.position].name)
	}
}

func DisplayMovement(antID int, roomName string) {
	fmt.Printf("L%v->%v", antID, roomName)
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
		if room.name == antFarm.endroom.name {
			continue
		}

		(*occupied)[room.name] = true
	}

	// [] braket
	// () paranthesis
	// {} curly braket

	return path
}
