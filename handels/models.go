package lemin

import "fmt"

// Antsfarm represents the parsed lem-in input.
type Antsfarm struct {
	TotalAnts int
	DoneAnts  int
	NextAnt   int
	Rooms     map[string]*Room
	StartRoom *Room
	EndRoom   *Room
}

// Room represents a single room node.
type Room struct {
	Name    string
	Links   []*Room
	Ant     *Ant
	X, Y    int
	Visited bool
}

func (r *Room) String() string {
	return r.Name
}

type Ant struct {
	ID       int
	Position int //position in path
	Path     *Path
}

func (ant *Ant) Move(antFarm *Antsfarm) bool {
	nextPosition := ant.Position + 1
	if nextPosition >= len(ant.Path.Rooms) {
		// Already at end
		return true
	}

	currentRoom := ant.Path.Rooms[ant.Position]
	nextRoom := ant.Path.Rooms[nextPosition]

	// If next room is occupied, wait
	if nextRoom.Ant != nil {
		fmt.Print("L", ant.ID, " waiting in room ", currentRoom.Name, " for room ", nextRoom.Name, " ")
		return false
	}

	// Move ant
	currentRoom.Ant = nil
	ant.Position++
	nextRoom.Ant = ant

	// Reached end?
	if nextRoom == antFarm.EndRoom {
		antFarm.DoneAnts++
		nextRoom.Ant = nil // free end room
		return true        // signal that this ant is finished
	}

	return false
}

type Path struct {
	Rooms []*Room
	Ants  []*Ant
}

func (p *Path) String() string {
	if p == nil {
		return "<nil path>"
	}
	roomNames := make([]string, len(p.Rooms))
	for i, room := range p.Rooms {
		roomNames[i] = room.Name
	}
	return "Path: " + fmt.Sprintf("%v", roomNames)
}
