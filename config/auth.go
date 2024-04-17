package config

type Auth struct {
	Name string
}

func GetAuthServiceName() string {
	return cfg.Auth.Name
}
