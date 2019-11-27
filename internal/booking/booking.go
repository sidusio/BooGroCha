package booking

import (
	"time"
)

type Booking struct {
	ServiceBooking
	Provider string
}

type ServiceBooking struct {
	Room     string
	Start    time.Time
	End      time.Time
	Text     string
	Id       string
}