package cli

import (
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"sidus.io/boogrocha/internal/cli/commands"
)

var RootCmd = &cobra.Command{
	Use:   "bgc",
	Short: "Manage your group room bookings at Chalmers",
	Long: `A lightweight, easy to use application for managing your
		   group room bookings at Chalmers University of Technology`,
	Run: nil,
}

const ApplicationName = "BooGroCha"

func init() {

	RootCmd.AddCommand(commands.BookCmd(getBookingService))
	RootCmd.AddCommand(commands.ConfigCmd())
	RootCmd.AddCommand(commands.DeleteCmd(getBookingService))
	RootCmd.AddCommand(commands.ListCmd(getBookingService))
	RootCmd.AddCommand(commands.VersionCmd())

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

var user string
func loadFlags() {
	RootCmd.PersistentFlags().StringVarP(&user, "cid", "", "", "Manually specify the user")
}

func bindFlags() error {
	err := viper.BindPFlag("chalmers.cid", RootCmd.PersistentFlags().Lookup("cid"))
	return err
}

func loadConfig() error {

	home, err := homedir.Dir()
	if err != nil {
		return err
	}
	configPath := fmt.Sprintf("%s/.%s/", home, ApplicationName)

	viper.SetConfigName("config")   // name of config file (without extension)
	viper.AddConfigPath(configPath) // call multiple times to add many search paths

	viper.SetDefault("chalmers.cid", "")
	viper.SetDefault("chalmers.pass", "")
	viper.SetDefault("chalmers.campus", "johanneberg")

	// Create config folder
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		err = os.MkdirAll(configPath, 0744)
		if err != nil {
			return err
		}
	}

	// Write config file if it doesn't exists
	if _, err := os.Stat(configPath + "config.toml"); os.IsNotExist(err) {
		err = viper.WriteConfigAs(configPath + "config.toml")
		if err != nil {
			return err
		}
		// Make sure no one else can reed the config file
		err = os.Chmod(configPath+"config.toml", 0600)
		if err != nil {
			return err
		}
	}

	return viper.ReadInConfig() // Find and read the config file

}



