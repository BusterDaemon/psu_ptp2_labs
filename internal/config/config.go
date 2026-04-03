package config

type Configer interface {
	ReadConfig(path string) error
	GetDBPath() string
	GetInitFilePath() string
}
