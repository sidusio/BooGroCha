package file

import (
	"encoding/json"
	"io/ioutil"
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
	return RankingService{
		path: fullPath,
	}, nil
}

type RankingService struct {
	path string
}

func (rs RankingService) GetRankings() (ranking.Rankings, error) {
	bytes, err := ioutil.ReadFile(rs.path)
	if err != nil {
		return nil, err
	}

	var data ranking.Rankings
	err = json.Unmarshal(bytes, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (rs RankingService) SaveRankings(rankings ranking.Rankings) error {
	data, _ := json.Marshal(rankings)

	err := ioutil.WriteFile(rs.path, data, 0644)
	return err
}
