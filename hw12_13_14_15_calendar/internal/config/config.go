package config

import (
	"os"

	"github.com/BurntSushi/toml" //nolint:depguard
)

type Config struct {
	Logger  LoggerConf `toml:"logger"`
	HTTP    HTTPConf   `toml:"http"`
	DB      DBConf     `toml:"db"`
	Storage string     `toml:"storage"`
}

type LoggerConf struct {
	Level string `toml:"level"`
}

type HTTPConf struct {
	Host string `toml:"host"`
	Port int    `toml:"port"`
}

type DBConf struct {
	Host     string `toml:"host"`
	Port     int    `toml:"port"`
	User     string `toml:"user"`
	Password string `toml:"password"`
	Name     string `toml:"name"`
}

func NewConfig(configFile string) Config {
	data, err := os.ReadFile(configFile)
	if err != nil {
		panic(err)
	}

	conf := Config{}
	if _, err := toml.Decode(string(data), &conf); err != nil {
		panic(err)
	}

	return conf
}
