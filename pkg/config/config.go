package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"log"
	"os"
)

type Application struct {
	Api      Api      `yaml:"api"`
	Auth     Auth     `yaml:"auth"`
	Database Database `yaml:"database"`
}

func LoadApplicationConfig(filepath string) (*Application, error) {
	var config Application
	appConfigData, err := os.ReadFile(filepath)
	if err != nil {
		return &config, err
	}
	err = yaml.Unmarshal(appConfigData, &config)
	if err != nil {
		log.Fatal(err)
		return &config, err
	}
	return &config, nil
}

type Api struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

type Auth struct {
	SecretKey string `yaml:"secret-key"`
	Duration  int    `yaml:"duration"`
}

type Database struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
}

func (d *Database) GetDataSource() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		d.Host,
		d.Port,
		d.Username,
		d.Password,
		d.Database,
	)
}
