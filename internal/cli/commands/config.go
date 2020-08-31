package commands

import (
	"encoding/base64"
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/crypto/ssh/terminal"
)

func ConfigCmd(getSavePassword func() func(string) error) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Manage settings",
		Long:  `Configure your campus and credentials`,
		Run:   nil,
	}

	cmd.AddCommand(configGetCmd)
	cmd.AddCommand(configSetCmd(getSavePassword))

	return cmd
}

var validGetArgs = []string{"campus", "cid"}
var validSetArgs = append(validGetArgs, "pass")
var validCampuses = []string{"johanneberg", "lindholmen"}

func configSetCmd(getSavePassword func() func(string) error) *cobra.Command {
	return &cobra.Command{
		Use:   fmt.Sprintf("set {%s} {value}", strings.Join(validSetArgs, "|")),
		Short: "Set config option",
		Long: fmt.Sprintf(
			"Set config option.\nValue should not be provided for the pass config option.\nValid campuses are (%s).",
			strings.Join(validCampuses, ", "),
		),
		Run: func(cmd *cobra.Command, args []string) {
			setConfig(cmd, args, getSavePassword)
		},
		ValidArgs: validSetArgs,
		Args: func(cmd *cobra.Command, args []string) error {

			if err := cobra.MinimumNArgs(1)(cmd, args); err != nil {
				return err
			}

			if err := cobra.OnlyValidArgs(cmd, args[:1]); err != nil {
				return err
			}
			if args[0] == "pass" {
				return cobra.ExactArgs(1)(cmd, args)
			}
			return cobra.ExactArgs(2)(cmd, args)
		},
	}
}

func setConfig(cmd *cobra.Command, args []string, getSavePassword func() func(string) error) {
	// Get the value or grab the password
	var value string
	var savePassword func(string) error
	if args[0] == "pass" {
		savePassword = getSavePassword()
		if savePassword == nil {
			fmt.Println("***************************************************")
			fmt.Println("* WARNING: The password won't be securely stored! *")
			fmt.Println("***************************************************")
		}
		value = credentials()
	} else {
		value = args[1]
	}

	// Verify that the campus is valid
	if args[0] == "campus" {
		valid := false
		for _, v := range validCampuses {
			if v == args[1] {
				valid = true
				break
			}
		}
		if !valid {
			fmt.Printf("%s is not a valid campus\n", args[1])
			os.Exit(1)
		}
	}

	if args[0] == "pass" && savePassword != nil {
		err := savePassword(value)
		if err != nil {
			fmt.Printf("Failed to save password: %s\n", err.Error())
			os.Exit(1)
		}
		// Remove any password from config
		value = ""
	}

	if args[0] == "pass" {
		// Base64 encode password for obscurity
		value = base64.StdEncoding.EncodeToString([]byte(value))
	}
	// Set and save config options
	viper.Set(fmt.Sprintf("chalmers.%s", args[0]), value)
	err := viper.WriteConfig()
	if err != nil {
		fmt.Printf("Failed to write to config: %s\n", err.Error())
		os.Exit(1)
	}
}

var configGetCmd = &cobra.Command{
	Use:   fmt.Sprintf("get {%s}", strings.Join(validGetArgs, "|")),
	Short: "Get config option",
	Long:  "Get config option",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(viper.GetString(fmt.Sprintf("chalmers.%s", args[0])))
	},
	Args:      cobra.ExactValidArgs(1),
	ValidArgs: validGetArgs,
}

func credentials() string {
	fmt.Print("Enter Password (default: \"\"): ")
	bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
	fmt.Println()
	if err != nil {
		fmt.Println("failed to read password")
		os.Exit(1)
	}
	return strings.TrimSpace(string(bytePassword))
}
