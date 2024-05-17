package e2e

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"rift/authz/client"
	"rift/httpsrv"
	"rift/memdb"
)

var tsURL string

func TestMain(m *testing.M) {
	if os.Getenv("E2E") != "true" {
		fmt.Println("skipping e2e tests, if you want to run them set env E2E=true")
		return
	}

	// containers
	containerSrv := die2(NewContainerService())
	containerSrv.HandleInterrupt()
	die(containerSrv.Cleanup())
	die(containerSrv.RunPostgres())
	die(containerSrv.RunSpicedb())

	// clients
	db := memdb.New()
	authzC := die2(client.New(containerSrv.SpicedbHostPort, "spicedb-super-secret"))
	die(authzC.MigrateSchema(context.Background()))

	// http rest
	mux := httpsrv.New(db, authzC)
	ts := httptest.NewUnstartedServer(mux)
	if os.Getenv("CI") == "true" {
		// if tests are run in github actions then bind to all ports
		// so local docker containers can connect to the server
		ts = &httptest.Server{
			Listener: die2(net.Listen("tcp", ":0")),
			Config:   &http.Server{Handler: mux},
		}
	}
	ts.Start()
	tsURL = ts.URL
	defer ts.Close()
	os.Exit(m.Run())
}

func die(err error) {
	if err != nil {
		panic(err)
	}
}

func die2[T any](t T, err error) T {
	if err != nil {
		panic(err)
	}
	return t
}
