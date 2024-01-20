package config

type NatsConfig struct {
	Urls string
}

type PostgresConfig struct {
	Host               string
	Port               string
	DBName             string
	Username           string
	Password           string
	SslMode            bool
	MaxOpenConns       int
	ConnMaxLifetimeSec int
	MaxIdleConns       int
	ConnMaxIdleTimeSec int
}

type Config struct {
	Nats     NatsConfig
	Postgres PostgresConfig
}

func ReadConfig() Config {
	return Config{}
}
