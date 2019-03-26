package main

import (
	b64 "encoding/base64"
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	B "sidus.io/boogrocha"
	"sidus.io/boogrocha/chalmers"
	"golang.org/x/crypto/ssh/terminal"
	"os"
	"strings"
	"syscall"
)

const ApplicationName = "BooGroCha"

func init() {

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

func loadFlags() {
	rootCmd.PersistentFlags().StringVarP(&user, "cid", "", "", "Manually specify the user")
}

func bindFlags() error {
	err := viper.BindPFlag("chalmers.cid", rootCmd.PersistentFlags().Lookup("cid"))
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

var user string

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func getBookingService() B.BookingService {
	bs, err := chalmers.NewBookingService(viper.GetString("chalmers.cid"), getPassword())
	if err != nil {
		fmt.Printf("Failed to login: %s\n", err.Error())
		os.Exit(1)
	}
	return bs
}

func credentials() string {
	fmt.Println("***************************************************")
	fmt.Println("* WARNING: The password won't be securely stored! *")
	fmt.Println("***************************************************")
	fmt.Print("Enter Password: ")
	bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
	fmt.Println()
	if err != nil {
		fmt.Println("failed to read password")
		os.Exit(1)
	}
	encPass := b64.StdEncoding.EncodeToString([]byte(strings.TrimSpace(string(bytePassword))))
	return encPass
}

func getPassword() string {
	password, err := b64.StdEncoding.DecodeString(viper.GetString("chalmers.pass"))
	if err != nil {
		fmt.Printf("Failed to read password: %s\n", err.Error())
		os.Exit(1)
	}
	return string(password)
}

var rootCmd = &cobra.Command{
	Use:   "bgc",
	Short: "Manage your group room bookings at Chalmers",
	Long: `A lightweight, easy to use application for managing your
		   group room bookings at Chalmers University of Technology`,
	Run: nil,
}
