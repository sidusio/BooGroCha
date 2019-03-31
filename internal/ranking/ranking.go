package ranking

import "sort"

type Rankings map[string]uint64

func (r Rankings) Sort(rooms []string) []string {
	sort.Slice(rooms, func(i, j int) bool {
		return r[rooms[i]] > r[rooms[j]]
	})
	return rooms
}

func (r Rankings) Update(selected string, pool []string) {
	for _, room := range pool {
		if room != selected {
			if r[room] > r[selected] {
				r[room] = r[room] + 5
			} else {
				r[room] = r[room] + 1
			}
		}
	}
}
