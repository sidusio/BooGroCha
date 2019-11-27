package booking

import (
	"fmt"
	"time"
)

type MockErrorService struct{}

func (*MockErrorService) Book(booking ServiceBooking) (string, error) {
	return "", fmt.Errorf("mock error")
}

func (*MockErrorService) UnBook(ServiceBooking ServiceBooking) error {
	return fmt.Errorf("mock error")
}

func (*MockErrorService) MyBookings() ([]ServiceBooking, error) {
	return nil, fmt.Errorf("mock error")
}

func (*MockErrorService) Available(start time.Time, end time.Time) ([]string, error) {
	return nil, fmt.Errorf("mock error")
}

type MockStaticService struct {
	bookings []ServiceBooking
	rooms    []string
}

func NewMockStaticService(bookings []ServiceBooking, rooms []string) *MockStaticService {
	return &MockStaticService{bookings: bookings, rooms: rooms}
}

func (*MockStaticService) Book(booking ServiceBooking) (string, error) {
	return "1", nil
}

func (*MockStaticService) UnBook(booking ServiceBooking) error {
	return nil
}

func (ms *MockStaticService) MyBookings() ([]ServiceBooking, error) {
	return ms.bookings, nil
}

func (ms *MockStaticService) Available(start time.Time, end time.Time) ([]string, error) {
	return ms.rooms, nil
}

type MockService struct {
	Bookings map[string]*ServiceBooking
	Rooms    []string
}

func NewMockService(rooms []string) *MockService {
	return &MockService{Bookings: make(map[string]*ServiceBooking), Rooms: rooms}
}

func (bs *MockService) Book(b ServiceBooking) (string, error) {
	if bs.Bookings[b.Room] != nil {
		return "", fmt.Errorf("room already booked")
	}

	bs.Bookings[b.Room] = &b
	return b.Room, nil
}

func (bs *MockService) UnBook(b ServiceBooking) error {
	if bs.Bookings[b.Room] == nil {
		return fmt.Errorf("room not booked")
	}

	delete(bs.Bookings, b.Room)
	return nil
}

func (bs *MockService) MyBookings() ([]ServiceBooking, error) {
	var bookings []ServiceBooking
	for _, v := range bs.Bookings {
		bookings = append(bookings, *v)
	}
	return bookings, nil
}

func (bs *MockService) Available(start time.Time, end time.Time) ([]string, error) {
	var available []string
	for _, room := range bs.Rooms {
		if bs.Bookings[room] == nil {
			available = append(available, room)
		}
	}
	return available, nil
}
