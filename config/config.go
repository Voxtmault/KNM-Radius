package config

import "github.com/tkanos/gonfig"

type Configuration struct {
	PHOENIX_DB_USERNAME string
	PHOENIX_DB_PASSWORD string
	PHOENIX_DB_PORT     string
	PHOENIX_DB_HOST     string
	PHOENIX_DB_NAME     string
	RADIUS_DB_USERNAME  string
	RADIUS_DB_PASSWORD  string
	RADIUS_DB_PORT      string
	RADIUS_DB_HOST      string
	RADIUS_DB_NAME      string
}

func GetConfig() Configuration {
	conf := Configuration{}
	gonfig.GetConf("config/config.json", &conf)

	return conf
}
