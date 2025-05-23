package config

import (
	"errors"
	"flag"
	"fmt"
	"github.com/caarlos0/env/v6"
	"os"
)

var ErrEnvConfigCreation = errors.New("config creation error via environments")

const (
	RunAddressFlag            = "a"
	DatabaseURIFlag           = "d"
	AccrualSystemAddressFlag  = "r"
	AuthKeyNameFlag           = "k"
	RunAddressDefaultValue    = "localhost:8081"
	DatabaseURIDefaultValue   = "postgres://localhost:5432/postgres?sslmode=disable"
	AccrualSystemDefaultValue = "http://localhost:8080"
	GophermartFlagName        = "gophermart"
)

type Config struct {
	RunAddress           string `env:"RUN_ADDRESS"`
	DatabaseURI          string `env:"DATABASE_URI"`
	AccrualSystemAddress string `env:"ACCRUAL_SYSTEM_ADDRESS"`
	AuthKey              string `env:"AUTH_KEY"`
}

func New() (*Config, error) {
	var cfg Config

	flagSet := flag.NewFlagSet(GophermartFlagName, flag.ContinueOnError)
	flagSet.StringVar(&cfg.RunAddress, RunAddressFlag, RunAddressDefaultValue, "help message for flagname")
	flagSet.StringVar(&cfg.DatabaseURI, DatabaseURIFlag, DatabaseURIDefaultValue, "help message for flagname")
	flagSet.StringVar(&cfg.AccrualSystemAddress, AccrualSystemAddressFlag, AccrualSystemDefaultValue, "help message for flagname")
	_ = flagSet.Parse(os.Args[1:])

	err := env.Parse(&cfg)
	if err != nil {
		return nil, errors.Join(ErrEnvConfigCreation, err)
	}

	return &cfg, nil
}

func (c *Config) String() string {
	return fmt.Sprintf(
		"Config{\n"+
			"  RunAddress: %s,\n"+
			"  DatabaseURI: %s,\n"+
			"  AccrualSystemAddress: %s,\n"+
			"  AuthKey: ****\n"+
			"}",
		c.RunAddress,
		c.DatabaseURI,
		c.AccrualSystemAddress,
	)
}
