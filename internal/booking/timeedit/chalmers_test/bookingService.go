package chalmers_test

import (
	"time"

	"sidus.io/boogrocha/internal/booking"
	"sidus.io/boogrocha/internal/booking/timeedit"
)

const timeEditVersion = timeedit.ChalmersTest

type BookingService struct {
	TimeEditService timeedit.BookingService
}

func NewBookingService(cid, pass string) (BookingService, error) {
	bs, err := timeedit.NewBookingService(cid, pass, timeEditVersion)
	if err != nil {
		return BookingService{}, err
	}

	return BookingService{
		TimeEditService: bs,
	}, err
}

func (bs BookingService) Book(booking booking.Booking) error {
	return bs.TimeEditService.Book(booking, timeEditVersion)
}

func (bs BookingService) UnBook(booking booking.Booking) error {
	return bs.TimeEditService.UnBook(booking, timeEditVersion)
}

func (bs BookingService) MyBookings() ([]booking.Booking, error) {
	return bs.TimeEditService.MyBookings(timeEditVersion)
}

func (bs BookingService) Available(start time.Time, end time.Time) ([]booking.Room, error) {
	return bs.TimeEditService.Available(start, end, timeEditVersion)
}
