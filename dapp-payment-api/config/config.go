package config

import (
	"encoding/json"
	"io/ioutil"
)

type ServerConfig struct {
	LogFilePath         string
	ListenPort          int
	RefreshCacheSeconds int
	Database            DatabaseConfig
	CORS                CORSConfig
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
}

type CORSConfig struct {
	AllowedOrigins   []string
	AllowedMethods   []string
	AllowedHeaders   []string
	AllowCredentials bool
	MaxAge           int
}

func NewServerConfig() *ServerConfig {
	ret := new(ServerConfig)
	return ret
}

func (c *ServerConfig) ReadFromFile(filename string) error {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	err = json.Unmarshal(b, c)
	if err != nil {
		return err
	}

	return nil
}
