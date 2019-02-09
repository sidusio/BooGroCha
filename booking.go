package booking_demo

import (
	"time"
)

type Booking struct {
	Room  string
	Start time.Time
	End   time.Time
	Text  string
	Id    string
}
