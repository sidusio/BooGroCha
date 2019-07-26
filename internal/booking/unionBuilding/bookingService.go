package unionBuilding

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/publicsuffix"

	"sidus.io/boogrocha/internal/booking"
)

const (
	providerName = "Kårhuset"

	bookURLFormat   = "http://aptus.chs.chalmers.se/AptusPortal/wwwashcommand.aspx?command=book&PanelId=3655&TypeId=18313&GroupId=%d&Date=%s&IntervalId=%d&NextPage"
	listURL         = "http://aptus.chs.chalmers.se/AptusPortal/wwwashbookings.aspx?"
	cancelURLFormat = "http://aptus.chs.chalmers.se/AptusPortal/wwwashcommand.aspx?command=cancel&PanelId=3655&TypeId=%d&GroupId=%d&Date=%s&IntervalId=%d&NextPage"

	loginURLPOST = "http://aptus.chs.chalmers.se/AptusPortal/login.aspx?ReturnUrl=%2FAptusPortal%2Findex.aspx"
	loginURL     = "http://aptus.chs.chalmers.se/AptusPortal/"

	// Login form keys
	lastFocusKey          = "__LASTFOCUS"
	eventTargetKey        = "__EVENTTARGET"
	eventArgumentKey      = "__EVENTARGUMENT"
	viewStateKey          = "__VIEWSTATE"
	viewStateGeneratorKey = "__VIEWSTATEGENERATOR"
	eventValidationKey    = "__EVENTVALIDATION"
	loginUsernameKey      = "LoginPortal$UserName"
	loginPasswordKey      = "LoginPortal$Password"
	loginButtonKey        = "LoginPortal$LoginButton"
)

/*
   IDs for rooms and and groups of rooms. All may not be used
   at first but they are stored here for possible future use.
*/
const (
	// Room ID:s which is passed in the query as "GroupId"
	room1GroupID = RoomID(40625)
	room2GroupID = RoomID(42943)
	room3GroupID = RoomID(42944)
	//exerciseHallGroupID = RoomID(40626)
	//musicRoomGroupID = RoomID(40627)

	// ID for the type of rooms available for booking
	groupRoomTypeID = TypeID(18313)
	//musicRoomTypeID = TypeID(18314)
	//exerciseHallTypeID = TypeID(18315)
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
}

func NewBookingService(pid, pass string) (BookingService, error) {
	// Setup http client with a cookie jar
	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		return BookingService{}, err
	}
	client := &http.Client{
		Jar: jar,
	}

	// Collect the form data required for login
	values, err := loginForm(client)
	if err != nil {
		return BookingService{}, err
	}

	// Add credentials
	values[loginUsernameKey] = []string{pid}
	values[loginPasswordKey] = []string{pass}

	// Login
	_, err = client.PostForm(loginURLPOST, values)
	if err != nil {
		return BookingService{}, err
	}

	loginurl, err := url.Parse(loginURL)
	fmt.Println(client.Jar.Cookies(loginurl))

	return BookingService{
		client: client,
	}, nil
}

func loginForm(client *http.Client) (url.Values, error) {
	resp, err := client.Get(loginURL)
	if err != nil {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	inputs := doc.Find("body > form > input")

	values := url.Values{}

	/*
	  Traverses through each input field and extracts the name and values
	  which is used as parameters in the login call.
	*/
	inputs.Each(func(i int, s *goquery.Selection) {
		name, exists := s.Attr("name")
		if exists {
			value, exists := s.Attr("value")
			if exists {
				values.Add(name, value)
			}
		}
	})

	for _, key := range []string{
		viewStateKey,
		viewStateGeneratorKey,
		eventValidationKey,
	} {
		if _, ok := values[key]; !ok {
			return nil, errors.New(fmt.Sprintf("value not found in form: %s", key))
		}
	}
	values.Add(loginButtonKey, "Enter")

	return values, nil
}

func (bs *BookingService) Book(booking booking.Booking) error {
	go bs.client.Get(fmt.Sprintf(bookURLFormat,
		room1GroupID,
		booking.Start.Format("2006-01-02"),
		booking.Start.Hour(),
	))
	return nil
}

func (*BookingService) UnBook(booking booking.Booking) error {
	panic("implement me")
}

func (bs *BookingService) MyBookings() ([]booking.Booking, error) {
	resp, err := bs.client.Get(listURL)
	if err != nil {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	bookingsTable := doc.Find("body > table > tbody > tr > td > table > tbody > tr > td > table > tbody > tr:nth-child(3) > td > table > tbody")

	return extractBookings(bookingsTable)
}

func (*BookingService) Available(start time.Time, end time.Time) ([]booking.Room, error) {
	panic("implement me")
}
