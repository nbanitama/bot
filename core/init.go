package core

import (
	config "github.com/tokopedia/chatbot-scheduler/config/chatbot"
	"github.com/tokopedia/chatbot-scheduler/postgresql"
)

type TaskModule struct {
	cfg *config.Config
}

var (
	postgresConnection *postgresql.Connection
	err                error
)

func NewTaskModule(cfg *config.Config) (*TaskModule, error) {
	task := TaskModule{
		cfg: cfg,
	}

	postgresConnection, err = postgresql.NewConnection(cfg.Postgresql.Host, cfg.Postgresql.DBName, cfg.Postgresql.Username, cfg.Postgresql.Password)
	if err != nil {
		return nil, err
	}

	return &task, nil
}
