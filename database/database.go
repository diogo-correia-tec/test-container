package database

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/docker/go-connections/nat"
	"github.com/jmoiron/sqlx"
	tc "github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type ContainerParams struct {
	Port     string
	Password string
	Username string
}

func StartPostgresTestContainer(params ContainerParams) (conn *sqlx.DB) {
	log.Println("Starting Postgres container...")
	postgresPort := nat.Port(params.Port + "/tcp")
	postgres, err := tc.GenericContainer(context.Background(),
		tc.GenericContainerRequest{
			ContainerRequest: tc.ContainerRequest{
				Image:        "postgis/postgis",
				ExposedPorts: []string{postgresPort.Port()},
				Env: map[string]string{
					"POSTGRES_PASSWORD": params.Password,
					"POSTGRES_USER":     params.Username,
				},
				WaitingFor: wait.ForAll(
					wait.ForLog("database system is ready to accept connections"),
					wait.ForListeningPort(postgresPort),
				),
			},
			Started: true, // auto-start the container
		})
	if err != nil {
		log.Fatal("start:", err)
	}

	hostPort, err := postgres.MappedPort(context.Background(), postgresPort)
	if err != nil {
		log.Fatal("map:", err)
	}
	postgresURLTemplate := "postgres://postgres:postgres@localhost:%s?sslmode=disable"
	postgresURL := fmt.Sprintf(postgresURLTemplate, hostPort.Port())
	log.Printf("Postgres container started, running at:  %s\n", postgresURL)

	conn, err = sqlx.Connect("postgres", postgresURL)
	if err != nil {
		log.Fatal("connect:", err)
	}

	return conn
}

func RunMigrations(conn *sqlx.DB, migrationsRootPath string) error {
	f, err := os.Open(migrationsRootPath)
	if err != nil {
		log.Fatal(err)
	}
	files, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if !strings.Contains(file.Name(), "up.sql") {
			continue
		}

		log.Printf("Running migrations. File: %s", file.Name())
		migrationByte, errRead := ioutil.ReadFile(migrationsRootPath + "/" + file.Name())
		if errRead != nil {
			log.Fatal("Error reading file: ", errRead)
		}

		migration := string(migrationByte)
		_, err := conn.Exec(migration)
		if err != nil {
			log.Fatal("Error running migration: ", err)
		}
	}

	return nil
}
