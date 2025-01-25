package senderconfig

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"                                           //nolint:depguard
	"github.com/slem7451/home_work/hw12_13_14_15_calendar/internal/config" //nolint:depguard
)

type Config struct {
	Logger  LoggerConf `toml:"logger"`
	DB      DBConf     `toml:"db"`
	Rabbit  RabbitConf `toml:"rabbit"`
	Storage string     `toml:"storage"`
}

type LoggerConf struct {
	Level string `toml:"level"`
}

type DBConf struct {
	Host     string `toml:"host"`
	Port     int    `toml:"port"`
	User     string `toml:"user"`
	Password string `toml:"password"`
	Name     string `toml:"name"`
}

type RabbitConf struct {
	Host     string `toml:"host"`
	Port     int    `toml:"port"`
	User     string `toml:"user"`
	Password string `toml:"password"`
	Exchange string `toml:"exchange"`
	Queue    string `toml:"queue"`
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

func (r RabbitConf) URL() string {
	return fmt.Sprintf("amqp://%s:%s@%s:%d", r.User, r.Password, r.Host, r.Port)
}

func (r RabbitConf) GetExchange() string {
	return r.Exchange
}

func (r RabbitConf) GetQueue() string {
	return r.Queue
}
