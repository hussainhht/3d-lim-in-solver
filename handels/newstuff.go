package lemin

func MinRounds(totalAnts int, pathLengths []int) int {
	maxLen := 0
	for _, l := range pathLengths {
		if l > maxLen {
			maxLen = l
		}
	}
	// left start at max path length; right is a rough upper bound
	left, right := maxLen, maxLen+totalAnts*len(pathLengths)

	for left < right {
		mid := (left + right) / 2
		antsDelivered := 0
		for _, L := range pathLengths {
			delivered := mid - L + 1
			if delivered > 0 {
				antsDelivered += delivered
			}
		}
		if antsDelivered >= totalAnts {
			right = mid
		} else {
			left = mid + 1
		}
	}
	return left
}


func PathLengths(paths []*Path) []int {
	lengths := make([]int, len(paths))
	for i, p := range paths {
		if p == nil || len(p.Rooms) == 0 {
			lengths[i] = 0
		} else {
			lengths[i] = len(p.Rooms)
		}
	}
	return lengths
}


