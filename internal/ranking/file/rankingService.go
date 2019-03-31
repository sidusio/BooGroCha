package file

import (
	"encoding/json"
	"github.com/spf13/viper"
	"os"
	"sidus.io/boogrocha/internal/ranking"
)

func NewRankingService(path string) (ranking.RankingService, error) {

	// Create folder if it doesnt exist
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err = os.MkdirAll(path, 0744)
		if err != nil {
			return nil, err
		}
	}

	fullPath := path + "rankings.json"

	// Write rankings file if it doesn't exists
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		file, err := os.Create(fullPath)
		if err != nil {
			return nil, err
		}
		_, err = file.Write([]byte("{}"))
		if err != nil {
			return nil, err
		}
	}

}

type RankingService struct {
}

func (RankingService) GetRankings() (ranking.Rankings, error) {
	panic("implement me")
}

func (RankingService) SaveRankings(rankings ranking.Rankings) error {
	panic("implement me")
}
