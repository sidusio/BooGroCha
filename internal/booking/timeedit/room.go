package timeedit

import "fmt"

type room struct {
	Name   string `json:"fields.Lokalsignatur"`
	Id     string `json:"idAndType"`
	Seats  int    `json:"seats"`
	Campus string `json:"campus"`
}

type Rooms []room

func (rs Rooms) idFromName(name string) (string, error) {
	for _, room := range rs {
		if room.Name == name {
			return room.Id, nil
		}
	}
	return "", fmt.Errorf("no such room")
}

func (rs Rooms) nameFromId(id string) (string, error) {
	for _, room := range rs {
		if room.Id == id {
			return room.Name, nil
		}
	}
	return "", fmt.Errorf("no such room")
}

func (rs Rooms) remove(i int) Rooms {
	return append(rs[:i], rs[i+1:]...)
}

func (rs Rooms) removeWithName(name string) Rooms {
	for i, r := range rs {
		if r.Name == name {
			return rs.remove(i)
		}
	}
	return rs
}

func (rs Rooms) removeWithNames(names []string) Rooms {
	var rooms = rs
	for _, n := range names {
		rooms = rooms.removeWithName(n)
	}
	return rooms
}

func (rs Rooms) keepWithNames(names []string) Rooms {
	var rooms = Rooms{}
	for _, r := range rs {
		for _, n := range names {
			if r.Name == n {
				rooms = append(rooms, r)
			}
		}
	}
	return rooms
}
