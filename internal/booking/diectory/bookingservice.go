package diectory

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"sidus.io/boogrocha/internal/log"

	"sidus.io/boogrocha/internal/booking"
)

const (
	prefixFormat = "%s/%s"
)

type BookingService struct {
	services map[string]booking.BookingService
	log      log.Logger
}

func NewBookingService(services map[string]booking.BookingService, log log.Logger) *BookingService {
	return &BookingService{services: services, log: log}
}

type availableResult struct {
	available []string
	err       *serviceError
}

type myBookingsResult struct {
	bookings []booking.Booking
	err      *serviceError
}

type serviceError struct {
	serviceName string
	err         error
}

func (e *serviceError) Error() string {
	return fmt.Sprintf("couldn't get available rooms from provider %s: %e", e.serviceName, e.err)
}

func (bs *BookingService) Book(booking booking.Booking) error {
	if len(bs.services) == 0 {
		return errors.New(ErrNoServices)
	}

	serviceName, b, err := unwrapBooking(booking)
	if err != nil {
		return fmt.Errorf("couldn't book booking %v: %e", booking, err)
	}

	if bs.services[serviceName] == nil {
		return fmt.Errorf("booking service not found: %s", serviceName)
	}

	return bs.services[serviceName].Book(b)
}

func (bs *BookingService) UnBook(booking booking.Booking) error {
	if len(bs.services) == 0 {
		return errors.New(ErrNoServices)
	}

	serviceName, b, err := unwrapBooking(booking)
	if err != nil {
		return fmt.Errorf("couldn't book booking %v: %e", booking, err)
	}

	if bs.services[serviceName] == nil {
		return fmt.Errorf("booking service not found: %s", serviceName)
	}

	return bs.services[serviceName].UnBook(b)
}

func (bs *BookingService) MyBookings() ([]booking.Booking, error) {
	if len(bs.services) == 0 {
		return nil, errors.New(ErrNoServices)
	}

	rooms, errs := bs.myBookings()
	for _, err := range errs {
		bs.log.Error(err)
	}

	if len(errs) == len(bs.services) {
		return nil, errors.New(ErrAllServicesFailed)
	}

	return rooms, nil
}

func (bs *BookingService) myBookings() ([]booking.Booking, []*serviceError) {
	incoming := make(chan myBookingsResult)

	wg := sync.WaitGroup{}
	go func() {
		wg.Wait()
		close(incoming)
	}()

	for name, service := range bs.services {
		wg.Add(1)
		go func(name string, service booking.BookingService) {
			bookings, err := service.MyBookings()
			if err != nil {
				incoming <- myBookingsResult{
					bookings: nil,
					err: &serviceError{
						serviceName: name,
						err:         err,
					},
				}
				return
			}
			for i, b := range bookings {
				bookings[i] = wrapBooking(name, b)
			}
			incoming <- myBookingsResult{
				bookings: bookings,
				err:      nil,
			}
		}(name, service)
	}

	var bookings []booking.Booking
	var errors []*serviceError
	for result := range incoming {
		wg.Done()
		if result.err != nil {
			errors = append(errors, result.err)
		}
		bookings = append(bookings, result.bookings...)
	}
	return bookings, errors
}

func (bs *BookingService) Available(start time.Time, end time.Time) ([]string, error) {
	if len(bs.services) == 0 {
		return nil, errors.New(ErrNoServices)
	}

	rooms, errs := bs.available(start, end)
	for _, err := range errs {
		bs.log.Error(err)
	}

	if len(errs) == len(bs.services) {
		return nil, errors.New(ErrAllServicesFailed)
	}

	return rooms, nil
}

func (bs *BookingService) available(start time.Time, end time.Time) ([]string, []*serviceError) {
	incoming := make(chan availableResult)

	wg := sync.WaitGroup{}
	go func() {
		wg.Wait()
		close(incoming)
	}()

	for name, service := range bs.services {
		wg.Add(1)
		go func(name string, service booking.BookingService) {
			a, err := service.Available(start, end)
			if err != nil {
				incoming <- availableResult{
					available: nil,
					err: &serviceError{
						serviceName: name,
						err:         err,
					},
				}
				return
			}
			for i, room := range a {
				a[i] = fmt.Sprintf(prefixFormat, name, room)
			}
			incoming <- availableResult{
				available: a,
				err:       nil,
			}
		}(name, service)
	}

	var rooms []string
	var errors []*serviceError
	for result := range incoming {
		wg.Done()
		if result.err != nil {
			errors = append(errors, result.err)
		}
		rooms = append(rooms, result.available...)
	}
	return rooms, errors
}

func wrapBooking(serviceName string, b booking.Booking) booking.Booking {
	return booking.Booking{
		Room:  fmt.Sprintf(prefixFormat, serviceName, b.Room),
		Start: b.Start,
		End:   b.End,
		Text:  b.Text,
		Id:    fmt.Sprintf(prefixFormat, serviceName, b.Id),
	}
}

func unwrapBooking(b booking.Booking) (string, booking.Booking, error) {
	parts := strings.Split(b.Room, "/")
	if len(parts) == 1 {
		return "", booking.Booking{}, fmt.Errorf("booking not formatted correctly")
	}
	room := parts[1]
	serviceName := parts[0]

	parts = strings.Split(b.Id, "/")
	if len(parts) == 1 {
		return "", booking.Booking{}, fmt.Errorf("booking not formatted correctly")
	}
	id := parts[1]

	if serviceName != parts[0] {
		return "", booking.Booking{}, fmt.Errorf("booking not formatted correctly")
	}

	return serviceName, booking.Booking{
		Room:  room,
		Start: b.Start,
		End:   b.End,
		Text:  b.Text,
		Id:    id,
	}, nil
}
