package store

type Config struct {
	DBURL string `toml:"db_url"`
}

func NewConfig() *Config {
	return &Config{}
}
