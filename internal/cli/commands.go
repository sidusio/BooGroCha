package cli

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"sidus.io/boogrocha/internal/cli/commands"
)

const ApplicationName = "BooGroCha"
const Version = "0.1"

func init() {

	BgcCmd.AddCommand(commands.BookCmd(getBookingService))
	BgcCmd.AddCommand(commands.ConfigCmd(getSavePassword))
	BgcCmd.AddCommand(commands.DeleteCmd(getBookingService))
	BgcCmd.AddCommand(commands.ListCmd(getBookingService))
	BgcCmd.AddCommand(commands.VersionCmd(ApplicationName, Version))

	loadFlags()

	err := loadConfig()
	if err != nil {
		fmt.Printf("Failed to load config: %s", err)
		os.Exit(1)
	}

	err = bindFlags()
	if err != nil {
		panic(err)
	}
}

var BgcCmd = &cobra.Command{
	Use:   "bgc",
	Short: "Manage your group room bookings at Chalmers",
	Long: `A lightweight, easy to use application for managing your
		   group room bookings at Chalmers University of Technology`,
	Run: nil,
}
