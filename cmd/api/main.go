package main

import (
	"github.com/sushantgahlot/articleapi/cmd/api/router"
	"github.com/sushantgahlot/articleapi/pkg/application"
	"github.com/sushantgahlot/articleapi/pkg/interrupthandler"
	"github.com/sushantgahlot/articleapi/pkg/logger"
	"github.com/sushantgahlot/articleapi/pkg/server"
)

func main() {
	app, err := application.GetApplication()

	if err != nil {
		logger.Error.Fatal(err.Error())
	}

	r := router.GetHandlers(app)

	srvr, err := server.GetServer(app.Config.GetAPIPort(), logger.Error, r)

	if err != nil {
		logger.Error.Fatal(err.Error())
	}

	logger.Info.Printf("Starting server at %s", app.Config.GetAPIPort())

	go func() {
		if err := srvr.Start(); err != nil {
			logger.Error.Fatal(err.Error())
		}
	}()

	interrupthandler.HandleInterrupt(
		func() {
			if err := srvr.Close(); err != nil {
				logger.Error.Println(err.Error())
			}

			app.DB.DBClient.Close()
		},
	)
}
