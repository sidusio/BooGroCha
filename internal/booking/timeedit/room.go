package timeedit

import "fmt"

type room struct {
	Name   string `json:"fields.Lokalsignatur"`
	Id     string `json:"idAndType"`
	Seats  int    `json:"seats"`
	Campus string `json:"campus"`
}

type rooms []room

func (rs rooms) idFromName(name string) (string, error) {
	for _, room := range rs {
		if room.Name == name {
			return room.Id, nil
		}
	}
	return "", fmt.Errorf("no such room")
}

func (rs rooms) nameFromId(id string) (string, error) {
	for _, room := range rs {
		if room.Id == id {
			return room.Name, nil
		}
	}
	return "", fmt.Errorf("no such room")
}

func (rs rooms) remove(i int) rooms {
	return append(rs[:i], rs[i+1:]...)
}

func (rs rooms) removeWithName(name string) rooms {
	for i, r := range rs {
		if r.Name == name {
			return rs.remove(i)
		}
	}
	return rs
}

func (rs rooms) removeWithNames(names []string) rooms {
	var rooms = rs
	for _, n := range names {
		rooms = rooms.removeWithName(n)
	}
	return rooms
}

func (rs rooms) keepWithNames(names []string) rooms {
	var rooms = rooms{}
	for _, r := range rs {
		for _, n := range names {
			if r.Name == n {
				rooms = append(rooms, r)
			}
		}
	}
	return rooms
}
