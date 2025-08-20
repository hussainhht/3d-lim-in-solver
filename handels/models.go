package lemin

import "fmt"

// Antsfarm represents the parsed lem-in input.
type Antsfarm struct {
	totalAnts int
	doneAnts  int
	nextAnt   int
	rooms     map[string]*Room
	startRoom *Room
	endRoom   *Room
}

// Room represents a single room node.
type Room struct {
	name    string
	links   []*Room
	ant     *Ant
	x, y    int
	visited bool
}

func (r *Room) String() string {
	return r.name
}

type Ant struct {
	id       int
	position int //position in path
	path     *Path
}

func (ant *Ant) Move(antFarm *Antsfarm) bool {
	nextPosition := ant.position + 1
	if nextPosition >= len(ant.path.rooms) {
		// Already at end
		return true
	}

	currentRoom := ant.path.rooms[ant.position]
	nextRoom := ant.path.rooms[nextPosition]

	// If next room is occupied, wait
	if nextRoom.ant != nil {
		fmt.Print("L", ant.id, " waiting in room ", currentRoom.name, " for room ", nextRoom.name, " ")
		return false
	}

	// Move ant
	currentRoom.ant = nil
	ant.position++
	nextRoom.ant = ant

	// Reached end?
	if nextRoom == antFarm.endRoom {
		antFarm.doneAnts++
		nextRoom.ant = nil // free end room
		return true        // signal that this ant is finished
	}

	return false
}


type Path struct {
	rooms []*Room
	ants  []*Ant
}

func (p *Path) String() string {
	if p == nil {
		return "<nil path>"
	}
	roomNames := make([]string, len(p.rooms))
	for i, room := range p.rooms {
		roomNames[i] = room.name
	}
	return "Path: " + fmt.Sprintf("%v", roomNames)
}