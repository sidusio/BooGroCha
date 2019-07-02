package diectory

import (
	"fmt"
	"reflect"
	"sort"
	"testing"
	"time"

	"sidus.io/boogrocha/internal/booking"
	"sidus.io/boogrocha/internal/log"
	fmtLog "sidus.io/boogrocha/internal/log/fmt"
)

func TestBookingService_Book(t *testing.T) {
	type fields struct {
		services map[string]booking.BookingService
		log      log.Logger
	}
	type args struct {
		booking booking.Booking
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "no services",
			fields: fields{},
			args: args{
				booking: booking.Booking{},
			},
			wantErr: true,
		},
		{
			name: "only failing services",
			fields: fields{
				services: map[string]booking.BookingService{
					"serviceA": &booking.MockErrorService{},
					"serviceB": &booking.MockErrorService{},
				},
				log: &fmtLog.Logger{},
			},
			args: args{
				booking: booking.Booking{
					Room: "serviceA/room1",
					Id:   "serviceA/room1",
				},
			},
			wantErr: true,
		},
		{
			name: "invalid prefix",
			fields: fields{
				services: map[string]booking.BookingService{
					"serviceA": booking.NewMockService([]string{"room1", "room2", "room3"}),
					"serviceB": &booking.MockErrorService{},
					"serviceC": booking.NewMockService([]string{"room1", "room2"}),
				},
				log: &fmtLog.Logger{},
			},
			args: args{
				booking: booking.Booking{
					Room: "serviceX/room1",
					Id:   "serviceX/room1",
				},
			},
			wantErr: true,
		},
		{
			name: "already booked",
			fields: fields{
				services: map[string]booking.BookingService{
					"serviceA": &booking.MockService{Bookings: map[string]*booking.Booking{
						"room1": {
							Room: "room1",
							Id:   "room1",
						}}, Rooms: []string{"room1", "room2", "room3"}},
					"serviceB": &booking.MockErrorService{},
					"serviceC": booking.NewMockService([]string{"room1", "room2"}),
				},
				log: &fmtLog.Logger{},
			},
			args: args{
				booking: booking.Booking{
					Room: "serviceA/room1",
					Id:   "serviceA/room1",
				},
			},
			wantErr: true,
		},
		{
			name: "successfully book",
			fields: fields{
				services: map[string]booking.BookingService{
					"serviceA": booking.NewMockService([]string{"room1", "room2", "room3"}),
					"serviceB": &booking.MockErrorService{},
					"serviceC": booking.NewMockService([]string{"room1", "room2"}),
				},
				log: &fmtLog.Logger{},
			},
			args: args{
				booking: booking.Booking{
					Room: "serviceA/room1",
					Id:   "serviceA/room1",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bs := BookingService{
				services: tt.fields.services,
				log:      tt.fields.log,
			}
			if err := bs.Book(tt.args.booking); (err != nil) != tt.wantErr {
				t.Errorf("BookingService.Book() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBookingService_UnBook(t *testing.T) {
	type fields struct {
		services map[string]booking.BookingService
		log      log.Logger
	}
	type args struct {
		booking booking.Booking
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "no services",
			fields: fields{
				log: &fmtLog.Logger{},
			},
			args: args{
				booking: booking.Booking{},
			},
			wantErr: true,
		},
		{
			name: "only failing services",
			fields: fields{
				services: map[string]booking.BookingService{
					"serviceA": &booking.MockErrorService{},
					"serviceB": &booking.MockErrorService{},
				},
				log: &fmtLog.Logger{},
			},
			args: args{
				booking: booking.Booking{
					Room: "serviceA/room1",
					Id:   "serviceA/room1",
				},
			},
			wantErr: true,
		},
		{
			name: "invalid prefix",
			fields: fields{
				services: map[string]booking.BookingService{
					"serviceA": booking.NewMockService([]string{"room1", "room2", "room3"}),
					"serviceB": &booking.MockErrorService{},
					"serviceC": booking.NewMockService([]string{"room1", "room2"}),
				},
				log: &fmtLog.Logger{},
			},
			args: args{
				booking: booking.Booking{
					Room: "serviceX/room1",
					Id:   "serviceX/room1",
				},
			},
			wantErr: true,
		},
		{
			name: "successfully unbook",
			fields: fields{
				services: map[string]booking.BookingService{
					"serviceA": &booking.MockService{Bookings: map[string]*booking.Booking{
						"room1": {
							Room: "room1",
							Id:   "room1",
						}}, Rooms: []string{"room1", "room2", "room3"}},
					"serviceB": &booking.MockErrorService{},
					"serviceC": booking.NewMockService([]string{"room1", "room2"}),
				},
				log: &fmtLog.Logger{},
			},
			args: args{
				booking: booking.Booking{
					Room: "serviceA/room1",
					Id:   "serviceA/room1",
				},
			},
			wantErr: false,
		},
		{
			name: "not booked",
			fields: fields{
				services: map[string]booking.BookingService{
					"serviceA": booking.NewMockService([]string{"room1", "room2", "room3"}),
					"serviceB": &booking.MockErrorService{},
					"serviceC": booking.NewMockService([]string{"room1", "room2"}),
				},
				log: &fmtLog.Logger{},
			},
			args: args{
				booking: booking.Booking{
					Room: "serviceA/room1",
					Id:   "serviceA/room1",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := BookingService{
				services: tt.fields.services,
				log:      tt.fields.log,
			}
			if err := b.UnBook(tt.args.booking); (err != nil) != tt.wantErr {
				t.Errorf("BookingService.UnBook() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBookingService_MyBookings(t *testing.T) {
	type fields struct {
		services map[string]booking.BookingService
		log      log.Logger
	}
	tests := []struct {
		name    string
		fields  fields
		want    []booking.Booking
		wantErr bool
	}{
		{
			name: "no services",
			fields: fields{
				log: &fmtLog.Logger{},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "only failing services",
			fields: fields{
				services: map[string]booking.BookingService{
					"serviceA": &booking.MockErrorService{},
					"serviceB": &booking.MockErrorService{},
				},
				log: &fmtLog.Logger{},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "some failing services",
			fields: fields{
				services: map[string]booking.BookingService{
					"serviceA": booking.NewMockStaticService([]booking.Booking{
						{
							Room: "room1",
							Id:   "room1",
						},
						{
							Room: "room2",
							Id:   "room2",
						},
						{
							Room: "room3",
							Id:   "room3",
						},
					}, nil),
					"serviceB": &booking.MockErrorService{},
					"serviceC": booking.NewMockStaticService([]booking.Booking{
						{
							Room: "room1",
							Id:   "room1",
						},
						{
							Room: "room2",
							Id:   "room2",
						},
					}, nil),
				},
				log: &fmtLog.Logger{},
			},
			want: []booking.Booking{
				{
					Room: "serviceA/room1",
					Id:   "serviceA/room1",
				},
				{
					Room: "serviceA/room2",
					Id:   "serviceA/room2",
				},
				{
					Room: "serviceA/room3",
					Id:   "serviceA/room3",
				},
				{
					Room: "serviceC/room1",
					Id:   "serviceC/room1",
				},
				{
					Room: "serviceC/room2",
					Id:   "serviceC/room2",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := BookingService{
				services: tt.fields.services,
				log:      tt.fields.log,
			}
			got, err := b.MyBookings()
			if (err != nil) != tt.wantErr {
				t.Errorf("BookingService.MyBookings() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			sort.Slice(got, func(i, j int) bool {
				return got[i].Id < got[j].Id
			})

			sort.Slice(tt.want, func(i, j int) bool {
				return tt.want[i].Id < tt.want[j].Id
			})

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BookingService.MyBookings() = %v, wantRooms %v", got, tt.want)
			}
		})
	}
}

func TestBookingService_myBookings(t *testing.T) {
	var result []booking.Booking
	services := make(map[string]booking.BookingService)
	nErrors := 0
	for i := 3; i >= 0; i-- {
		serviceName := fmt.Sprintf("service%d", i)
		if i%7 == 0 {
			services[serviceName] = &booking.MockErrorService{}
			nErrors += 1
		} else {
			var bookings []booking.Booking
			for j := i % 20; j >= 0; j-- {
				roomName := fmt.Sprintf("room%d", j)
				id := fmt.Sprintf("%d", j)
				bookings = append(bookings, booking.Booking{
					Room: roomName,
					Id:   id,
				})
				result = append(result, booking.Booking{
					Room: fmt.Sprintf(prefixFormat, serviceName, roomName),
					Id:   fmt.Sprintf(prefixFormat, serviceName, id),
				})
			}
			services[serviceName] = booking.NewMockStaticService(bookings, nil)
		}
	}

	bs := NewBookingService(services, &fmtLog.Logger{})

	bookings, errors := bs.myBookings()

	if len(errors) != nErrors {
		t.Errorf("BookingService.availabe() len(errors) = %d, nErrors %d", len(errors), nErrors)
		return
	}

	sort.Slice(bookings, func(i, j int) bool {
		return bookings[i].Id < bookings[j].Id
	})

	sort.Slice(result, func(i, j int) bool {
		return result[i].Id < result[j].Id
	})
	//sort.Strings(bookings)
	//sort.Strings(result)

	if !reflect.DeepEqual(bookings, result) {
		t.Errorf("BookingService.myBookings() = %v, wantBookings %v", bookings, result)
	}
}

func TestBookingService_Available(t *testing.T) {
	type fields struct {
		services map[string]booking.BookingService
		log      log.Logger
	}
	type args struct {
		start time.Time
		end   time.Time
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []string
		wantErr bool
	}{
		{
			name: "no services",
			fields: fields{
				log: &fmtLog.Logger{},
			},
			args:    args{},
			want:    nil,
			wantErr: true,
		},
		{
			name: "only failing services",
			fields: fields{
				services: map[string]booking.BookingService{
					"serviceA": &booking.MockErrorService{},
					"serviceB": &booking.MockErrorService{},
				},
				log: &fmtLog.Logger{},
			},
			args:    args{},
			want:    nil,
			wantErr: true,
		},
		{
			name: "some failing services",
			fields: fields{
				services: map[string]booking.BookingService{
					"serviceA": booking.NewMockStaticService(nil, []string{"room1"}),
					"serviceB": &booking.MockErrorService{},
					"serviceC": booking.NewMockStaticService(nil, []string{"room1", "room2"}),
				},
				log: &fmtLog.Logger{},
			},
			args: args{
				start: time.Time{},
				end:   time.Time{},
			},
			want:    []string{"serviceA/room1", "serviceC/room1", "serviceC/room2"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bs := &BookingService{
				services: tt.fields.services,
				log:      tt.fields.log,
			}
			got, err := bs.Available(tt.args.start, tt.args.end)
			if (err != nil) != tt.wantErr {
				t.Errorf("BookingService.Available() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			sort.Strings(got)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BookingService.Available() = %v, wantRooms %v", got, tt.want)
			}
		})
	}
}

func TestBookingService_available(t *testing.T) {
	var result []string
	services := make(map[string]booking.BookingService)
	nErrors := 0
	for i := 100; i >= 0; i-- {
		serviceName := fmt.Sprintf("service%d", i)
		if i%7 == 0 {
			services[serviceName] = &booking.MockErrorService{}
			nErrors += 1
		} else {
			var rooms []string
			for j := i % 20; j >= 0; j-- {
				roomName := fmt.Sprintf("room%d", j)
				rooms = append(rooms, roomName)
				result = append(result, fmt.Sprintf(prefixFormat, serviceName, roomName))
			}
			services[serviceName] = booking.NewMockStaticService(nil, rooms)
		}
	}

	bs := NewBookingService(services, &fmtLog.Logger{})

	rooms, errors := bs.available(time.Time{}, time.Time{})

	if len(errors) != nErrors {
		t.Errorf("BookingService.availabe() len(errors) = %d, nErrors %d", len(errors), nErrors)
		return
	}

	sort.Strings(rooms)
	sort.Strings(result)

	if !reflect.DeepEqual(rooms, result) {
		t.Errorf("BookingService.available() = %v, wantRooms %v", rooms, result)
	}
}

func Test_wrapBooking(t *testing.T) {
	type args struct {
		serviceName string
		b           booking.Booking
	}
	tests := []struct {
		name string
		args args
		want booking.Booking
	}{
		{
			name: "wraps booking",
			args: args{
				serviceName: "serviceA",
				b: booking.Booking{
					Room: "room1",
					Id:   "room1",
				},
			},
			want: booking.Booking{
				Room: "serviceA/room1",
				Id:   "serviceA/room1",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := wrapBooking(tt.args.serviceName, tt.args.b); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("wrapBooking() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_unwrapBooking(t *testing.T) {
	type args struct {
		b booking.Booking
	}
	tests := []struct {
		name    string
		args    args
		want    string
		want1   booking.Booking
		wantErr bool
	}{
		{
			name: "room has wrong format",
			args: args{
				b: booking.Booking{
					Room: "room1",
					Id:   "serviceA/room1",
				},
			},
			want:    "",
			want1:   booking.Booking{},
			wantErr: true,
		},
		{
			name: "id has wrong format",
			args: args{
				b: booking.Booking{
					Room: "serviceA/room1",
					Id:   "room1",
				},
			},
			want:    "",
			want1:   booking.Booking{},
			wantErr: true,
		},
		{
			name: "different services for room and id",
			args: args{
				b: booking.Booking{
					Room: "serviceA/room1",
					Id:   "serviceB/room1",
				},
			},
			want:    "",
			want1:   booking.Booking{},
			wantErr: true,
		},
		{
			name: "",
			args: args{
				b: booking.Booking{
					Room: "serviceA/room1",
					Id:   "serviceA/room1",
				},
			},
			want: "serviceA",
			want1: booking.Booking{
				Room: "room1",
				Id:   "room1",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := unwrapBooking(tt.args.b)
			if (err != nil) != tt.wantErr {
				t.Errorf("unwrapBooking() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("unwrapBooking() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("unwrapBooking() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
