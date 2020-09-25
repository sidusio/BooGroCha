package timeedit

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"

	"golang.org/x/net/publicsuffix"

	"github.com/PuerkitoBio/goquery"

	"sidus.io/boogrocha/internal/booking"
)

const samlURLFormat = "https://cloud.timeedit.net/%s/web/timeedit/sso/%s?back=https://cloud.timeedit.net/%s/web/b1/"
const bookURLFormat = "https://cloud.timeedit.net/%s/web/b1/ri1Q5008.html"
const objectsURLFormat = "https://cloud.timeedit.net/%s/web/b1/objects.json?part=t&types=186&step=1"
const bookingsURLFormat = "https://cloud.timeedit.net/%s/web/b1/my.html"
const roomInfoURL = "https://boogrocha.sidus.io/rooms.json"

const studentUnionRoomFilter = "&sid=1010"
const otherPurpose = "203460.192"
const BaseProvider = "TimeEdit"

type BookingService struct {
	client  *http.Client
	rooms   rooms
	version version
}

func (bs BookingService) Book(booking booking.Booking) error {
	bookingURL := fmt.Sprintf(bookURLFormat, bs.version.String())

	formData := url.Values{}
	roomId, err := bs.rooms.idFromName(booking.Room.Id)
	if err != nil {
		return err
	}
	formData.Add("o", roomId)       // Denotes the room
	formData.Add("o", otherPurpose) // Denotes the purpose always "other"
	formData.Add("dates", booking.Start.Format("20060102"))
	formData.Add("starttime", booking.Start.Format("15:04"))
	formData.Add("endtime", booking.End.Format("15:04"))
	formData.Add("fe2", booking.Text)
	formData.Add("fe8", "Booked with BookingDemo") // Todo
	formData.Add("url", bookingURL)
	resp, err := bs.client.PostForm(bookingURL, formData)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to book room (%d)", resp.StatusCode)
		}
		return fmt.Errorf(string(body))
	}
	return nil
}

func (bs BookingService) UnBook(booking booking.Booking) error {
	bookingsURL := fmt.Sprintf(bookingsURLFormat, bs.version.String())

	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s?id=%s", bookingsURL, booking.Id), nil)
	if err != nil {
		return err
	}

	// Fetch Request
	resp, err := bs.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("failed to unbook")
	}

	return nil
}

func (bs BookingService) MyBookings() ([]booking.Booking, error) {
	bookingsURL := fmt.Sprintf(bookingsURLFormat, bs.version.String())
	resp, err := bs.client.Get(bookingsURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	// Find the table with the bookings in
	selections := doc.Find("body #texttable table tr")

	// The important information starts on the third row
	trs := make([]*goquery.Selection, selections.Length()-2)
	selections.Each(func(i int, selection *goquery.Selection) {
		if i >= 2 {
			trs[i-2] = selection
		}
	})
	bookings := make([]booking.Booking, 0, 4)
	selectedDate := ""
	for _, tr := range trs {
		// Check if the row is a date row
		headline := tr.Find(".headline.t")
		if headline.Length() > 0 {
			// If it is a date row we extract the date and move to the next row
			text := headline.Text()
			selectedDate = strings.Split(text, " ")[1]
		} else {
			// If it isn't and we have no selected date somethings wrong
			if selectedDate == "" {
				return nil, fmt.Errorf("parsing failure")
			}
			id, found := tr.Attr("data-id")
			if !found {
				return nil, fmt.Errorf("parsing failure")
			}

			roomInfo := strings.Split(tr.Find(".column0").Text(), ", ")[0]

			startTime, endTime, err := getBookingPeriod(tr, selectedDate)
			if err != nil {
				return nil, err
			}

			text, err := bs.getText(id)
			if err != nil {
				return nil, err
			}

			bookings = append(bookings, booking.Booking{
				Text:  text,
				Start: startTime,
				End:   endTime,
				Room: booking.Room{
					Provider: bs.Provider(),
					Id:       roomInfo,
				},
				Id: id,
			})
		}
	}
	return bookings, nil
}

func (bs BookingService) Available(start time.Time, end time.Time) ([]booking.Room, error) {
	date := start.Format("20060102")
	dates := fmt.Sprintf("%s-%s", date, date)

	startTime := start.Format("15:04")
	endTime := end.Format("15:04")

	rooms, err := bs.getRooms(fmt.Sprintf("dates=%s&starttime=%s&endtime=%s", dates, startTime, endTime))
	if err != nil {
		return nil, err
	}
	var result []booking.Room

	for _, room := range rooms {
		result = append(result, booking.Room{
			Provider: bs.Provider(),
			Id:       room.Name,
			Seats:    room.Seats,
			Campus:   room.Campus,
		})
	}
	return result, nil
}

func (bs BookingService) Provider() string {
	return BaseProvider + bs.version.String()
}

func NewBookingService(cid string, pass string, timeEditVersion version) (BookingService, error) {
	bs, err := login(cid, pass, timeEditVersion)
	if err != nil {
		return BookingService{}, err
	}

	rs, err := bs.getRooms("")
	if err != nil {
		return BookingService{}, err
	}
	bs.rooms = rs

	return bs, nil
}

func (bs BookingService) getText(id string) (string, error) {
	bookingsURL := fmt.Sprintf(bookingsURLFormat, bs.version.String())
	resp, err := bs.client.Get(fmt.Sprintf("%s?step=3&id=%s", bookingsURL, id))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", err
	}
	text := ""
	rows := doc.Find(".detailedResObjects tr")
	rows.Each(func(i int, selection *goquery.Selection) {
		if selection.Find(".columnname").Text() == "Egen text" {
			text = selection.Find(".pr").Text()
		}
	})
	return text, nil
}

