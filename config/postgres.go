package config

import "time"

type Postgres struct {
	URL            string
	RequestTimeOut time.Duration
}

func PostgresURL() string {
	return cfg.Postgres.URL
}
