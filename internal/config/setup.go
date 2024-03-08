package config

import (
	"github.com/caarlos0/env/v10"
)

const (
	ProductionLog  = "production"
	DevelopmentLog = "development"
)

// Application contains data related to application configuration parameters.
type Application struct {
	DryRun          bool   `env:"DRY_RUN" envDefault:"false"`
	ApplicationPort string `env:"APPLICATION_PORT" envDefault:":8080"`
	LogLevel        string `env:"LOG_ENVIRONMENT" envDefault:"development"`
	Repository      RepositoryParameters
}

// RepositoryParameters contains data related to a repository.
type RepositoryParameters struct {
	Host     string `env:"DB_HOST" envDefault:"localhost"`
	Port     string `env:"DB_PORT" envDefault:"3306"`
	User     string `env:"DB_USER" envDefault:"user"`
	Password string `env:"DB_PASSWORD" envDefault:"payments"`
	DBName   string `env:"DBNAME" envDefault:"payments_db"`
}

// Load load application configuration
func Load() (Application, error) {
	cfg := Application{}
	if err := env.Parse(&cfg); err != nil {
		return cfg, err
	}
	repository := RepositoryParameters{}
	if err := env.Parse(&repository); err != nil {
		return cfg, err
	}
	cfg.Repository = repository
	return cfg, nil
}
