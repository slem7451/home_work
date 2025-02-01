package calendarconfig

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"                                           //nolint:depguard
	"github.com/slem7451/home_work/hw12_13_14_15_calendar/internal/config" //nolint:depguard
)

type Config struct {
	Logger  LoggerConf `toml:"logger"`
	HTTP    HTTPConf   `toml:"http"`
	GRPC    GRPCConf   `toml:"grpc"`
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

type GRPCConf struct {
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

func (c Config) GetStorage() string {
	return c.Storage
}

func (c Config) GetDB() config.DBConf {
	return c.DB
}

func (db DBConf) DSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s", db.User, db.Password, db.Host, db.Port, db.Name)
}
