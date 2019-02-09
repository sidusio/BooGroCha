package main

import (
	"fmt"
	"github.com/spf13/viper"
	"github.com/williamleven/BooGroCha/chalmers"
)

const ApplicationName = "BooGroCha"

func init() {
	err := loadConfig()
	if err != nil {
		fmt.Println("Failed to load config")
		panic(err)
	}
	fmt.Println("Loaded config")
}

func main() {

	_, err := chalmers.NewBookingService(viper.GetString("chalmers.cid"), viper.GetString("chalmers.pass"))
	if err != nil {
		panic(err)
	}

	// TODO some cli?

	fmt.Println("Done")
}
