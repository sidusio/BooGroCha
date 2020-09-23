package cli

import (
	"fmt"
	"os"

	"github.com/spf13/viper"

	"sidus.io/boogrocha/internal/booking"
	"sidus.io/boogrocha/internal/booking/directory"
	"sidus.io/boogrocha/internal/booking/timeedit/chalmers"
	"sidus.io/boogrocha/internal/booking/timeedit/chalmers_test"
	logfmt "sidus.io/boogrocha/internal/log/fmt"
)

func getBookingService() booking.BookingService {
	if viper.GetString("chalmers.cid") == "" {
		fmt.Println("No cid specified, set it permanently with 'bgc config set cid' or use the '--cid' flag")
		os.Exit(1)
	}
	chalmersBS, err := chalmers.NewBookingService(viper.GetString("chalmers.cid"), getPassword())
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}

	chalmersTestBS, err := chalmers_test.NewBookingService(viper.GetString("chalmers.cid"), getPassword())
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}

	bs := directory.NewBookingService(map[string]booking.BookingService{
		"TimeEditchalmers":      chalmersBS,
		"TimeEditchalmers_test": chalmersTestBS,
	}, &logfmt.Logger{})

	return bs
}
