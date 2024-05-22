package postgrescontainer

import (
	"context"
	"strings"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

var (
	dbName     = "test"
	dbUser     = "test"
	dbPassword = "test"
)

func StartPostgresContainer(ctx context.Context) (*PostgresContainer, error) {
	container, err := postgres.RunContainer(
		ctx, testcontainers.WithImage("docker.io/postgres:15.2-alpine"),
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPassword),
	)
	if err != nil {
		return nil, err
	}
	endpoint, err := container.Endpoint(ctx, "")
	if err != nil {
		return nil, err
	}
	hostPort := strings.Split(endpoint, ":")
	return &PostgresContainer{
		PostgresContainer: container,
		Host:              hostPort[0],
		Port:              hostPort[1],
		Name:              dbName,
		User:              dbUser,
		Password:          dbPassword,
	}, nil
}

type PostgresContainer struct {
	*postgres.PostgresContainer
	Host     string
	Port     string
	Name     string
	User     string
	Password string
}
