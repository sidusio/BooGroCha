package directory

import (
	"fmt"
	"sync"
	"time"

	"sidus.io/boogrocha/internal/log"

	"sidus.io/boogrocha/internal/booking"
)

const (
	prefixFormat          = "%s/%s"
	aggregatorServiceName = "aggregator"
)

type BookingAggregator struct {
	providers map[string]booking.BookingService
	log       log.Logger
}

func NewBookingService(services map[string]booking.BookingService, log log.Logger) *BookingAggregator {
	return &BookingAggregator{providers: services, log: log}
}

type availableResult struct {
	available []booking.Room
	err       *ServiceError
}

type myBookingsResult struct {
	bookings []booking.Booking
	err      *ServiceError
}

type ServiceError struct {
	ServiceName string
	Err         error
}

func (e *ServiceError) Error() string {
	return fmt.Sprintf("couldn't get available rooms from provider %s: %s", e.ServiceName, e.Err.Error())
}

func (bs *BookingAggregator) Book(b booking.Booking) (string, error) {
	if len(bs.providers) == 0 {
		err := ErrNoServices
		bs.log.Error(err.Error())
		return "", err
	}

	p := b.Room.Provider
	if bs.providers[p] == nil {
		return "", fmt.Errorf("booking service not found: %s", p)
	}

	return bs.providers[p].Book(b)
}

func (bs *BookingAggregator) UnBook(b booking.Booking) error {
	if len(bs.providers) == 0 {
		return ErrNoServices
	}

	p := b.Room.Provider
	if bs.providers[p] == nil {
		return fmt.Errorf("booking provider not found: %s", p)
	}

	return bs.providers[p].UnBook(b)
}

func (bs *BookingAggregator) MyBookings() ([]booking.Booking, []*ServiceError) {
	if len(bs.providers) == 0 {
		return nil, []*ServiceError{
			{
				ServiceName: aggregatorServiceName,
				Err:         ErrNoServices,
			},
		}
	}

	rooms, errs := bs.myBookings()
	for _, err := range errs {
		bs.log.Error(err.Error())
	}

	return rooms, errs
}

func (bs *BookingAggregator) myBookings() ([]booking.Booking, []*ServiceError) {
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
					err: &ServiceError{
						ServiceName: name,
						Err:         err,
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
	var errors []*ServiceError
	for result := range incoming {
		wg.Done()
		if result.err != nil {
			errors = append(errors, result.err)
		}
		bookings = append(bookings, result.bookings...)
	}
	return bookings, errors
}

func (bs *BookingAggregator) Available(start time.Time, end time.Time) ([]booking.Room, []*ServiceError) {
	if len(bs.providers) == 0 {
		return nil, []*ServiceError{
			{
				ServiceName: aggregatorServiceName,
				Err:         ErrNoServices,
			},
		}
	}

	rooms, errs := bs.available(start, end)
	for _, err := range errs {
		bs.log.Error(err.Error())
	}

	return rooms, errs
}

func (bs *BookingAggregator) available(start time.Time, end time.Time) ([]booking.Room, []*ServiceError) {
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
					err: &ServiceError{
						ServiceName: name,
						Err:         err,
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
	var errors []*ServiceError
	for result := range incoming {
		wg.Done()
		if result.err != nil {
			errors = append(errors, result.err)
		}
		rooms = append(rooms, result.available...)
	}
	return rooms, errors
}
