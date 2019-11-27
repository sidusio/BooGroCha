package cli

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
	"sidus.io/boogrocha/internal/booking"
	"sidus.io/boogrocha/internal/booking/chalmers"
	fmtlog "sidus.io/boogrocha/internal/log/fmt"
)

func getBookingService() *booking.Directory {
	if viper.GetString("chalmers.cid") == "" {
		fmt.Println("No cid specified, set it permanently with 'bgc config set cid' or use the '--cid' flag")
		os.Exit(1)
	}
	bs, err := chalmers.NewBookingService(viper.GetString("chalmers.cid"), getPassword())
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}

	ba := booking.NewDirectory(map[string]booking.BookingService{
		"chalmers": bs,
	}, &fmtlog.Logger{})
	return ba
}
