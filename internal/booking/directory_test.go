package booking

import (
	"fmt"
	"reflect"
	"sort"
	"testing"
	"time"

	"sidus.io/boogrocha/internal/log"
	fmtLog "sidus.io/boogrocha/internal/log/fmt"
)

const (
	providerA = "providerA"
	providerB = "providerB"
	providerC = "providerC"
	providerX = "providerX"
	roomA     = "roomA"
	roomB     = "roomB"
	roomC     = "roomC"
)

func TestBookingService_Book(t *testing.T) {
	type fields struct {
		services map[string]BookingService
		log      log.Logger
	}
	type args struct {
		booking Booking
	}
	var tests = []struct {
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
				booking: Booking{},
			},
			wantErr: true,
		},
		{
			name: "only failing services",
			fields: fields{
				services: map[string]BookingService{
					providerA: &MockErrorService{},
					providerB: &MockErrorService{},
				},
				log: &fmtLog.Logger{},
			},
			args: args{
				booking: Booking{
					ServiceBooking: ServiceBooking{
						Room: roomA,
					},
					Provider: providerA,
				},
			},
			wantErr: true,
		},
		{
			name: "invalid prefix",
			fields: fields{
				services: map[string]BookingService{
					providerA: NewMockService([]string{roomA, roomB, roomC}),
					providerB: &MockErrorService{},
					providerC: NewMockService([]string{roomA, roomB}),
				},
				log: &fmtLog.Logger{},
			},
			args: args{
				booking: Booking{
					ServiceBooking: ServiceBooking{
						Room: roomA,
					},
					Provider: providerX,
				},
			},
			wantErr: true,
		},
		{
			name: "already booked",
			fields: fields{
				services: map[string]BookingService{
					providerA: &MockService{
						Bookings: map[string]*ServiceBooking{
							roomA: {
								Room: roomA,
							},
						},
						Rooms: []string{
							roomA,
							roomB,
							roomC,
						}},
					providerB: &MockErrorService{},
					providerC: NewMockService([]string{
						roomA, roomB,
					}),
				},
				log: &fmtLog.Logger{},
			},
			args: args{
				booking: Booking{
					ServiceBooking: ServiceBooking{
						Room: roomA,
					},
					Provider: providerA,
				},
			},
			wantErr: true,
		},
		{
			name: "successfully book",
			fields: fields{
				services: map[string]BookingService{
					providerA: NewMockService([]string{roomA, roomB, roomC}),
					providerB: &MockErrorService{},
					providerC: NewMockService([]string{roomA, roomB}),
				},
				log: &fmtLog.Logger{},
			},
			args: args{
				booking: Booking{
					ServiceBooking: ServiceBooking{
						Room: roomA,
					},
					Provider: providerA,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bs := Directory{
				providers: tt.fields.services,
				log:       tt.fields.log,
			}
			if _, err := bs.Book(tt.args.booking); (err != nil) != tt.wantErr {
				t.Errorf("BookingService.Book() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBookingService_UnBook(t *testing.T) {
	type fields struct {
		services map[string]BookingService
		log      log.Logger
	}
	type args struct {
		booking Booking
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
				booking: Booking{},
			},
			wantErr: true,
		},
		{
			name: "only failing services",
			fields: fields{
				services: map[string]BookingService{
					providerA: &MockErrorService{},
					providerB: &MockErrorService{},
				},
				log: &fmtLog.Logger{},
			},
			args: args{
				booking: Booking{
					ServiceBooking: ServiceBooking{
						Room: roomA,
					},
					Provider: providerA,
				},
			},
			wantErr: true,
		},
		{
			name: "invalid prefix",
			fields: fields{
				services: map[string]BookingService{
					providerA: NewMockService([]string{roomA, roomB, roomC}),
					providerB: &MockErrorService{},
					providerC: NewMockService([]string{roomA, roomB}),
				},
				log: &fmtLog.Logger{},
			},
			args: args{
				booking: Booking{
					ServiceBooking: ServiceBooking{
						Room: roomA,
					},
					Provider: providerX,
				},
			},
			wantErr: true,
		},
		{
			name: "successfully unbook",
			fields: fields{
				services: map[string]BookingService{
					providerA: &MockService{Bookings: map[string]*ServiceBooking{
						roomA: {
							Room: roomA,
						}},
						Rooms: []string{
							roomB,
							roomA,
							roomC,
						},
					},
					providerB: &MockErrorService{},
					providerC: NewMockService([]string{roomA, roomB}),
				},
				log: &fmtLog.Logger{},
			},
			args: args{
				booking: Booking{
					ServiceBooking: ServiceBooking{
						Room: roomA,
					},
					Provider: providerA,
				},
			},
			wantErr: false,
		},
		{
			name: "not booked",
			fields: fields{
				services: map[string]BookingService{
					providerA: NewMockService([]string{roomA, roomB, roomC}),
					providerB: &MockErrorService{},
					providerC: NewMockService([]string{roomA, roomB}),
				},
				log: &fmtLog.Logger{},
			},
			args: args{
				booking: Booking{
					ServiceBooking: ServiceBooking{
						Room: roomA,
					},
					Provider: providerA,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := Directory{
				providers: tt.fields.services,
				log:       tt.fields.log,
			}
			if err := b.UnBook(tt.args.booking); (err != nil) != tt.wantErr {
				t.Errorf("BookingService.UnBook() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBookingService_MyBookings(t *testing.T) {
	type fields struct {
		services map[string]BookingService
		log      log.Logger
	}
	tests := []struct {
		name    string
		fields  fields
		want    []Booking
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
				services: map[string]BookingService{
					providerA: &MockErrorService{},
					providerB: &MockErrorService{},
				},
				log: &fmtLog.Logger{},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "some failing services",
			fields: fields{
				services: map[string]BookingService{
					providerA: NewMockStaticService([]ServiceBooking{
						{
							Room: roomA,
						},
						{
							Room: roomB,
						},
						{
							Room: roomC,
						},
					}, nil),
					providerB: &MockErrorService{},
					providerC: NewMockStaticService([]ServiceBooking{
						{
							Room: roomA,
						},
						{
							Room: roomB,
						},
					}, nil),
				},
				log: &fmtLog.Logger{},
			},
			want: []Booking{
				{
					ServiceBooking: ServiceBooking{
						Room: roomA,
					},
					Provider: providerA,
				},
				{
					ServiceBooking: ServiceBooking{
						Room: roomB,
					},
					Provider: providerA,
				},
				{
					ServiceBooking: ServiceBooking{
						Room: roomC,
					},
					Provider: providerA,
				},
				{
					ServiceBooking: ServiceBooking{
						Room: roomA,
					},
					Provider: providerC,
				},
				{
					ServiceBooking: ServiceBooking{
						Room: roomB,
					},
					Provider: providerC,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := Directory{
				providers: tt.fields.services,
				log:       tt.fields.log,
			}
			got, err := b.MyBookings()
			if (err != nil) != tt.wantErr {
				t.Errorf("BookingService.MyBookings() error = %v, wantErr %v", err, tt.wantErr)
			}

			sort.Slice(got, func(i, j int) bool {
				if got[i].Provider == got[j].Provider {
					return got[i].Id < got[j].Id
				}
				return got[i].Provider < got[j].Provider
			})

			sort.Slice(tt.want, func(i, j int) bool {
				if tt.want[i].Provider == tt.want[j].Provider {
					return tt.want[i].Id < tt.want[j].Id
				}
				return tt.want[i].Provider < tt.want[j].Provider
			})

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BookingService.MyBookings() = %v, wantRooms %v", got, tt.want)
			}
		})
	}
}

func TestBookingService_myBookings(t *testing.T) {
	var result []Booking
	services := make(map[string]BookingService)
	nErrors := 0
	for i := 3; i >= 0; i-- {
		providerName := fmt.Sprintf("service%d", i)
		if i%7 == 0 {
			services[providerName] = &MockErrorService{}
			nErrors += 1
		} else {
			var bookings []ServiceBooking
			for j := i % 20; j >= 0; j-- {
				roomName := fmt.Sprintf("room%d", j)
				id := fmt.Sprintf("%d", j)
				bookings = append(bookings, ServiceBooking{
					Room: roomName,
					Id:   fmt.Sprintf(prefixFormat, providerName, id),
				})
				result = append(result, Booking{
					ServiceBooking: ServiceBooking{
						Room: roomName,
						Id:   fmt.Sprintf(prefixFormat, providerName, id),
					},
					Provider: providerName,
				})
			}
			services[providerName] = NewMockStaticService(bookings, nil)
		}
	}

	bs := NewDirectory(services, &fmtLog.Logger{})

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

	if !reflect.DeepEqual(bookings, result) {
		t.Errorf("BookingService.myBookings() = %v, wantBookings %v", bookings, result)
	}
}

func TestBookingService_Available(t *testing.T) {
	type fields struct {
		services map[string]BookingService
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
		want    []AvailableRoom
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
				services: map[string]BookingService{
					providerA: &MockErrorService{},
					providerB: &MockErrorService{},
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
				services: map[string]BookingService{
					providerA: NewMockStaticService(nil, []string{roomA}),
					providerB: &MockErrorService{},
					providerC: NewMockStaticService(nil, []string{roomA, roomB}),
				},
				log: &fmtLog.Logger{},
			},
			args: args{
				start: time.Time{},
				end:   time.Time{},
			},
			want: []AvailableRoom{
				{
					Provider: providerA,
					Name:     roomA,
				},
				{
					Provider: providerC,
					Name:     roomA,
				},
				{
					Provider: providerC,
					Name:     roomB,
				}},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bs := &Directory{
				providers: tt.fields.services,
				log:       tt.fields.log,
			}
			got, err := bs.Available(tt.args.start, tt.args.end)
			if (err != nil) != tt.wantErr {
				t.Errorf("BookingService.Available() error = %v, wantErr %v", err, tt.wantErr)
			}
			sort.Slice(got, func(i, j int) bool {
				if got[i].Provider == got[j].Provider {
					return got[i].Name < got[j].Name
				}
				return got[i].Provider < got[j].Provider
			})
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BookingService.Available() = %v, wantRooms %v", got, tt.want)
			}
		})
	}
}

func TestBookingService_available(t *testing.T) {
	var result []AvailableRoom
	services := make(map[string]BookingService)
	nErrors := 0
	for i := 100; i >= 0; i-- {
		providerName := fmt.Sprintf("service%d", i)
		if i%7 == 0 {
			services[providerName] = &MockErrorService{}
			nErrors += 1
		} else {
			var rooms []string
			for j := i % 20; j >= 0; j-- {
				roomName := fmt.Sprintf("room%d", j)
				room := AvailableRoom{
					Provider: providerName,
					Name:     roomName,
				}
				rooms = append(rooms, roomName)
				result = append(result, room)
			}
			services[providerName] = NewMockStaticService(nil, rooms)
		}
	}

	bs := NewDirectory(services, &fmtLog.Logger{})

	rooms, errors := bs.available(time.Time{}, time.Time{})

	if len(errors) != nErrors {
		t.Errorf("BookingService.availabe() len(errors) = %d, nErrors %d", len(errors), nErrors)
		return
	}

	sort.Slice(rooms, func(i, j int) bool {
		if rooms[i].Provider == rooms[j].Provider {
			return rooms[i].Name < rooms[j].Name
		}
		return rooms[i].Provider < rooms[j].Provider
	})

	sort.Slice(result, func(i, j int) bool {
		if result[i].Provider == result[j].Provider {
			return result[i].Name < result[j].Name
		}
		return result[i].Provider < result[j].Provider
	})

	if !reflect.DeepEqual(rooms, result) {
		t.Errorf("BookingService.available() = %v, wantRooms %v", rooms, result)
	}
}
