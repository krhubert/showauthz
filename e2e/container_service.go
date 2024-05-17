package e2e

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

const containerE2ELabel = "authz-e2e"

const containerExpireSeconds = 240

type ContainerService struct {
	*dockertest.Pool
	dockerHost string

	PostgresCtnr     *dockertest.Resource
	PostgresHostPort string
	PostgresPort     string

	SpicedbMigraitonCtnr *dockertest.Resource
	SpicedbCtnr          *dockertest.Resource
	SpicedbHostPort      string
	SpicedbPort          string
}

func NewContainerService() (*ContainerService, error) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		return nil, err
	}

	dockerHost := "host.docker.internal"
	if os.Getenv("CI") == "true" {
		// github action always assign 172.17.0.1 as docker host
		dockerHost = "172.17.0.1"
	}
	return &ContainerService{Pool: pool, dockerHost: dockerHost}, nil
}

func (cs *ContainerService) HandleInterrupt() {
	terminate := make(chan os.Signal, 1)
	signal.Notify(terminate, os.Interrupt)
	go func() {
		<-terminate
		cs.Purge()
		os.Exit(1)
	}()
}

func (cs ContainerService) Purge() {
	if cs.SpicedbMigraitonCtnr != nil {
		_ = cs.Pool.Purge(cs.SpicedbMigraitonCtnr)
	}

	if cs.SpicedbCtnr != nil {
		_ = cs.Pool.Purge(cs.SpicedbCtnr)
	}

	if cs.PostgresCtnr != nil {
		_ = cs.Pool.Purge(cs.PostgresCtnr)
	}
}

