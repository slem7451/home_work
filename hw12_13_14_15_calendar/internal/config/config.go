package config

type Config interface {
	GetStorage() string
	GetDB() DBConf
}

type DBConf interface {
	DSN() string
}

type RabbitConf interface {
	URL() string
	GetExchange() string
	GetQueue() string
}
