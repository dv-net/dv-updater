package config

import (
	"time"

	"github.com/tkcrm/mx/logger"
)

type (
	Config struct {
		App        AppConfig        `yaml:"app"`
		HTTP       HTTPConfig       `yaml:"http"`
		Log        logger.Config    `yaml:"log"`
		AutoUpdate AutoUpdateConfig `yaml:"auto_update"`
	}

	AppConfig struct {
		Profile string `yaml:"profile" default:"dev"`
	}

	HTTPCorsConfig struct {
		Enabled        bool     `yaml:"enabled" default:"true" usage:"allows to disable cors" example:"true / false"`
		AllowedOrigins []string `yaml:"allowed_origins"`
	}

	HTTPConfig struct {
		Host               string         `yaml:"host" default:"localhost"`
		Port               string         `yaml:"port" default:"8081"`
		FetchInterval      time.Duration  `yaml:"fetch_interval" env:"FETCH_INTERVAL" default:"30s"`
		ConnectTimeout     time.Duration  `yaml:"connect_timeout" env:"CONNECT_TIMEOUT" default:"5s"`
		ReadTimeout        time.Duration  `yaml:"read_timeout" env:"READ_TIMEOUT" default:"10s"`
		WriteTimeout       time.Duration  `yaml:"write_timeout" env:"WRITE_TIMEOUT" default:"10s"`
		MaxHeaderMegabytes int            `yaml:"max_header_megabytes" env:"MAX_HEADER_MEGABYTES" default:"1"`
		Cors               HTTPCorsConfig `yaml:"cors"`
	}

	SeedConfig struct {
		Base string `yaml:"base" default:"seeds"`
	}

	AutoUpdateConfig struct {
		Enabled bool `yaml:"enabled" default:"true"`
	}
)
