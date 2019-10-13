package srv

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/go-chi/chi"

	"sidus.io/boogrocha/internal/booking"
	"sidus.io/boogrocha/internal/booking/directory"

	"sidus.io/boogrocha/internal/booking/chalmers"
	"sidus.io/boogrocha/internal/credentials"
	fmtLog "sidus.io/boogrocha/internal/log/fmt"
)

type server struct {
	credentialsSecret []byte
}

func NewServer(credentialsSecret []byte) *server {
	return &server{credentialsSecret: credentialsSecret}
}

func (s *server) Run(address string) error {
	return http.ListenAndServe(address, s.newRouter())
}

func (s *server) available(w http.ResponseWriter, r *http.Request) {
	creds, ok := r.Context().Value(credentialsContextKey).(credentials.Credentials)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	fromTimeString := r.URL.Query().Get("from")
	toTimeString := r.URL.Query().Get("to")

	if fromTimeString == "" || toTimeString == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	fromTime, err := time.Parse(time.RFC3339, fromTimeString)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	toTime, err := time.Parse(time.RFC3339, toTimeString)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	cbs, err := chalmers.NewBookingService(creds.CID, creds.Password)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	bs := directory.NewBookingService(map[string]booking.BookingService{
		"timeedit": cbs,
	}, &fmtLog.Logger{})

	rooms, errs := bs.Available(fromTime, toTime)

	data, err := json.Marshal(struct {
		Rooms  []booking.Room
		Errors []*directory.ServiceError
	}{
		Rooms:  rooms,
		Errors: errs,
	})

	_, _ = w.Write(data)

}

func (s *server) bookings(w http.ResponseWriter, r *http.Request) {
	creds, ok := r.Context().Value(credentialsContextKey).(credentials.Credentials)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	cbs, err := chalmers.NewBookingService(creds.CID, creds.Password)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	bs := directory.NewBookingService(map[string]booking.BookingService{
		"timeedit": cbs,
	}, &fmtLog.Logger{})

	bookings, errs := bs.MyBookings()

	data, err := json.Marshal(struct {
		Rooms  []booking.Booking
		Errors []*directory.ServiceError
	}{
		Rooms:  bookings,
		Errors: errs,
	})

	_, _ = w.Write(data)
}

func (*server) delete(w http.ResponseWriter, r *http.Request) {
	creds, ok := r.Context().Value(credentialsContextKey).(credentials.Credentials)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	bookingID := chi.URLParam(r, "bookingID")

	cbs, err := chalmers.NewBookingService(creds.CID, creds.Password)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	bs := directory.NewBookingService(map[string]booking.BookingService{
		"timeedit": cbs,
	}, &fmtLog.Logger{})

	err = bs.UnBook(booking.Booking{
		Room: booking.Room{
			Provider: "timeedit",
		},
		Id: bookingID,
	})

	data, err := json.Marshal(struct {
		Error error
	}{
		Error: err,
	})

	_, _ = w.Write(data)
}

func (*server) book(w http.ResponseWriter, r *http.Request) {
	creds, ok := r.Context().Value(credentialsContextKey).(credentials.Credentials)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	data := struct {
		From string
		To   string
		Text string
		Room struct {
			Provider string
			Id       string
		}
	}{}

	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(bytes, &data)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	from, err := time.Parse(time.RFC3339, data.From)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	to, err := time.Parse(time.RFC3339, data.To)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	cbs, err := chalmers.NewBookingService(creds.CID, creds.Password)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	bs := directory.NewBookingService(map[string]booking.BookingService{
		"timeedit": cbs,
	}, &fmtLog.Logger{})

	bookingId, err := bs.Book(booking.Booking{
		Room: booking.Room{
			Provider: data.Room.Provider,
			Id:       data.Room.Id,
		},
		Start: from,
		End:   to,
		Text:  data.Text,
	})

	response, jsonErr := json.Marshal(struct {
		Id    string
		Error error
	}{
		Id:    bookingId,
		Error: err,
	})
	if jsonErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, _ = w.Write(response)

}
