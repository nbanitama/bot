package chatbot

import (
	"log"

	logging "gopkg.in/tokopedia/logging.v1"
)

type ServerConfig struct {
	Name string
}

type PostgresqlConfig struct {
	Host     string
	Username string
	Password string
	DBName   string
}

type Config struct {
	Server     ServerConfig
	Postgresql PostgresqlConfig
}

var (
	Main Config
)

func NewMainConfig() {
	ok := logging.ReadModuleConfig(&Main, "/etc/chatbot", "app") || logging.ReadModuleConfig(&Main, "files/etc/chatbot", "app")
	if !ok {
		log.Fatal("failed to read config")
	}

	log.Println("Config loaded")
}
