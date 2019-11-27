package booking

import "time"

type BookingService interface {
	Book(booking ServiceBooking) (string, error)
	UnBook(booking ServiceBooking) error
	MyBookings() ([]ServiceBooking, error)
	Available(start time.Time, end time.Time) ([]string, error)
}
