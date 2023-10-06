package main

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

func parseEnv() (*Config, error) {
	c := new(Config)
	if err := envconfig.Process(applicationID, c); err != nil {
		return nil, err
	}
	return c, nil
}

type Config struct {
	AccessToken   string        `envconfig:"access_token" required:"true"`
	SleepInterval time.Duration `envconfig:"sleep_interval" default:"1s"`
	UpdateConfig  UpdateConfig
}

type UpdateConfig struct {
	Offset         int      `envconfig:"update_config_offset" default:"0"`
	Limit          int      `envconfig:"update_config_limit" default:"0"`
	Timeout        int      `envconfig:"update_config_timeout" default:"1"`
	AllowedUpdates []string `envconfig:"update_config_allowed_updates"`
}
