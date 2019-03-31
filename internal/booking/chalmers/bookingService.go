package chalmers

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/publicsuffix"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"sidus.io/boogrocha/internal/booking"
	"strings"
	"time"
)

const samelURL = "https://se.timeedit.net/web/chalmers/db1/timeedit/sso/saml2?back=https%3A%2F%2Fcloud.timeedit.net%2Fchalmers%2Fweb%2Fb1%2F"
const bookURL = "https://cloud.timeedit.net/chalmers/web/b1/ri1Q5008.html"
const objectsURL = "https://cloud.timeedit.net/chalmers/web/b1/objects.json?part=t&types=186&step=1"
const bookingsURL = "https://cloud.timeedit.net/chalmers/web/b1/my.html"
const otherPurpose = "203460.192"

type BookingService struct {
	client *http.Client
	rooms  rooms
}

func (bs BookingService) Book(booking booking.Booking) error {
	formData := url.Values{}
	roomId, err := bs.rooms.idFromName(booking.Room)
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
	formData.Add("url", bookURL)
	resp, err := bs.client.PostForm(bookURL, formData)
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
	resp, err := bs.client.Get(bookingsURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	selections := doc.Find("body #texttable table tr")

	trs := make([]*goquery.Selection, selections.Length()-2)
	selections.Each(func(i int, selection *goquery.Selection) {
		if i >= 2 {
			trs[i-2] = selection
		}
	})
	bookings := make([]booking.Booking, 0, 4)
	selectedDate := ""
	for _, tr := range trs {
		headline := tr.Find(".headline.changeDateLink")
		if headline.Length() > 0 {
			selectedDate = strings.Split(headline.Text(), " ")[1]
		} else {
			if selectedDate == "" {
				return nil, fmt.Errorf("parsing failure")
			}
			id, found := tr.Attr("data-id")
			if !found {
				return nil, fmt.Errorf("parsing failure")
			}

			timeInfo := tr.Find(".time").Text()
			roomInfo := strings.Split(tr.Find(".column0").Text(), ", ")[0]

			timeStrings := strings.Split(timeInfo, " - ")
			startTime, err := time.Parse("2006-01-02T15:04", fmt.Sprintf("%sT%s", selectedDate, timeStrings[0]))
			if err != nil {
				return nil, fmt.Errorf("parsing failure")
			}
			endTime, err := time.Parse("2006-01-02T15:04", fmt.Sprintf("%sT%s", selectedDate, timeStrings[1]))
			if err != nil {
				return nil, fmt.Errorf("parsing failure")
			}

			text, err := bs.getText(id)
			if err != nil {
				return nil, err
			}

			bookings = append(bookings, booking.Booking{
				Text:  text,
				Start: startTime,
				End:   endTime,
				Room:  roomInfo,
				Id:    id,
			})
		}
	}
	return bookings, nil
}

func (bs BookingService) getText(id string) (string, error) {
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

func (bs BookingService) Available(start time.Time, end time.Time) ([]string, error) {
	date := start.Format("20060102")
	dates := fmt.Sprintf("%s-%s", date, date)

	startTime := start.Format("15:04")
	endTime := end.Format("15:04")

	rooms, err := bs.fetchRooms(fmt.Sprintf("dates=%s&starttime=%s&endtime=%s", dates, startTime, endTime))
	if err != nil {
		return nil, err
	}
	var result []string

	for _, room := range rooms {
		result = append(result, room.Name)
	}
	return result, nil
}

func NewBookingService(cid, pass string) (BookingService, error) {
	// Setup http client with a cookie jar
	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		return BookingService{}, err
	}
	client := &http.Client{
		Jar: jar,
	}

	// Initiate SAML auth flow
	resp, err := client.Get(samelURL)
	if err != nil {
		return BookingService{}, err
	}

	// Extract login form from request to cover XSS prevention values
	form, err := getForm(resp, "form")
	_ = resp.Body.Close()
	if err != nil {
		return BookingService{}, err
	}

	// Populate form with username and password
	form.Values.Add("ctl00$ContentPlaceHolder1$UsernameTextBox", cid)
	form.Values.Add("ctl00$ContentPlaceHolder1$PasswordTextBox", pass)

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
		if cookie.Name == "TEchalmersweb" {
			success = true
			break
		}
	}
	if !success {
		return BookingService{}, fmt.Errorf("failed to retrive cookie")
	}

	bs := BookingService{
		client: client,
	}

	rs, err := bs.fetchRooms("") // TODO
	if err != nil {
		return BookingService{}, err
	}
	bs.rooms = rs

	return bs, nil
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

func (bs BookingService) fetchRooms(extra string) (rooms, error) {
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

	var rs rooms

	for {
		url := fmt.Sprintf("%s&max=%d&start=%d", objectsURL, max, start)
		if extra != "" {
			url = fmt.Sprintf("%s&%s", url, extra)
		}
		resp, err := bs.client.Get(url)
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
				Name: r.Fields.Name,
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
