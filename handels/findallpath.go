// TODO: find  all path we can take from start to end
package lemin

func FindPaths(antFarm *Antsfarm) []*Path {
	startRoom := antFarm.startRoom
	var paths []*Path
	occupied := make(map[string]bool)
	for _, room := range startRoom.links {
		path := Path{rooms: BFS(room, antFarm.endRoom, &occupied, antFarm)}
		if path.rooms != nil {
			paths = append(paths, &path)
		}
	}
	paths = SortPaths(paths)

	paths = FilterCheckCompare(paths, antFarm.totalAnts)

	return paths
}

// FilterCheckCompare compares paths and only keeps the useful ones.
func FilterCheckCompare(paths []*Path, totalAnts int) []*Path {
	if len(paths) == 0 {
		return paths
	}

	shortest := len(paths[0].rooms)
	var worthwhile []*Path

	for i, p := range paths {
	
		if (len(p.rooms)-shortest) <= (totalAnts/(i+1)) {
			worthwhile = append(worthwhile, p)
		}
	}

	return worthwhile
}


func UpdatePath(antFarm *Antsfarm, path *Path) {
	for i := len(path.ants) - 1; i >= 0; i-- {
		ant := path.ants[i]
		finished := ant.Move(antFarm)
		DisplayMovement(ant.id, ant.path.rooms[ant.position].name)
		if finished {
			// safely remove
			path.ants = append(path.ants[:i], path.ants[i+1:]...)
		}
	}
}