// TODO: find  all path we can take from start to end
package lemin

import "sort"

func FindPaths(antFarm *Antsfarm) []*Path {
	startRoom := antFarm.StartRoom
	var paths []*Path
	occupied := make(map[string]bool)
	occupied[startRoom.Name] = true
	for _, room := range startRoom.Links {
		path := Path{Rooms: BFS(room, antFarm.EndRoom, &occupied, antFarm)}
		if path.Rooms != nil {
			paths = append(paths, &path)
		}
	}
	
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

func UpdatePath(antFarm *Antsfarm, path *Path) {
	for i := len(path.Ants) - 1; i >= 0; i-- {
		ant := path.Ants[i]
		finished := ant.Move(antFarm)
		DisplayMovement(ant.ID, ant.Path.Rooms[ant.Position].Name)
		if finished {
			// safely remove
			path.Ants = append(path.Ants[:i], path.Ants[i+1:]...)
		}
	}
}
