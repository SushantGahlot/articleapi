package main

import (
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/sushantgahlot/articleapi/pkg/config"
)

// https://hackernoon.com/how-to-create-golang-rest-api-project-layout-configuration-part-2-wh2z3y5z
func main() {
	cfg := config.GetConfig()

	direction := "up"
	if direction != "down" && direction != "up" {
		log.Println("-migrate accepts [up, down] values only")
		return
	}

	fmt.Println(cfg.GetDBConnectionString(), "This is connection string")

	m, err := migrate.New("file://db/migrations", cfg.GetDBConnectionString())
	if err != nil {
		log.Printf("%s", err)
		return
	}

	if direction == "up" {
		if err := m.Up(); err != nil {
			log.Printf("failed migrate up: %s", err)
			return
		}
	}

	if direction == "down" {
		if err := m.Down(); err != nil {
			log.Printf("failed migrate down: %s", err)
			return
		}
	}
}
