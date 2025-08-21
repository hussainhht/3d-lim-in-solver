//TODO : use number of ants dont move it .
//TODO : all path are unique with their room

package lemin

// import "fmt"

// ? bader func  merge  sort
// ? func SortPaths(paths []*Path) []*Path {
// 	if len(paths) <= 1 {
// 		return paths
// 	}
// 	// Merge sort
// 	left := SortPaths(paths[:len(paths)/2])
// 	right := SortPaths(paths[len(paths)/2:])
// 	var result []*Path
// 	i := 0
// 	j := 0
// //
// 	for i < len(left) && j < len(right) {
// 		if len(left[i].Rooms) <= len(right[j].Rooms) {
// 			result = append(result, left[i])
// 			i++
// 		} else {
// 			result = append(result, right[j])
// 			j++
// 		}
// 	}
// 	for i < len(left) {
// 		result = append(result, left[i])
// 		i++
// 	}
// 	for j < len(right) {
// 		result = append(result, right[j])
// 		j++
// 	}
// 	return result
// }

func AssignAntsToPath(antFarm *Antsfarm, paths []*Path, dist *map[string]int) {
	for _, path := range paths {
		// fmt.Println(dist)

		if (*dist)[path.String()] <= 0 {
			continue
		}
		// fmt.Println("1")

		if len(path.Rooms) == 1 {
			// fmt.Println("1")

			newAnt := &Ant{
				ID:       antFarm.NextAnt + 1,
				Path:     path,
				Position: 0, // start room
			}
			antFarm.NextAnt++
			(*dist)[path.String()]--
			antFarm.DoneAnts++
			DisplayMovement(newAnt.ID, antFarm.EndRoom.Name)
			continue
		}
		// Only assign ant in start room (Position = 0)
		newAnt := &Ant{
			ID:       antFarm.NextAnt + 1,
			Path:     path,
			Position: 0, // start room
		}
		path.Ants = append([]*Ant{newAnt}, path.Ants...)
		antFarm.NextAnt++
		(*dist)[path.String()]--
		DisplayMovement(newAnt.ID, path.Rooms[0].Name)

		if antFarm.NextAnt >= antFarm.TotalAnts {
			break
		}
	}
}
