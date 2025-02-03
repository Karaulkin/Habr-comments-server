package main

import (
	"errors"
	"flag"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	var host, port, username, password, database, migrationsPath, mode string

	flag.StringVar(&host, "host", "", "host to storage")
	flag.StringVar(&port, "port", "", "port to host")
	flag.StringVar(&username, "login", "", "login to storage")
	flag.StringVar(&password, "password", "", "password to storage")
	flag.StringVar(&database, "db", "", "name of PostgreSQL database")
	flag.StringVar(&migrationsPath, "path", "./migrations", "path to migrations files")
	flag.StringVar(&mode, "mode", "up", "migration mode: up or down") // "up" по умолчанию
	flag.Parse()

	if host == "" || port == "" || username == "" || password == "" || database == "" {
		log.Fatal("All database connection parameters are required")
	}
	if migrationsPath == "" {
		log.Fatal("Migrations path is required")
	}

	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		username, password, host, port, database)

	m, err := migrate.New(
		"file://"+migrationsPath,
		dsn,
	)
	if err != nil {
		log.Fatal("Failed to initialize migrations: ", err)
	}

	if mode == "down" {
		if err := m.Down(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			log.Fatal("Failed to apply down migrations: ", err)
		}
		fmt.Println("Migrations rolled back successfully")
		return
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatal("Migration error: ", err)
	}

	fmt.Println("Migrations applied successfully")
}
