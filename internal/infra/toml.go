package infra

import (
	"github.com/BurntSushi/toml"
)

type (
	ConfigFile struct {
		Token        string
		BaseUrl      string
		RepositoryId int
		PipelineId   int
	}
)

const (
	_FILENAME string = "config.toml"
)

func LoadConfig() ConfigFile {
	var config ConfigFile
	if _, err := toml.DecodeFile(_FILENAME, &config); err != nil {
		panic(err)
	}
	return config
}
