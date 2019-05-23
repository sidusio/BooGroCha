package booking

import (
	"time"
)

type Booking struct {
	Room  Room
	Start time.Time
	End   time.Time
	Text  string
	Id    string
}
