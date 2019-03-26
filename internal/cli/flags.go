package cli

import "github.com/spf13/viper"

var user string

func loadFlags() {
	BgcCmd.PersistentFlags().StringVarP(&user, "cid", "", "", "Manually specify the user")
}

func bindFlags() error {
	err := viper.BindPFlag("chalmers.cid", BgcCmd.PersistentFlags().Lookup("cid"))
	return err
}
