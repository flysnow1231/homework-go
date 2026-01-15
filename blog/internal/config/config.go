package config

import (
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	App      AppConfig      `mapstructure:"app"`
	Log      LogConfig      `mapstructure:"log"`
	MySQL    MySQLConfig    `mapstructure:"mysql"`
	Redis    RedisConfig    `mapstructure:"redis"`
	RabbitMQ RabbitMQConfig `mapstructure:"rabbitmq"`
	Security SecurityConfig `mapstructure:"security"`
}

type AppConfig struct {
	Name               string `mapstructure:"name"`
	Env                string `mapstructure:"env"`
	HTTPAddr           string `mapstructure:"http_addr"`
	ShutdownTimeoutSec int    `mapstructure:"shutdown_timeout_sec"`
}

func (a AppConfig) ShutdownTimeout() time.Duration {
	if a.ShutdownTimeoutSec <= 0 {
		return 15 * time.Second
	}
	return time.Duration(a.ShutdownTimeoutSec) * time.Second
}

type LogConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"` // json/console
	File   string `mapstructure:"file"`   // empty -> stdout
}

type MySQLConfig struct {
	DSN                string `mapstructure:"dsn"`
	MaxOpenConns       int    `mapstructure:"max_open_conns"`
	MaxIdleConns       int    `mapstructure:"max_idle_conns"`
	ConnMaxLifetimeSec int    `mapstructure:"conn_max_lifetime_sec"`
}

type RedisConfig struct {
	Addr     string `mapstructure:"addr"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

type RabbitMQConfig struct {
	URL         string `mapstructure:"url"`
	Exchange    string `mapstructure:"exchange"`
	Queue       string `mapstructure:"queue"`
	RoutingKey  string `mapstructure:"routing_key"`
	ConsumerTag string `mapstructure:"consumer_tag"`
	Prefetch    int    `mapstructure:"prefetch"`
	Concurrency int    `mapstructure:"concurrency"`
}

type SecurityConfig struct {
	JWTSecret string `mapstructure:"jwt_secret"`
}

// Load reads configs/config.yaml and allows env override with prefix BANGKOK_
func Load() (*Config, error) {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath("./configs")
	v.AddConfigPath(".")

	v.SetDefault("app.name", "blog")
	v.SetDefault("app.env", "dev")
	v.SetDefault("app.http_addr", ":8080")
	v.SetDefault("app.shutdown_timeout_sec", 15)

	v.SetDefault("log.level", "info")
	v.SetDefault("log.format", "json")
	v.SetDefault("log.file", "")

	v.SetEnvPrefix("BANGKOK")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	_ = v.ReadInConfig() // optional
	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
