package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Config struct {
	Server struct {
		Host      string `mapstructure:"host"`
		Port      int    `mapstructure:"port"`
		Prefix    string `mapstructure:"prefix"`
		Migration string `mapstructure:"migration"`
	} `mapstructure:"server"`

	Psql struct {
		Name string `mapstructure:"name"`
		Host string `mapstructure:"host"`
		Port int    `mapstructure:"port"`
		User string `mapstructure:"user"`
		Pass string `mapstructure:"pass"`
	} `mapstructure:"psql"`

	Redis struct {
		Name string `mapstructure:"name"`
		Host string `mapstructure:"host"`
		Port int    `mapstructure:"port"`
		User string `mapstructure:"user"`
		Pass string `mapstructure:"pass"`
	} `mapstructure:"redis"`

	Jwt struct {
		AccessSecret  string `mapstructure:"access_secret"`
		AccessExpire  string `mapstructure:"access_expire"`
		RefreshSecret string `mapstructure:"refresh_secret"`
		RefreshExpire string `mapstructure:"refresh_expire"`
	} `mapstructure:"jwt"`
}

const (
	CONFIG_ARG = "config"
)

func NewConfig(path, name, file string) *Config {
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

	cfg := &Config{}
	if err := viper.Unmarshal(cfg); err != nil {
		panic(fmt.Errorf("failed on reading config file: %w", err))
	}

	return cfg
}