// Cleanup finds all containers with the label "authz-e2e" and removes them.
func (cs *ContainerService) Cleanup() error {
	images, err := cs.Client.ListContainers(docker.ListContainersOptions{
		Filters: map[string][]string{
			"label": {containerE2ELabel},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to list containers: %w", err)
	}

	for _, image := range images {
		if err := cs.Client.RemoveContainer(docker.RemoveContainerOptions{
			ID:            image.ID,
			Force:         true,
			RemoveVolumes: true,
		}); err != nil {
			return fmt.Errorf("failed to remove container: %w", err)
		}
	}
	return nil
}

func (cs *ContainerService) RunPostgres() error {
	ctnr, err := cs.RunWithOptions(
		&dockertest.RunOptions{
			Name:         "authz-e2e-postgres",
			Repository:   "postgres",
			Cmd:          []string{"-c", "max_connections=300"},
			Tag:          "15.3-alpine3.18",
			Env:          []string{"POSTGRES_HOST_AUTH_METHOD=trust"},
			ExposedPorts: []string{"5432/tcp"},
			Labels:       map[string]string{containerE2ELabel: "true"},
		},
		func(config *docker.HostConfig) {
			config.AutoRemove = true
			config.RestartPolicy = docker.RestartPolicy{Name: "no"}
		},
	)
	if err != nil {
		return fmt.Errorf("failed to start postgres container: %w", err)
	}

	if err := ctnr.Expire(containerExpireSeconds); err != nil {
		return fmt.Errorf("failed to set expiration on postgres container: %w", err)
	}

	cs.PostgresCtnr = ctnr
	cs.PostgresPort = ctnr.GetPort("5432/tcp")
	cs.PostgresHostPort = ctnr.GetHostPort("5432/tcp")

	if err := cs.WaitForPostgres(); err != nil {
		return fmt.Errorf("failed to wait for postgres: %w", err)
	}

	return nil
}

func (cs *ContainerService) WaitForPostgres() error {
	return backoff.Retry(func() error {
		db, err := sql.Open("pgx", fmt.Sprintf("postgres://postgres:@%s/postgres", cs.PostgresHostPort))
		if err != nil {
			return err
		}
		defer db.Close()
		_, err = db.ExecContext(context.Background(), "SELECT 1;")
		return err
	},
		backoff.WithMaxRetries(backoff.NewConstantBackOff(1*time.Second), 20),
	)
}

func (cs *ContainerService) RunSpicedbMigration() error {
	ctnr, err := cs.RunWithOptions(
		&dockertest.RunOptions{
			Name:       "authz-e2e-spicedb-migration",
			Repository: "authzed/spicedb",
			Tag:        "v1.31.0",
			Env: []string{
				"SPICEDB_DATASTORE_ENGINE=postgres",
				fmt.Sprintf("SPICEDB_DATASTORE_CONN_URI=postgres://postgres:@%s:%s/postgres", cs.dockerHost, cs.PostgresPort),
			},
			Labels: map[string]string{containerE2ELabel: "true"},
			Cmd:    []string{"migrate", "head"},
		},
		func(config *docker.HostConfig) {
			config.AutoRemove = true
			config.RestartPolicy = docker.RestartPolicy{Name: "no"}
		},
	)
	if err != nil {
		return fmt.Errorf("failed to start spicedb migration container: %w", err)
	}

	if err := ctnr.Expire(containerExpireSeconds); err != nil {
		return fmt.Errorf("failed to set expiration on spicedb migration container: %w", err)
	}

	cs.SpicedbMigraitonCtnr = ctnr
	if err := cs.WaitForSpicedbMigration(); err != nil {
		return fmt.Errorf("failed to wait for spicedb migration: %w", err)
	}
	return nil
}

func (cs *ContainerService) WaitForSpicedbMigration() error {
	return backoff.Retry(func() error {
		db, err := sql.Open("pgx", fmt.Sprintf("postgres://postgres:@%s/postgres", cs.PostgresHostPort))
		if err != nil {
			return err
		}
		defer db.Close()
		_, err = db.ExecContext(context.Background(), "SELECT 1 from alembic_version")
		return err
	},
		backoff.WithMaxRetries(backoff.NewConstantBackOff(1*time.Second), 40),
	)
}

func (cs *ContainerService) RunSpicedb() error {
	if err := cs.RunSpicedbMigration(); err != nil {
		return fmt.Errorf("failed to run spicedb migration: %w", err)
	}

	ctnr, err := cs.RunWithOptions(
		&dockertest.RunOptions{
			Name:       "authz-e2e-spicedb",
			Repository: "authzed/spicedb",
			Tag:        "v1.31.0",
			Env: []string{
				"SPICEDB_GRPC_PRESHARED_KEY=spicedb-super-secret",
				"SPICEDB_DATASTORE_ENGINE=postgres",
				fmt.Sprintf("SPICEDB_DATASTORE_CONN_URI=postgres://postgres:@%s:%s/postgres", cs.dockerHost, cs.PostgresPort),
			},
			ExposedPorts: []string{"50051/tcp"},
			Labels:       map[string]string{containerE2ELabel: "true"},
			Cmd:          []string{"serve", "--log-level", "warn"},
		},
		func(config *docker.HostConfig) {
			config.AutoRemove = true
			config.RestartPolicy = docker.RestartPolicy{Name: "no"}
		},
	)
	if err != nil {
		return fmt.Errorf("failed to start spicedb container: %w", err)
	}

	if err := ctnr.Expire(containerExpireSeconds); err != nil {
		return fmt.Errorf("failed to set expiration on spicedb container: %w", err)
	}

	cs.SpicedbCtnr = ctnr
	cs.SpicedbHostPort = ctnr.GetHostPort("50051/tcp")
	if err := cs.WaitForSpicedb(); err != nil {
		return fmt.Errorf("failed to wait for spicedb: %w", err)
	}
	return nil
}

func (cs *ContainerService) WaitForSpicedb() error {
	return backoff.Retry(func() error {
		conn, err := grpc.Dial(
			cs.SpicedbHostPort,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			return err
		}
		resp, err := healthpb.NewHealthClient(conn).
			Check(
				context.Background(),
				&healthpb.HealthCheckRequest{Service: "authzed.api.v1.SchemaService"},
			)
		if err != nil {
			return err
		}
		if resp.GetStatus() != healthpb.HealthCheckResponse_SERVING {
			return errors.New("no ready")
		}

		return nil
	},
		backoff.WithMaxRetries(backoff.NewConstantBackOff(1*time.Second), 20),
	)
}
