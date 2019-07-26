package unionBuilding

import (
	"errors"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"sidus.io/boogrocha/internal/booking"
)

func extractBookings(body *goquery.Selection) ([]booking.Booking, error) {
	if body.Nodes == nil {
		return []booking.Booking{}, nil
	}

	tds := body.Find("tr > td")

	// The rows are structured in two data-rows followed by one separator
	if len(tds.Nodes)%6 != 0 {
		return nil, errors.New("couldn't extract bookings")
	}

	bookings := make([]booking.Booking, len(tds.Nodes)/6)

	// Each booking is described by 6 consecutive td elements
	bookingSize := 6
	for i := range bookings {
		tdFirst := i * bookingSize
		b, err := extractBooking(tds, tdFirst)
		if err != nil {
			return nil, err
		}
		bookings[i] = b
	}

	return bookings, nil
}

func extractBooking(selection *goquery.Selection, startIndex int) (booking.Booking, error) {
	var b booking.Booking

	// Extract date
	text := selection.Eq(startIndex + 2).Text()
	date, err := extractDate(text)
	if err != nil {
		return b, err
	}

	// Extract room name
	text = selection.Eq(startIndex + 3).Text()
	room, err := extractRoomName(text)
	if err != nil {
		return b, err
	}

	// Extract time
	text = selection.Eq(startIndex + 4).Text()
	start, end, err := extractTime(text)
	if err != nil {
		return b, err
	}

	b.Room = room
	b.Start = date.Add(start)
	b.End = date.Add(end)

	return b, nil
}

func extractDate(text string) (time.Time, error) {
	return time.Parse("02/01/2006", text)
}

func extractRoomName(text string) (booking.Room, error) {
	texts := strings.Split(text, " - ")
	if len(texts) != 2 {
		return booking.Room{}, errors.New("couldn't parse room")
	}

	return booking.Room{
		Provider: providerName,
		Id:       texts[1],
	}, nil
}

func extractTime(text string) (time.Duration, time.Duration, error) {
	startText := strings.Split(text, "-")[0]
	endText := strings.Split(text, "-")[1]
	start, err := time.Parse("15:04", startText)
	if err != nil {
		return 0, 0, errors.New("couldn't parse start time")
	}
	end, err := time.Parse("15:04", endText)
	if err != nil {
		return 0, 0, errors.New("couldn't parse end time")
	}

	return time.Duration(start.Hour()) * time.Hour, time.Duration(end.Hour()) * time.Hour, nil
}
