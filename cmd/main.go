package main

import (
	b64 "encoding/base64"
	"fmt"
	B "github.com/williamleven/BooGroCha"
	"github.com/williamleven/BooGroCha/chalmers"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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

	viper.SetConfigName("config")                                   // name of config file (without extension)
	viper.AddConfigPath(fmt.Sprintf("/etc/%s/", ApplicationName))   // path to look for the config file in
	viper.AddConfigPath(fmt.Sprintf("$HOME/.%s/", ApplicationName)) // call multiple times to add many search paths
	viper.AddConfigPath(".")                                        // optionally look for config in the working directory

	err := viper.ReadInConfig() // Find and read the config file
	return err

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
	Short: "Hugo is a very fast static site generator",
	Long: `A Fast and Flexible Static Site Generator built with
                love by spf13 and friends in Go.
                Complete documentation is available at http://hugo.spf13.com`,
	Run: nil,
}