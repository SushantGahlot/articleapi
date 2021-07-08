package config

import (
	"fmt"
	"os"
)

type Config struct {
	dbUser  string
	dbPass  string
	dbHost  string
	dbPort  string
	dbName  string
	apiPort string
}

func GetConfig() *Config {
	conf := Config{
		dbUser:  os.Getenv("POSTGRES_USER"),
		dbPass:  os.Getenv("POSTGRES_PASSWORD"),
		dbHost:  os.Getenv("POSTGRES_HOST"),
		dbPort:  os.Getenv("POSTGRES_PORT"),
		dbName:  os.Getenv("POSTGRES_DB"),
		apiPort: os.Getenv("API_PORT"),
	}
	return &conf
}

func (c *Config) GetDBConnectionString() string {
	var connStr string = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", c.dbUser, c.dbPass, c.dbHost, c.dbPort, c.dbName)
	return connStr
}

func (c *Config) GetDBPort() string {
	return c.dbPort
}

func (c *Config) GetAPIPort() string {
	return c.apiPort
}
