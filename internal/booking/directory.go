package booking

import (
	"fmt"
	"sync"
	"time"

	"sidus.io/boogrocha/internal/log"
)

const (
	prefixFormat          = "%s/%s"
	aggregatorServiceName = "aggregator"
)

type Directory struct {
	providers map[string]BookingService
	log       log.Logger
}

func NewDirectory(services map[string]BookingService, log log.Logger) *Directory {
	return &Directory{providers: services, log: log}
}

type availableResult struct {
	available []AvailableRoom
	err       *ServiceError
}

type myBookingsResult struct {
	bookings []Booking
	err      *ServiceError
}

type ServiceError struct {
	ServiceName string
	Err         error
}

func (e *ServiceError) Error() string {
	return fmt.Sprintf("couldn't get available rooms from provider %s: %s", e.ServiceName, e.Err.Error())
}

func (bs *Directory) Book(b Booking) (string, error) {
	if len(bs.providers) == 0 {
		err := ErrNoServices
		bs.log.Error(err.Error())
		return "", err
	}

	p := b.Provider
	if bs.providers[p] == nil {
		return "", fmt.Errorf("booking service not found: %s", p)
	}

	return bs.providers[p].Book(b.ServiceBooking)
}

func (bs *Directory) UnBook(b Booking) error {
	if len(bs.providers) == 0 {
		return ErrNoServices
	}

	p := b.Provider
	if bs.providers[p] == nil {
		return fmt.Errorf("booking provider not found: %s", p)
	}

	return bs.providers[p].UnBook(b.ServiceBooking)
}

func (bs *Directory) MyBookings() ([]Booking, []*ServiceError) {
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

func (bs *Directory) myBookings() ([]Booking, []*ServiceError) {
	incoming := make(chan myBookingsResult)

	wg := sync.WaitGroup{}
	go func() {
		wg.Wait()
		close(incoming)
	}()

	for name, provider := range bs.providers {
		wg.Add(1)
		go func(name string, provider BookingService) {
			servicebookings, err := provider.MyBookings()
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
			var bookings []Booking
			for _, sb := range servicebookings {
				bookings = append(bookings, Booking{
					ServiceBooking: sb,
					Provider:       name,
				})
			}
			incoming <- myBookingsResult{
				bookings: bookings,
				err:      nil,
			}
		}(name, provider)
	}

	var bookings []Booking
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

func (bs *Directory) Available(start time.Time, end time.Time) ([]AvailableRoom, []*ServiceError) {
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

func (bs *Directory) available(start time.Time, end time.Time) ([]AvailableRoom, []*ServiceError) {
	incoming := make(chan availableResult)

	wg := sync.WaitGroup{}
	go func() {
		wg.Wait()
		close(incoming)
	}()

	for name, provider := range bs.providers {
		wg.Add(1)
		go func(name string, provider BookingService) {
			roomNames, err := provider.Available(start, end)
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
			var available []AvailableRoom
			for _, roomName := range roomNames {
				available = append(available, AvailableRoom{
					Provider: name,
					Name:     roomName,
				})
			}
			incoming <- availableResult{
				available: available,
				err:       nil,
			}
		}(name, provider)
	}

	var rooms []AvailableRoom
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
