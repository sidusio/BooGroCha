package chalmers

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
