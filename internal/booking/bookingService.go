package booking

import "time"

type BookingService interface {
	Book(booking Booking) (string, error)
	UnBook(booking Booking) error
	MyBookings() ([]Booking, error)
	Available(start time.Time, end time.Time) ([]Room, error)
}
