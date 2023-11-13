package reback

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
)

var db *sql.DB

func TestMain(m *testing.M) {
	// uses a sensible default on windows (tcp/http) and linux/osx (socket)
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not construct pool: %s", err)
	}

	// uses pool to try to connect to Docker
	err = pool.Client.Ping()
	if err != nil {
		log.Fatalf("Could not connect to Docker: %s", err)
	}

	// pulls an image, creates a container based on it and runs it
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "15.4",
		Env: []string{
			"POSTGRES_USER=username",
			"POSTGRES_PASSWORD=password",
			"POSTGRES_DB=dbname",
			"listen_addresses = '*'",
		},
	}, func(config *docker.HostConfig) {
		// set AutoRemove to true so that stopped container goes away by itself
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{
			Name: "no",
		}
	})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	time.Sleep(30 * time.Second)
	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	if err := pool.Retry(func() error {
		hostAndPort := resource.GetHostPort("5432/tcp")
		var err error
		pool, err := pgxpool.New(context.Background(), fmt.Sprintf("postgres://username:password@%s/dbname?sslmode=disable", hostAndPort))
		if err != nil {
			return err
		}

		db = stdlib.OpenDBFromPool(pool)
		if err != nil {
			return err
		}
		return db.Ping()
	}); err != nil {
		log.Fatalf("Could not connect to database: %s", err)
	}

	b, err := os.ReadFile("./migration/up.sql")
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(string(b))
	if err != nil {
		log.Fatal(err)
	}
	// set db conn
	SetDBConn(context.TODO(), db)

	code := m.Run()

	// You can't defer this because os.Exit doesn't care for defer
	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}
