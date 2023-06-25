package main

import (
	"fmt"
	"io"
	"os"
	"strconv"

	"gopkg.in/yaml.v3"
)

/* Environment utility */

func loadEnvStr(key string, result *string) {
	s, ok := os.LookupEnv(key)
	if !ok {
		return
	}

	*result = s
}

func loadEnvUint(key string, result *uint) {
	s, ok := os.LookupEnv(key)
	if !ok {
		return
	}

	n, err := strconv.Atoi(s)

	if err != nil {
		return
	}

	*result = uint(n) // will clamp the negative value
}

/* Configuration */

type pgConfig struct {
	Host string `yaml:"host" json:"host"`
	Port uint   `yaml:"port" json:"port"`

	DBName  string `yaml:"db_name" json:"db_name"`
	SslMode string `yaml:"ssl_mode" json:"ssl_mode"`
}

func (p pgConfig) ConnStr() string {
	return fmt.Sprintf("host=%s port=%d database=%s sslmode=%s", p.Host, p.Port, p.DBName, p.SslMode)
}

func defaultPgConfig() pgConfig {
	return pgConfig{
		Host:    "localhost",
		Port:    5432,
		DBName:  "todo",
		SslMode: "disable",
	}
}

func (p *pgConfig) loadFromEnv() {
	loadEnvStr("KAD_DB_HOST", &p.Host)
	loadEnvUint("KAD_DB_PORT", &p.Port)
	loadEnvStr("KAD_DB_NAME", &p.DBName)
	loadEnvStr("KAD_DB_SSL", &p.SslMode)

}

type listenConfig struct {
	Host string `yaml:"host" json:"host"`
	Port uint   `yaml:"port" json:"port"`
}

func (l listenConfig) Addr() string {
	return fmt.Sprintf("%s:%d", l.Host, l.Port)
}

func defaultListenConfig() listenConfig {
	return listenConfig{
		Host: "127.0.0.1",
		Port: 8080,
	}
}

func (l *listenConfig) loadFromEnv() {
	loadEnvStr("KAD_LISTEN_HOST", &l.Host)
	loadEnvUint("KAD_LISTEN_PORT", &l.Port)
}

type config struct {
	Listen   listenConfig `yaml:"listen" json:"listen"`
	DBConfig pgConfig     `yaml:"db" json:"db"`
}

func (c *config) loadFromEnv() {
	c.Listen.loadFromEnv()
	c.DBConfig.loadFromEnv()
}

func defaultConfig() config {
	return config{
		Listen:   defaultListenConfig(),
		DBConfig: defaultPgConfig(),
	}
}

func loadConfigFromReader(r io.Reader, c *config) error {
	return yaml.NewDecoder(r).Decode(c)
}

func loadConfigFromFile(fn string, c *config) error {
	_, err := os.Stat(fn)

	if err != nil {
		return err
	}

	f, err := os.Open(fn)

	if err != nil {
		return err
	}

	defer f.Close()

	return loadConfigFromReader(f, c)
}

/* How to load the configuration, the highest priority loaded last
 * First: Initialise to default config
 * Second: Replace with environment variables
 * Third: Replace with configuration file
 */

func loadConfig(fn string) config {
	cfg := defaultConfig()
	cfg.loadFromEnv()

	loadConfigFromFile(fn, &cfg)

	return cfg
}
