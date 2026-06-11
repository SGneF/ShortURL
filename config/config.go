package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server ServerConfig `yaml:"server"`
	MySQL  MySQLConfig  `yaml:"mysql"`
	Redis  RedisConfig  `yaml:"redis"`
	Bloom  BloomConfig  `yaml:"bloom"`
}

type ServerConfig struct {
	Port   string `yaml:"port"`
	Domain string `yaml:"domain"`
}

type MySQLConfig struct {
	DSN string `yaml:"dsn"`
}

type RedisConfig struct {
	Addr     string `yaml:"addr"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

type BloomConfig struct {
	Capacity      uint    `yaml:"capacity"`
	FalsePositive float64 `yaml:"false_positive"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	cfg := &Config{}
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

var Default = &Config{
	Server: ServerConfig{Port: "8080", Domain: "http://localhost:8080"},
	MySQL:  MySQLConfig{DSN: "root:password@tcp(127.0.0.1:3306)/shorturl?charset=utf8mb4&parseTime=True&loc=Local"},
	Redis:  RedisConfig{Addr: "127.0.0.1:6379", Password: "", DB: 0},
	Bloom:  BloomConfig{Capacity: 1000000, FalsePositive: 0.001},
}
