// TODO: find  all path we can take from start to end
package lemin



func FindPaths(antFarm *Antsfarm) []*Path {
	startRoom := antFarm.startroom
	var paths []*Path
	var occupied map[string]bool
	for _, room := range startRoom.links {
		path := Path{rooms: BFS(room, antFarm.endroom, &occupied, antFarm)}
		if path.rooms != nil {
			paths = append(paths, &path)
		}
	}
	paths = SortPaths(paths)
	return paths
}



func UpdatePath(path *Path) {
	for _, ant := range path.ants {
		if ant.position >= len(ant.path.rooms) {
			continue
		}
		ant.Move()
		DisplayMovement(ant.id, ant.path.rooms[ant.position].name)
	}
}
