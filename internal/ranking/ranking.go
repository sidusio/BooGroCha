package ranking

import (
	"sidus.io/boogrocha/internal/booking"
	"sort"
	"strings"
)

/**
Ranking representation, a low rank is good
*/
type Rankings map[booking.AvailableRoom]uint64

func (r Rankings) Sort(rooms []booking.AvailableRoom) []booking.AvailableRoom {
	sort.Slice(rooms, func(i, j int) bool {
		if r[rooms[i]]-r[rooms[j]] != 0 {
			return r[rooms[i]] < r[rooms[j]]
		} else if strings.Compare(rooms[i].Provider, rooms[j].Provider) != 0 {
			return strings.Compare(rooms[i].Provider, rooms[j].Provider) > 0
		} else {
			return strings.Compare(rooms[i].Name, rooms[j].Name) > 0
		}
	})
	return rooms
}

func (r Rankings) Update(selected booking.AvailableRoom, pool []booking.AvailableRoom) {
	for _, room := range pool {
		if room != selected {
			diff := uint64(0)
			if r[room] > r[selected] {
				diff = 1
			} else {
				diff = 5
			}
			if uint64(diff+r[room]) < r[room] { // Handle overflow
				r.Normalize(1000)
			}
			r[room] = diff + r[room]
		}
	}
}

func (r Rankings) Normalize(amount uint64) {
	for key := range r {
		if r[key] < uint64(r[key]-amount) {
			r[key] = 0
		} else {
			r[key] = r[key] - amount
		}
	}
}
