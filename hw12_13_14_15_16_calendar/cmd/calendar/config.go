package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

// При желании конфигурацию можно вынести в internal/config.
// Организация конфига в main принуждает нас сужать API компонентов, использовать
// при их конструировании только необходимые параметры, а также уменьшает вероятность циклической зависимости.
type Config struct {
	Logger LoggerConf `yaml:"logger"`
	DB     DBConf     `yaml:"database"`
	Server ServerConf `yaml:"server"`
}

type LoggerConf struct {
	Level   string `yaml:"level"`
	File    string `yaml:"file"`
	MaxSize int    `yaml:"maxSize"`
}

type DBConf struct {
	Type           string `yaml:"type"`
	Name           string `yaml:"name"`
	User           string `yaml:"user"`
	Password       string `yaml:"password"`
	Host           string `yaml:"host"`
	Port           string `yaml:"port"`
	MigrationsPath string `yaml:"migrationsPath"`
}

func (d *DBConf) DSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", d.User, d.Password, d.Host, d.Port, d.Name)
}

func (d *DBConf) updateEnv() {
	d.Name = expandEnv(d.Name)
	d.User = expandEnv(d.User)
	d.Password = expandEnv(d.Password)
	d.Host = expandEnv(d.Host)
	d.Port = expandEnv(d.Port)
}

func expandEnv(value string) string {
	if strings.HasPrefix(value, "${") && strings.HasSuffix(value, "}") {
		return os.Getenv(value[2 : len(value)-1])
	}
	return value
}

type ServerConf struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

func NewConfig(pathConfig string) Config {
	data, err := os.ReadFile(pathConfig)
	if err != nil {
		log.Fatal("no config file")
	}
	var conf Config
	err = yaml.Unmarshal(data, &conf)
	if err != nil {
		log.Fatal("failed unmarshal config file")
	}
	conf.DB.updateEnv()
	return conf
}
