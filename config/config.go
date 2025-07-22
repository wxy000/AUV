package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

type AppConfig struct {
	DB struct {
		Driver string `yaml:"driver"`
		Sqlite struct {
			Path string `yaml:"path"`
		} `yaml:"sqlite"`
		Mysql struct {
			Host     string `yaml:"host"`
			Port     int    `yaml:"port"`
			User     string `yaml:"user"`
			Password string `yaml:"password"`
			Name     string `yaml:"name"`
			SSLMode  string `yaml:"ssl_mode"`
		} `yaml:"mysql"`
	} `yaml:"database"`
	Server struct {
		Port         string `yaml:"port"`
		LimitRate    int    `yaml:"limit_rate"`
		StaticPrefix string `yaml:"static_prefix"`
	} `yaml:"server"`
	JWT struct {
		Secret             string `yaml:"secret"`
		ExpiresHours       int    `yaml:"expires_hours"`
		RefreshWindowHours int    `yaml:"refresh_window_hours"`
	} `yaml:"jwt"`
	Admin struct {
		Username string `yaml:"username"`
		Password string `yaml:"password"`
	} `yaml:"admin"`
	User struct {
		DefaultPassword string `yaml:"default_password"`
	} `yaml:"user"`
	HitokotoFile string `yaml:"hitokoto_file"`
}

var Cfg *AppConfig

func LoadConfig(path string) (*AppConfig, error) {
	config := &AppConfig{}

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func(file *os.File) {
		file.Close()
	}(file)

	if err := yaml.NewDecoder(file).Decode(config); err != nil {
		return nil, err
	}
	Cfg = config
	return config, nil
}
