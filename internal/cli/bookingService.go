package cli

import (
	"encoding/base64"
	"fmt"
	"github.com/spf13/viper"
	"os"
	"sidus.io/boogrocha/internal/booking"
	"sidus.io/boogrocha/internal/booking/chalmers"
)

func getBookingService() booking.BookingService {
	bs, err := chalmers.NewBookingService(viper.GetString("chalmers.cid"), getPassword())
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}
	return bs
}

func getPassword() string {
	password, err := base64.StdEncoding.DecodeString(viper.GetString("chalmers.pass"))
	if err != nil {
		fmt.Printf("Failed to read password: %s\n", err.Error())
		os.Exit(1)
	}
	return string(password)
}
