package config

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
)

type DBCreds struct {
	Host         string `mapstructure:"DB_HOST"`
	Password     string `mapstructure:"DB_PASSWORD"`
	Username     string `mapstructure:"DB_USERNAME"`
	Port         int    `mapstructure:"DB_PORT"`
	DatabaseName string `mapstructure:"DB_DATABASE_NAME"`
}

type GRPCServer struct {
	Port int `mapstructure:"GRPC_SERVER_PORT"`
}

type Config struct {
	viper    *viper.Viper
	Database DBCreds    `mapstructure:",squash"`
	Server   GRPCServer `mapstructure:",squash"`
}

func NewConfig() *Config {
	return &Config{
		viper: viper.New(),
	}
}

func (c *Config) Load(configFile string) error {
	viper.SetConfigFile(configFile)
	err := viper.ReadInConfig()
	if err != nil {
		return fmt.Errorf("read in config: %w", err)
	}
	if err := viper.Unmarshal(&c); err != nil {
		log.Fatal(err)
	}
	return nil
}
