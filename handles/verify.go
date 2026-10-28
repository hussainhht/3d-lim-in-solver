package lemin

import "fmt"

// Step is one movement of a single ant: at the given turn (1-based) it enters Room.
type Step struct {
	Turn int
	Room *Room
}

// Tracks groups the movements of a Result by ant. The returned slice is indexed
// by ant ID - 1 and each track is ordered by turn. An ant that never moved has
// an empty track. Movements to unknown rooms are skipped.
func (r *Result) Tracks(farm *Antsfarm) [][]Step {
	tracks := make([][]Step, farm.TotalAnts)
	for t, turn := range r.Turns {
		for _, m := range turn {
			room, ok := farm.Rooms[m.Room]
			if !ok || m.AntID < 1 || m.AntID > farm.TotalAnts {
				continue
			}
			tracks[m.AntID-1] = append(tracks[m.AntID-1], Step{Turn: t + 1, Room: room})
		}
	}
	return tracks
}

// Verify checks a simulation result against the lem-in rules and returns the
// first violation found, or nil when the result is a legal solution:
//   - every ant starts in the start room and moves at most once per turn
//   - every move follows an existing tunnel
//   - a tunnel carries at most one ant per turn
//   - no room other than start/end holds more than one ant at the end of a turn
//   - every ant ends in the end room
func Verify(farm *Antsfarm, res *Result) error {
	position := make(map[int]*Room, farm.TotalAnts)
	for id := 1; id <= farm.TotalAnts; id++ {
		position[id] = farm.StartRoom
	}

	for t, turn := range res.Turns {
		moved := make(map[int]bool, len(turn))
		tunnels := make(map[[2]string]bool, len(turn))

		for _, m := range turn {
			from, known := position[m.AntID]
			if !known {
				return fmt.Errorf("turn %d: unknown ant L%d", t+1, m.AntID)
			}
			to, ok := farm.Rooms[m.Room]
			if !ok {
				return fmt.Errorf("turn %d: L%d moves to unknown room %q", t+1, m.AntID, m.Room)
			}
			if moved[m.AntID] {
				return fmt.Errorf("turn %d: L%d moves twice", t+1, m.AntID)
			}
			moved[m.AntID] = true

			linked := false
			for _, l := range from.Links {
				if l == to {
					linked = true
					break
				}
			}
			if !linked {
				return fmt.Errorf("turn %d: L%d moves %s -> %s but there is no tunnel", t+1, m.AntID, from.Name, to.Name)
			}

			key := [2]string{from.Name, to.Name}
			if from.Name > to.Name {
				key = [2]string{to.Name, from.Name}
			}
			if tunnels[key] {
				return fmt.Errorf("turn %d: tunnel %s-%s used twice", t+1, key[0], key[1])
			}
			tunnels[key] = true

			position[m.AntID] = to
		}

		occupant := make(map[*Room]int)
		for id, room := range position {
			if room == farm.StartRoom || room == farm.EndRoom {
				continue
			}
			if other, taken := occupant[room]; taken {
				return fmt.Errorf("turn %d: room %s holds both L%d and L%d", t+1, room.Name, other, id)
			}
			occupant[room] = id
		}
	}

	for id, room := range position {
		if room != farm.EndRoom {
			return fmt.Errorf("ant L%d finished in room %s instead of the end room", id, room.Name)
		}
	}
	return nil
}