// This function gets more information about the rooms, like on which
// campus it is or how many seats it has. This information doesn't
// exists on TimeEdit at the time of writing this so therefore it has been
// collected from chalmers maps. This process requires multiple api calls per room
// has therefore been summarised into a json and is hosted by us.
func (bs BookingService) getRoomInfo(rs rooms) (rooms, error) {
	var roomInfos map[string]struct {
		Seats  int    `json:"seats"`
		Campus string `json:"campus"`
	}

	resp, err := bs.client.Get(roomInfoURL)
	if err != nil {
		fmt.Println("couldn't get room info json")
		return rs, err
	}
	jsonBytes, err := ioutil.ReadAll(resp.Body)
	_ = resp.Body.Close()
	if err != nil {
		fmt.Println("couldn't read response body")
		return rs, err
	}

	err = json.Unmarshal(jsonBytes, &roomInfos)
	if err != nil {
		fmt.Println("couldn't unmarshal room info json")
		return rs, err
	}

	for i, r := range rs {
		rs[i].Seats = roomInfos[r.Name].Seats
		rs[i].Campus = roomInfos[r.Name].Campus
	}

	return rs, nil
}

func (bs BookingService) getRooms(extra string) (rooms, error) {

	objectsURL := fmt.Sprintf(objectsURLFormat, bs.version)

	if extra != "" {
		objectsURL = fmt.Sprintf("%s&%s", objectsURL, extra)
	}

	if bs.version == VersionChalmers {
		extra += studentUnionRoomFilter
	}

	objectsURL += extra

	rs, err := bs.fetchRooms(objectsURL)
	if err != nil {
		return nil, err
	}

	// Since student union room shouldn't be booked on chalmers_test we
	// remove them from this list.
	if bs.version == VersionChalmersTest {
		studentUnionRooms, err := bs.fetchRooms(objectsURL + studentUnionRoomFilter)
		if err != nil {
			return nil, err
		}
		rs = rs.removeMany(studentUnionRooms)
	}

	rs, err = bs.getRoomInfo(rs)
	if err != nil {
		fmt.Println("couldn't get room info")
		return nil, err
	}

	return rs, nil
}

