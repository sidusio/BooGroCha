package srv

import (
	"encoding/json"
	"net/http"
	"time"

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
