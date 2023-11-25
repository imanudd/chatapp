package config

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"log"
	"path"
	"runtime"
)

type Config struct {
	Server   Server
	Database Database
}

type Server struct {
	ENV   string `envconfig:"APP_ENV" required:"true"`
	NAME  string `envconfig:"APP_NAME" required:"true"`
	PORT  string `envconfig:"APP_PORT" required:"true"`
	DEBUG bool   `envconfig:"APP_DEBUG" default:"false"`
}

type Database struct {
	Mongo *Mongo
}

type Mongo struct {
	URI string `envconfig:"MONGO_URI" required:"true"`
	Db  string `envconfig:"MONGO_DB" required:"true"`
}

func LoadDefault() *Config {
	return load("default", ".env")
}

// load config and populate to config struct
func load(file string, env string) *Config {
	var config Config

	readEnv(&config, env)
	err := envconfig.Process("", &config)
	if err != nil {
		panic(err)
	}
	return &config
}

func readEnv(cfg *Config, env string) {
	err := godotenv.Overload(getSourcePath() + "/../" + env)
	if err != nil {
		log.Print(err)
	}
}

func getSourcePath() string {
	_, filename, _, _ := runtime.Caller(1)
	return path.Dir(filename)
}