func (bs BookingService) fetchRooms(objectsURL string) (rooms, error) {
	var jsonResponse struct {
		HasMore bool `json:"hasMore"`
		Rooms   []struct {
			Id     string `json:"idAndType"`
			Fields struct {
				Name string `json:"Lokalsignatur"`
			} `json:"fields"`
		} `json:"objects"`
	}

	start := 0
	max := 50

	var rs = rooms{}

	for {
		requestURL := fmt.Sprintf("%s&max=%d&start=%d", objectsURL, max, start)
		resp, err := bs.client.Get(requestURL)
		if err != nil {
			return nil, err
		}

		jsonBytes, err := ioutil.ReadAll(resp.Body)
		_ = resp.Body.Close()
		if err != nil {
			return nil, err
		}

		/*
			The HasMore boolean isn't always reliable
			The API will return "Inga sökresultat" when there are no more rooms
		*/
		if string(jsonBytes) == "\"Inga sökresultat\"" {
			break
		}

		err = json.Unmarshal(jsonBytes, &jsonResponse)
		if err != nil {
			return nil, err
		}

		for _, r := range jsonResponse.Rooms {
			rs = append(rs, room{
				Name: strings.Trim(r.Fields.Name, " "),
				Id:   r.Id,
			})
		}

		if jsonResponse.HasMore {
			start += max
		} else {
			break
		}
	}
	return rs, nil
}

func login(cid string, pass string, timeEditVersion version) (BookingService, error) {
	// Setup http client with a cookie jar
	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		return BookingService{}, err
	}
	client := &http.Client{
		Jar: jar,
	}
	var saml string
	if timeEditVersion == VersionChalmersTest {
		saml = "saml2_test"
	} else {
		saml = "saml2"
	}

	samlURL := fmt.Sprintf(samlURLFormat, timeEditVersion, saml, timeEditVersion)

	// Initiate SAML auth flow
	resp, err := client.Get(samlURL)
	if err != nil {
		return BookingService{}, err
	}

	// Extract login form from request to cover XSS prevention values
	form, err := getForm(resp, "#loginForm")
	_ = resp.Body.Close()
	if err != nil {
		return BookingService{}, err
	}

	// Populate form with username and password
	form.Values.Add("UserName", toUsername(cid))
	form.Values.Add("Password", pass)

	// Submit login form
	resp, err = form.Post(client)
	if err != nil {
		return BookingService{}, err
	}

	// The IDP responds with a form that redirects to the original site,
	// this form is usually auto submitted by a script snippet but we have to submit it ourselves
	form, err = getForm(resp, "form")
	_ = resp.Body.Close()
	if err != nil {
		return BookingService{}, err
	}

	// Check if login was successful
	success := false
	for key := range form.Values {
		if key == "SAMLResponse" {
			success = true
			break
		}
	}
	if !success {
		return BookingService{}, fmt.Errorf("failed to login")
	}

	// Submit the redirect form
	resp, err = form.Post(client)
	if err != nil {
		return BookingService{}, err
	}
	_ = resp.Body.Close()

	// Check that we got the auth cookie
	u, err := url.Parse(form.Action)
	if err != nil {
		return BookingService{}, err
	}
	success = false
	for _, cookie := range jar.Cookies(u) {
		if cookie.Name == fmt.Sprintf("TE%sweb", timeEditVersion) {
			success = true
			break
		}
	}
	if !success {
		return BookingService{}, fmt.Errorf("failed to retrive cookie")
	}

	return BookingService{
		client:  client,
		version: timeEditVersion,
	}, nil
}

func printCookies(jar http.CookieJar, u string) {
	ur, err := url.Parse(u)
	if err != nil {
		log.Fatal(err)
	}
	for _, cookie := range jar.Cookies(ur) {
		fmt.Printf("  %s: %s\n", cookie.Name, cookie.Value)
	}
}

func toUsername(cid string) string {
	if strings.Contains(cid, "@net.chalmers.se") {
		return cid
	}
	return cid + "@net.chalmers.se"
}

func getBookingPeriod(tr *goquery.Selection, selectedDate string) (time.Time, time.Time, error) {
	timeInfo := tr.Find(".time").Text()

	timeStrings := strings.Split(timeInfo, " - ")
	startTime, err := time.Parse("2006-01-02T15:04", fmt.Sprintf("%sT%s", selectedDate, timeStrings[0]))
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("parsing failure")
	}
	endTime, err := time.Parse("2006-01-02T15:04", fmt.Sprintf("%sT%s", selectedDate, timeStrings[1]))
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("parsing failure")
	}
	return startTime, endTime, nil
}
