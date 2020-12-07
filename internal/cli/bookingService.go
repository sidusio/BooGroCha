package cli

import (
	"fmt"
	"os"

	"sidus.io/boogrocha/internal/booking/timeedit"

	"github.com/spf13/viper"

	"sidus.io/boogrocha/internal/booking"
	"sidus.io/boogrocha/internal/booking/directory"
	logfmt "sidus.io/boogrocha/internal/log/fmt"
)

func getBookingService() booking.BookingService {
	if viper.GetString("chalmers.cid") == "" {
		fmt.Println("No cid specified, set it permanently with 'bgc config set cid' or use the '--cid' flag")
		os.Exit(1)
	}
	chalmersBS, err := timeedit.NewBookingService(viper.GetString("chalmers.cid"), getPassword(), timeedit.VersionChalmers)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}

	chalmersCovidBS, err := timeedit.NewBookingService(viper.GetString("chalmers.cid"), getPassword(), timeedit.VersionChalmersCovid)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}

	bs := directory.NewBookingService(map[string]booking.BookingService{
		chalmersBS.Provider():     chalmersBS,
		chalmersCovidBS.Provider(): chalmersCovidBS,
	}, &logfmt.Logger{})

	return bs
}
