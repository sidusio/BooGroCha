package unionBuilding

import (
	"net/http"
	"time"

	"sidus.io/boogrocha/internal/booking"
)

const (
	bookURLFormat   = "http://aptus.chs.chalmers.se/AptusPortal/wwwashcommand.aspx?command=book&PanelId=3655&TypeId=18313&GroupId=%d&Date=%s&IntervalId=%d&NextPage"
	listURLFormat   = "http://aptus.chs.chalmers.se/AptusPortal/wwwashbookings.aspx?panelId=3655&weekoffset=%d&type=%d&group=%d&date=%s"
	cancelURLFormat = "http://aptus.chs.chalmers.se/AptusPortal/wwwashcommand.aspx?command=cancel&PanelId=3655&TypeId=%d&GroupId=%d&Date=%s&IntervalId=%d&NextPage"
)

/*
   IDs for rooms and and groups of rooms. All may not be used
   at first but they are stored here for possible future use.
*/
const (
	// Room ID:s which is passed in the query as "GroupId"
	room1GroupID        = RoomID(40625)
	room2GroupID        = RoomID(42943)
	room3GroupID        = RoomID(42944)
	exerciseHallGroupID = RoomID(40626)
	musicRoomGroupID    = RoomID(40627)

	// ID for the type of rooms available for booking
	groupRoomTypeID    = TypeID(18313)
	musicRoomTypeID    = TypeID(18314)
	exerciseHallTypeID = TypeID(18315)
)

var (
	room1 = room{
		roomID: room1GroupID,
		typeID: groupRoomTypeID,
	}
	room2 = room{
		roomID: room2GroupID,
		typeID: groupRoomTypeID,
	}
	room3 = room{
		roomID: room3GroupID,
		typeID: groupRoomTypeID,
	}
)

type RoomID int

type TypeID int

type BookingService struct {
	client *http.Client
	rooms  rooms
}

func NewBookingService(personalIDNumber, pass string) *BookingService {
	panic("implement me")
}

func (*BookingService) Book(booking booking.Booking) error {
	panic("implement me")
}

func (*BookingService) UnBook(booking booking.Booking) error {
	panic("implement me")
}

func (*BookingService) MyBookings() ([]booking.Booking, error) {
	panic("implement me")
}

func (*BookingService) Available(start time.Time, end time.Time) ([]booking.Room, error) {
	panic("implement me")
}
