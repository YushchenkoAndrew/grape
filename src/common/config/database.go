package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type DatabaseConfig struct {
	Organization struct {
		Name string `mapstructure:"name"`
	} `mapstructure:"organization"`

	User struct {
		Name     string `mapstructure:"name"`
		Password string `mapstructure:"password"`
	} `mapstructure:"users"`
}

func NewDatabaseConfig(path, name, file string) *DatabaseConfig {
	if value, ok := os.LookupEnv(CONFIG_ARG); ok {
		viper.AddConfigPath(filepath.Join(value, path))
	} else {
		viper.AddConfigPath(filepath.Join(".", path))
	}

	viper.SetConfigName(name)
	viper.SetConfigType(file)

	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("failed on reading config file: %w", err))
	}

	cfg := &DatabaseConfig{}
	if err := viper.Unmarshal(cfg); err != nil {
		panic(fmt.Errorf("failed on reading config file: %w", err))
	}

	return cfg
}
