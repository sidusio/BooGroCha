package cli

import (
	"fmt"
	"os"

	"github.com/spf13/viper"

	"sidus.io/boogrocha/internal/booking"
	"sidus.io/boogrocha/internal/booking/chalmers"
)

func getBookingService() booking.BookingService {
	if viper.GetString("chalmers.cid") == "" {
		fmt.Println("No cid specified, set it permanently with 'bgc config set cid' or use the '--cid' flag")
		os.Exit(1)
	}
	bs, err := chalmers.NewBookingService(viper.GetString("chalmers.cid"), getPassword())
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}
	return bs
}
