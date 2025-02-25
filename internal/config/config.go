package config

import (
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/ilyakaznacheev/cleanenv"
)

var ErrInField = errors.New("no require field")

type Server struct {
	Port uint16 `yaml:"port" validate:"required"`
}

type Config struct {
	Server      Server `yaml:"server" validate:"required"`
	Env         string `yaml:"env" validate:"required"`
}

func LoadConfig(path string) (*Config, error) {
	cfg := new(Config)

	err := cleanenv.ReadConfig(path, cfg)
	if err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	validate := validator.New()

	if err := validate.Struct(cfg); err != nil {
		//nolint:forcetypeassert, errorlint
		for _, err := range err.(validator.ValidationErrors) {
			return nil, fmt.Errorf("%w - %s: %s", ErrInField, err.StructField(), err.Tag())
		}

		return nil, fmt.Errorf("validate config: %w", err)
	}

	return cfg, nil
}
