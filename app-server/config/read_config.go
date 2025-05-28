package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type ServerConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type DatabaseConfig struct {
	Type     string `yaml:"type"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Database string `yaml:"database"`
	Charset  string `yaml:"charset"`
}

type AuthenticationConfig struct {
	JwtSecret string `yaml:"jwtsecret"`
}

type Config struct {
	Server         ServerConfig         `yaml:"server"`
	Database       DatabaseConfig       `yaml:"database"`
	Authentication AuthenticationConfig `yaml:"authentication"`
}

// 默认配置
func defaultConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Host: "localhost",
			Port: 8081,
		},
		Database: DatabaseConfig{
			Type:     "mysql",
			Username: "root",
			Password: "123456",
			Host:     "localhost",
			Port:     3306,
			Database: "app_server",
			Charset:  "utf8mb4",
		},
		Authentication: AuthenticationConfig{
			JwtSecret: "defaultsecret",
		},
	}
}

var cfg *Config

// GetConfig 返回全局配置
func GetConfig() *Config {
	if cfg == nil {
		cfg = defaultConfig()
	}
	return cfg
}

// 保存配置到文件
func SaveConfig(path string, cfg *Config) error {
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

// 加载配置，如果不存在则生成默认配置文件
func LoadConfig(path string) error {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		cfg := defaultConfig()
		if err := SaveConfig(path, cfg); err != nil {
			return err
		}
		return nil
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return err
	}
	return nil
}
