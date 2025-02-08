package config

import (
	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type Config struct {
	Auth struct {
		Encrypt struct {
			SecretKey string `mapstructure:"secret_key" validate:"required"`
		} `mapstructure:"encrypt"`
	} `mapstructure:"auth"`

	Database struct {
		Postgres struct {
			ConnectionString string `mapstructure:"connection_string" validate:"required"`
		} `mapstructure:"postgres"`
	} `mapstructure:"database"`
}

func NewConfig() (*Config, error) {
	c := &Config{}

	err := viper.UnmarshalExact(c)
	if err != nil {
		return nil, err
	}

	validate := validator.New()

	err = validate.Struct(c)
	if err != nil {
		return nil, err
	}

	return c, nil
}
