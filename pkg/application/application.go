package application

import (
	"github.com/sushantgahlot/articleapi/pkg/config"
	"github.com/sushantgahlot/articleapi/pkg/database"
)

type Application struct {
	DB     *database.DB
	Config *config.Config
}

func GetApplication() (*Application, error) {
	conf := config.GetConfig()
	db, err := database.GetDBConnection(conf.GetDBConnectionString())

	if err != nil {
		return nil, err
	}

	return &Application{
		DB:     db,
		Config: conf,
	}, nil
}
