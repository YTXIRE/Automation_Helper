package config

import (
	"backend/pkg/logging"
	"github.com/ilyakaznacheev/cleanenv"
	"sync"
)

type Config struct {
	IsDebug *bool `yaml:"is_debug" env-default:"true"`
	Listen  struct {
		Type   string `yaml:"type" env-default:"port"`
		BindIP string `yaml:"bind_ip" env-default:"127.0.0.1"`
		Port   string `yaml:"port" env-default:"8080"`
	} `yaml:"listen"`
	MongoDB struct {
		Host     string `yaml:"host" env-default:"localhost"`
		Port     string `yaml:"port" env-default:"27017"`
		Database string `yaml:"database"`
		AuthDB   string `yaml:"auth_db"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
	} `yaml:"mongodb"`
	Jwt struct {
		AccessToken struct {
			TTL  int64  `yaml:"ttl" env-default:"12"`
			Type string `yaml:"type" env-default:"h"`
		} `yaml:"access_token"`
		RefreshToken struct {
			TTL  int64  `yaml:"ttl" env-default:"30"`
			Type string `yaml:"type" env-default:"d"`
		} `yaml:"refresh_token"`
		SecretKey string `yaml:"secret_key" env-default:"verySecretKey"`
	} `yaml:"jwt"`
}

var instance *Config
var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		logger := logging.GetLogger()
		logger.Info("Read application configuration")
		instance = &Config{}
		if err := cleanenv.ReadConfig("config.yaml", instance); err != nil {
			help, _ := cleanenv.GetDescription(instance, nil)
			logger.Info(help)
			logger.Fatal(err)
		}
	})
	return instance
}
