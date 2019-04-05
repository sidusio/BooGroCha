package ranking

import (
	"sort"
)

/**
Ranking representation, a low rank is good
*/
type Rankings map[string]uint64

func (r Rankings) Sort(rooms []string) []string {
	sort.Slice(rooms, func(i, j int) bool {
		if r[rooms[i]]-r[rooms[j]] != 0 {
			return r[rooms[i]] < r[rooms[j]]
		} else {
			return rooms[i] > rooms[j]
		}
	})
	return rooms
}

func (r Rankings) Update(selected string, pool []string) {
	for _, room := range pool {
		if room != selected {
			if r[room] > r[selected] {
				r[room] = r[room] + 1
			} else {
				r[room] = r[room] + 5
			}
		}
	}
}
