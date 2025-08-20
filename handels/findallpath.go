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
	return paths
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
