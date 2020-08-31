package filter

import "sidus.io/boogrocha/internal/booking"

type RoomFilter func(booking.Room) bool

func Filter(rooms []booking.Room, filters []RoomFilter) []booking.Room {
	return filter(rooms, func(r booking.Room) bool {
		for _, filter := range filters {
			if !filter(r) {
				return false
			}
		}
		return true
	})
}

func filter(rooms []booking.Room, filter RoomFilter) []booking.Room {
	var filteredRooms []booking.Room
	for _, room := range rooms {
		if filter(room) {
			filteredRooms = append(filteredRooms, room)
		}
	}
	return filteredRooms
}
