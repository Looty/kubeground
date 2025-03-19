package config

import (
	"flag"
)

type Config struct {
	Port                int
	InClusterConnection bool
}

func ParseFlags() *Config {
	cfg := &Config{}

	flag.BoolVar(&cfg.InClusterConnection, "incluster", false, "Specifies whether to use in-cluster configuration")
	flag.IntVar(&cfg.Port, "port", 8080, "Specifies the port number")
	flag.Parse()

	return cfg
}
