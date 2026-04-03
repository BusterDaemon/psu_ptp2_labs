package config

import (
	"encoding/json"
	"os"
)

type GoodOldConfig struct {
	dbPath       string
	initFilePath string
}

func NewGoodOldConfig(cf_path string) (Configer, error) {
	var cf GoodOldConfig
	err := cf.ReadConfig(cf_path)
	if err != nil {
		return nil, err
	}

	return &cf, nil
}

func (conf *GoodOldConfig) ReadConfig(path string) error {
	type config struct {
		DataBaseFilePath string `json:"DataBaseFilePath"`
		InitFilePath     string `json:"InitFilePath"`
	}
	var cf config

	f, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	err = json.Unmarshal(f, &cf)
	if err != nil {
		return err
	}

	conf.dbPath = cf.DataBaseFilePath
	conf.initFilePath = cf.InitFilePath

	return nil
}

func (cf GoodOldConfig) GetDBPath() string {
	return cf.dbPath
}

func (cf GoodOldConfig) GetInitFilePath() string {
	return cf.initFilePath
}
