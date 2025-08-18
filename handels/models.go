package lemin

// Antsfarm represents the parsed lem-in input.
type Antsfarm struct {
	antsnum   int
	rooms     map[string]*Room
	startroom *Room
	endroom   *Room
}

// Room represents a single room node.
type Room struct {
	name    string
	links   []*Room
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

func (ant *Ant) Move() {
	ant.position++
}

type Path struct {
	rooms []*Room
	ants  []*Ant
}
