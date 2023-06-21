# Database Management for Postgres Containers

This repository contains a Go package for managing Postgres test containers and running migrations. The code allows you to automatically setup a Postgres container, connect to it and run SQL migrations. The operations are wrapped inside easy-to-use functions for convenience.

## Features

- **StartPostgresTestContainer:** This function takes a struct of `ContainerParams` type and spins up a new Postgres container using the parameters from the struct. The Postgres container uses the `postgis/postgis` Docker image. It returns a connection to the database.
- **RunMigrations:** This function accepts a connected database and a root path to your migrations. It executes all migrations found in the given path on the connected database. The migrations should be SQL files with 'up.sql' in their name.
- **SortFilesByName:** A helper function that sorts files in ascending order based on their names. Used to ensure migrations are run in the correct order.

## Usage

### StartPostgresTestContainer

```go
conn := StartPostgresTestContainer(ContainerParams{
    Port:     "5432",
    Password: "myPassword",
    Username: "myUsername",
})
```

### RunMigrations

```go
err := RunMigrations(conn, "/path/to/migrations")
if err != nil {
    log.Fatal("Error running migrations: ", err)
}
```

## Requirements

This package uses several dependencies:

- `github.com/docker/go-connections/nat`
- `github.com/jmoiron/sqlx`
- `github.com/testcontainers/testcontainers-go`
- `github.com/testcontainers/testcontainers-go/wait`

Please make sure these packages are installed in your Go environment.

## Notes

Please ensure that your migrations are SQL files with 'up.sql' in their names. The `RunMigrations` function will only execute those files.

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.
