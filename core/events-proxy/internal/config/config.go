package config

import (
	"github.com/cristalhq/aconfig"
	"github.com/cristalhq/aconfig/aconfigyaml"
)

type ProjectConfig struct {
	Version string `env:"VERSION" yaml:"version" default:"v1"`
}

type NatsConfig struct {
	Urls string `env:"URLS" yaml:"urls" default:"nats://nats:4333"`
}

type PostgresConfig struct {
	Host               string `env:"HOST" yaml:"host" default:"postgres"`
	Port               string `env:"PORT" yaml:"port" default:"5432"`
	DBName             string `env:"DBNAME" yaml:"dbName" default:"events_proxy"`
	Username           string `env:"USERNAME" yaml:"username" default:"postgres"`
	Password           string `env:"PASSWORD" yaml:"password" default:"postgres"`
	SslMode            string `env:"SSL_MODE" yaml:"sslMode" default:"disable"`
	MaxOpenConns       int    `env:"MAX_OPEN_CONNS" yaml:"maxOpenConns" default:"0"`
	ConnMaxLifetimeSec int    `env:"CONN_MAX_LIFETIME_SEC" yaml:"connMaxLifetimeSec" default:"0"`
	MaxIdleConns       int    `env:"MAX_IDLE_CONNS" yaml:"maxIdleConns" default:"2"`
	ConnMaxIdleTimeSec int    `env:"CONN_MAX_IDLE_TIME_SEC" yaml:"connMaxIdleTimeSec" default:"0"`
	MigrationsFolder   string `env:"MIGRATIONS" yaml:"migrations" default:"migrations"`
}

type MqttServerConfig struct {
	Id   string `env:"ID" yaml:"id" default:"t1"`
	Host string `env:"HOST" yaml:"host" default:"0.0.0.0"`
	Port string `env:"PORT" yaml:"port" default:"1883"`
}

type JwtConfig struct {
	Secret      []byte `env:"SECRET" yaml:"secret" required:"true"`
	DurationMin int    `env:"DURATION_MIN" yaml:"durationMin" default:"120"`
}

type MetricsServerConfig struct {
	Host string `env:"HOST" yaml:"host" default:"0.0.0.0"`
	Port int    `env:"PORT" yaml:"port" default:"9090"`
	Path string `env:"PATH" yaml:"path" default:"/metrics"`
}

type StatusServerConfig struct {
	Host           string `env:"HOST" yaml:"host" default:"0.0.0.0"`
	Port           int    `env:"PORT" yaml:"port" default:"8000"`
	LivelinessPath string `env:"LIVELINESS_PATH" yaml:"path" default:"/live"`
	ReadinessPath  string `env:"READINESS_PATH" yaml:"path" default:"/ready"`
}

type Config struct {
	Project    ProjectConfig       `env:"PROJECT" yaml:"project"`
	Nats       NatsConfig          `env:"NATS" yaml:"nats"`
	Database   PostgresConfig      `env:"DATABASE" yaml:"database"`
	MqttServer MqttServerConfig    `env:"MQTT_SERVER" yaml:"mqttServer"`
	Jwt        JwtConfig           `env:"JWT" yaml:"jwt"`
	Metrics    MetricsServerConfig `env:"METRICS" yaml:"metrics"`
	Status     StatusServerConfig  `env:"STATUS" yaml:"statust"`
}

func ReadConfig() (cfg Config, err error) {
	loader := aconfig.LoaderFor(&cfg, aconfig.Config{
		SkipEnv:   false,
		SkipFlags: true,
		Files:     []string{"config.yaml"},
		FileDecoders: map[string]aconfig.FileDecoder{
			".yaml": aconfigyaml.New(),
		},
	})

	err = loader.Load()
	return
}
