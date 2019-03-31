package cli

import (
	"encoding/base64"
	"fmt"
	"github.com/spf13/viper"
	"golang.org/x/crypto/ssh/terminal"
	"os"
	"sidus.io/boogrocha/internal/booking"
	"sidus.io/boogrocha/internal/booking/chalmers"
	"strings"
	"syscall"
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

func getPassword() string {
	var password string
	if viper.GetString("chalmers.pass") == "" {
		fmt.Println("No password set in config. You can set it with 'bgc config set pass'")
	}

	// Prompt for password if it's not set in config or the cid flag is specified
	if viper.GetString("chalmers.pass") == "" || BgcCmd.Flag("cid").Value.String() != "" {
		fmt.Print("Enter Password: ")
		bytes, err := terminal.ReadPassword(int(syscall.Stdin))
		fmt.Println()
		if err != nil {
			fmt.Println("failed to read password")
			os.Exit(1)
		}
		password = strings.TrimSpace(string(bytes))
	} else {
		bytes, err := base64.StdEncoding.DecodeString(viper.GetString("chalmers.pass"))
		if err != nil {
			fmt.Printf("Failed to read password: %s\n", err.Error())
			fmt.Printf("Try to reset or unset your password with 'gbc config set pass'\n")
			os.Exit(1)
		}
		password = string(bytes)
	}

	return password
}
