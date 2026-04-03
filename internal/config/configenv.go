package config

import "os"

type EnvConfig struct {
	dbPath       string
	initFilePath string
}

func (ec *EnvConfig) ReadConfig(_ string) error {
	envDbPath := os.Getenv("API_SQLITE_PATH")
	envInitFilePath := os.Getenv("API_INIT_DB_PATH")

	if envDbPath == "" {
		ec.dbPath = "./store.db"
	}
	ec.initFilePath = envInitFilePath

	return nil
}

func (ec EnvConfig) GetDBPath() string {
	return ec.dbPath
}

func (ec EnvConfig) GetInitFilePath() string {
	return ec.initFilePath
}
