package main

import (
	"github.com/kelseyhightower/envconfig"
)

type ConfigSpec struct {
	RedisAddr string `envconfig:"redis_addr" default:"localhost:6379"`
	RedisDB   int `envconfig:"redis_db" default:"0"`
}

var Config ConfigSpec

func LoadConfig() {
	err := envconfig.Process("hooker", &Config)
	if err != nil {
		panic(err)
	}
}
