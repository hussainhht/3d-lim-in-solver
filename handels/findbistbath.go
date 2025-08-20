//TODO : use number of ants dont move it .
//TODO : all path are unique with their room

package lemin

func SortPaths(paths []*Path) []*Path {
	if len(paths) <= 1 {
		return paths
	}
	// Merge sort
	left := SortPaths(paths[:len(paths)/2])
	right := SortPaths(paths[len(paths)/2:])
	var result []*Path
	i := 0
	j := 0

	for i < len(left) && j < len(right) {
		if len(left[i].Rooms) <= len(right[j].Rooms) {
			result = append(result, left[i])
			i++
		} else {
			result = append(result, right[j])
			j++
		}
	}
	for i < len(left) {
		result = append(result, left[i])
		i++
	}
	for j < len(right) {
		result = append(result, right[j])
		j++
	}
	return result
}

func AssignAntsToPath(antFarm *Antsfarm, paths []*Path, dist map[*Path]int) {
	for _, path := range paths {
		if dist[path] <= 0 {
			continue // quota used
		}

		firstRoom := path.Rooms[0]
		// Only place an ant if room is free (except start)
		if firstRoom.Ant != nil && firstRoom != antFarm.StartRoom {
			continue
		}

		newAnt := &Ant{
			ID:       antFarm.NextAnt,
			Path:     path,
			Position: 0,
		}
		path.Ants = append([]*Ant{newAnt}, path.Ants...)
		firstRoom.Ant = newAnt
		antFarm.NextAnt++
		dist[path]--
		DisplayMovement(newAnt.ID, firstRoom.Name)

		if antFarm.NextAnt >= antFarm.TotalAnts {
			break
		}
	}
}
