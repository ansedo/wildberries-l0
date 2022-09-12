package config

import (
	"context"
	"flag"

	"github.com/caarlos0/env/v6"
	"go.uber.org/zap"

	"github.com/ansedo/wildberries-l0/internal/logger"
)

type Config struct {
	RunAddress  string `env:"RUN_ADDRESS" envDefault:":8080"`
	DatabaseURI string `env:"DATABASE_URI" envDefault:"postgres://postgres:password@localhost:5432/postgres"`
	Stan        Stan
}

type Stan struct {
	ClusterID string `env:"STAN_CLUSTER_ID" envDefault:"test-cluster"`
	Subject   string `env:"STAN_SUBJECT" envDefault:"test-subject"`
}

func New(ctx context.Context) *Config {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		logger.FromCtx(ctx).Fatal("load Config: parse flags", zap.Error(err))
	}
	flag.StringVar(&cfg.RunAddress, "a", cfg.RunAddress, `server address to listen on`)
	flag.StringVar(&cfg.DatabaseURI, "d", cfg.DatabaseURI, `uri to connect to database`)
	flag.StringVar(&cfg.Stan.ClusterID, "c", cfg.Stan.ClusterID, `nats streaming cluster name`)
	flag.StringVar(&cfg.Stan.Subject, "s", cfg.Stan.Subject, `nats streaming subject (channel) name`)
	flag.Parse()
	return &cfg
}
