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
		if len(left[i].rooms) <= len(right[j].rooms) {
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
