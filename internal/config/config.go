package config

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type ServerConfig struct {
	Host        string `yaml:"host"`
	Port        int    `yaml:"port"`
	SwaggerHost string `yaml:"swagger"`
}

func (c ServerConfig) Addr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

type PostgresConfig struct {
	Host         string `yaml:"host"`
	Port         int    `yaml:"port"`
	User         string `yaml:"user"`
	Password     string `yaml:"password"`
	DatabaseName string `yaml:"databasename"`
}

func (p PostgresConfig) URL() string {
	return fmt.Sprintf(
		`postgres://%s:%s@%s:%d/%s?sslmode=disable`,
		p.User,
		p.Password,
		p.Host,
		p.Port,
		p.DatabaseName,
	)
}

type RedisConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Password string `yaml:"password"`
}

type Config interface {
	ServerConfig() ServerConfig
	PostgresConfig() PostgresConfig
	RedisConfig() RedisConfig
}

type cnf struct {
	Server   ServerConfig   `yaml:"server"`
	Postgres PostgresConfig `yaml:"postgres"`
	Redis    RedisConfig    `yaml:"redis"`
}

func (c *cnf) ServerConfig() ServerConfig {
	return c.Server
}

func (c *cnf) PostgresConfig() PostgresConfig {
	return c.Postgres
}

func (c *cnf) RedisConfig() RedisConfig {
	return c.Redis
}

func ParseFromFile(filePath string) Config {
	file, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("failed open file, %s", err)
	}
	var config cnf
	err = yaml.Unmarshal(file, &config)
	if err != nil {
		log.Fatalf("failed marhsal config, %s", err)
	}
	return &config
}
