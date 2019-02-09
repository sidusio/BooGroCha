package main

import (
	"fmt"
	"github.com/spf13/viper"
)

func loadConfig() error {

	viper.SetConfigName("config")                                   // name of config file (without extension)
	viper.AddConfigPath(fmt.Sprintf("/etc/%s/", ApplicationName))   // path to look for the config file in
	viper.AddConfigPath(fmt.Sprintf("$HOME/.%s/", ApplicationName)) // call multiple times to add many search paths
	viper.AddConfigPath(".")                                        // optionally look for config in the working directory

	err := viper.ReadInConfig() // Find and read the config file
	return err

}
