package diectory

import (
	"fmt"
	"sync"
	"time"

	"sidus.io/boogrocha/internal/log"

	"sidus.io/boogrocha/internal/booking"
)

const (
	prefixFormat = "%s/%s"
)

type BookingService struct {
	providers map[string]booking.BookingService
	log       log.Logger
}

func NewBookingService(services map[string]booking.BookingService, log log.Logger) *BookingService {
	return &BookingService{providers: services, log: log}
}

type availableResult struct {
	available []booking.Room
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
	return fmt.Sprintf("couldn't get available rooms from provider %s: %s", e.serviceName, e.err.Error())
}

func (bs *BookingService) Book(b booking.Booking) error {
	if len(bs.providers) == 0 {
		err := ErrNoServices
		bs.log.Error(err.Error())
		return err
	}

	p := b.Room.Provider
	if bs.providers[p] == nil {
		return fmt.Errorf("booking service not found: %s", p)
	}

	return bs.providers[p].Book(b)
}

func (bs *BookingService) UnBook(b booking.Booking) error {
	if len(bs.providers) == 0 {
		return ErrNoServices
	}

	p := b.Room.Provider
	if bs.providers[p] == nil {
		return fmt.Errorf("booking provider not found: %s", p)
	}

	return bs.providers[p].UnBook(b)
}

func (bs *BookingService) MyBookings() ([]booking.Booking, error) {
	if len(bs.providers) == 0 {
		return nil, ErrNoServices
	}

	rooms, errs := bs.myBookings()
	for _, err := range errs {
		bs.log.Error(err.Error())
	}

	if len(errs) == len(bs.providers) {
		return nil, ErrAllServicesFailed
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

	for name, provider := range bs.providers {
		wg.Add(1)
		go func(name string, provider booking.BookingService) {
			bookings, err := provider.MyBookings()
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
			incoming <- myBookingsResult{
				bookings: bookings,
				err:      nil,
			}
		}(name, provider)
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

func (bs *BookingService) Available(start time.Time, end time.Time) ([]booking.Room, error) {
	if len(bs.providers) == 0 {
		return nil, ErrNoServices
	}

	rooms, errs := bs.available(start, end)
	for _, err := range errs {
		bs.log.Error(err.Error())
	}

	if len(errs) == len(bs.providers) {
		return nil, ErrAllServicesFailed
	}

	return rooms, nil
}

func (bs *BookingService) available(start time.Time, end time.Time) ([]booking.Room, []*serviceError) {
	incoming := make(chan availableResult)

	wg := sync.WaitGroup{}
	go func() {
		wg.Wait()
		close(incoming)
	}()

	for name, provider := range bs.providers {
		wg.Add(1)
		go func(name string, provider booking.BookingService) {
			a, err := provider.Available(start, end)
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
			incoming <- availableResult{
				available: a,
				err:       nil,
			}
		}(name, provider)
	}

	var rooms []booking.Room
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
