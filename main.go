package main

import (
	"context"
	"log"
	"net/http"

	"rift/authz/client"
	"rift/httpsrv"
	"rift/memdb"
)

func main() {
	db := memdb.New()

	// NOTE: start authz in docker-compose before running this
	authzC, err := client.New("localhost:50051", "spicedb-super-secret")
	if err != nil {
		log.Fatal(err)
	}

	if err := authzC.MigrateSchema(context.Background()); err != nil {
		log.Fatal(err)
	}

	srv := httpsrv.New(db, authzC)
	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", srv))
}
