package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/jinzhu/copier"
	"github.com/spf13/viper"
)

type Config struct {
	Server struct {
		Address    string `mapstructure:"address"`
		Prefix     string `mapstructure:"prefix"`
		Migrations string `mapstructure:"migrations"`
	} `mapstructure:"server"`

	Psql struct {
		Name     string `mapstructure:"name"`
		Host     string `mapstructure:"host"`
		Port     int    `mapstructure:"port"`
		Username string `mapstructure:"username"`
		Password string `mapstructure:"password"`
		Logger   bool   `mapstructure:"logger"`
	} `mapstructure:"psql" copier:"-"`

	Redis struct {
		DB       int    `mapstructure:"db"`
		Address  string `mapstructure:"address"`
		Username string `mapstructure:"username"`
		Password string `mapstructure:"password"`
	} `mapstructure:"redis" copier:"-"`

	Jwt struct {
		AccessSecret  string `mapstructure:"access_secret"`
		AccessExpire  string `mapstructure:"access_expire"`
		RefreshSecret string `mapstructure:"refresh_secret"`
		RefreshExpire string `mapstructure:"refresh_expire"`
	} `mapstructure:"jwt" copier:"-"`

	Void struct {
		Url      string `mapstructure:"url"`
		Username string `mapstructure:"username"`
		Password string `mapstructure:"password"`
	} `mapstructure:"void" copier:"-"`
}

const (
	CONFIG_ARG = "config"
)

var config *Config

func GetGlobalConfig() *Config {
	return config
}

func GetConfigFile() string {
	return viper.ConfigFileUsed()
}

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

	config = &Config{}
	copier.Copy(config, cfg)
	return cfg
}
