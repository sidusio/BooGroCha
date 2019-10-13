package commands

import (
	"testing"
	"time"

	"github.com/magiconair/properties/assert"
	"sidus.io/boogrocha/internal/booking"
)

func TestFormatDateWithWeekdayMonday(t *testing.T) {
	date, _ := time.Parse("Mon Jan 2 15:04:05 -0700 MST 2006", "Mon Oct 14 15:04:05 -0700 MST 2019")
	booking := booking.Booking{Start: date}
	assert.Equal(t, formatDateWithWeekday(booking), "Mon 14/10")
}

func TestFormatDateWithWeekdayTuesday(t *testing.T) {
	date, _ := time.Parse("Mon Jan 2 15:04:05 -0700 MST 2006", "Sun Oct 13 15:04:05 -0700 MST 2019")
	booking := booking.Booking{Start: date}
	assert.Equal(t, formatDateWithWeekday(booking), "Sun 13/10")
}

func TestFormatTime(t *testing.T) {
	start, _ := time.Parse("Mon Jan 2 15:04:05 -0700 MST 2006", "Mon Oct 14 15:04:05 -0700 MST 2019")
	end, _ := time.Parse("Mon Jan 2 15:04:05 -0700 MST 2006", "Mon Oct 14 16:34:05 -0700 MST 2019")
	booking := booking.Booking{Start: start, End: end}
	assert.Equal(t, formatTime(booking), "15:04-16:34")
}
