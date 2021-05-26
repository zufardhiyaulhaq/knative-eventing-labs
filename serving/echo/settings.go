package main

import (
	"log"

	"github.com/kelseyhightower/envconfig"
)

type Settings struct {
	Port     string `envconfig:"PORT"`
	Revision string `envconfig:"K_REVISION"`
	Service  string `envconfig:"K_SERVICE"`
}

func NewSettings() Settings {
	var settings Settings

	err := envconfig.Process("", &settings)
	if err != nil {
		log.Fatalln(err)
	}

	return settings
}
