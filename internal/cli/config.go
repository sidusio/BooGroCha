package cli

import (
	"fmt"
	"os"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

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
