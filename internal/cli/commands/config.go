package commands

import (
	"encoding/base64"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/crypto/ssh/terminal"
	"os"
	"strings"
	"syscall"
)

func ConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Manage settings",
		Long:  `Configure your campus and credentials`,
		Run:   nil,
	}

	cmd.AddCommand(configGetCmd)
	cmd.AddCommand(configSetCmd)

	return cmd
}

var validGetArgs = []string{"campus", "cid"}
var validSetArgs = append(validGetArgs, "pass")
var validCampuses = []string{"johanneberg", "lindholmen"}

var configSetCmd = &cobra.Command{
	Use:   fmt.Sprintf("set {%s} {value}", strings.Join(validSetArgs, "|")),
	Short: "Set config option",
	Long: fmt.Sprintf(
		"Set config option.\nValue should not be provided for the pass config option.\nValid campuses are (%s).",
		strings.Join(validCampuses, ", "),
	),
	Run: func(cmd *cobra.Command, args []string) {

		// Get the value or grab the password
		var value string
		if args[0] == "pass" {
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

		// Set and save config options
		viper.Set(fmt.Sprintf("chalmers.%s", args[0]), value)
		err := viper.WriteConfig()
		if err != nil {
			fmt.Printf("Failed to write to config: %s\n", err.Error())
			os.Exit(1)
		}
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
		} else {
			return cobra.ExactArgs(2)(cmd, args)
		}

	},
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
	encPass := base64.StdEncoding.EncodeToString([]byte(strings.TrimSpace(string(bytePassword))))
	return encPass
}