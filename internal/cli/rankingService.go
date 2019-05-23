package cli

import (
	"fmt"
	"github.com/mitchellh/go-homedir"
	"os"
	"sidus.io/boogrocha/internal/ranking"
	"sidus.io/boogrocha/internal/ranking/file"
)

func getRankingService() ranking.RankingService {
	home, err := homedir.Dir()
	if err != nil {
		// TODO
		os.Exit(1)
	}
	path := fmt.Sprintf("%s/.%s/", home, ApplicationName)
	rs, err := file.NewRankingService(path)
	if err != nil {
		// TODO
		os.Exit(1)
	}
	return rs
}
