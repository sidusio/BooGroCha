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

	chalmersTestBS, err := timeedit.NewBookingService(viper.GetString("chalmers.cid"), getPassword(), timeedit.VersionChalmersTest)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}

	bs := directory.NewBookingService(map[string]booking.BookingService{
		chalmersBS.Provider():     chalmersBS,
		chalmersTestBS.Provider(): chalmersTestBS,
	}, &logfmt.Logger{})

	return bs
}
