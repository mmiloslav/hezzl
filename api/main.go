package main

import (
	"natsq"
	"net/http"
	"os"
	"postgres"
	"redisdb"

	"github.com/sirupsen/logrus"
)

func main() {
	postgresConnParams, err := postgres.GetConnectionParams()
	if err != nil {
		logrus.Errorf("Error getting Connection Params [%s]", err.Error())
		os.Exit(1)
	}

	err = postgres.OpenConnection(postgresConnParams)
	if err != nil {
		logrus.Errorf("Error Open Postgres Connection [%s]", err.Error())
		os.Exit(1)
	}

	redisConnParams, err := redisdb.GetConnectionParams()
	if err != nil {
		logrus.Errorf("Error getting Redis Connection Params [%s]", err.Error())
		os.Exit(1)
	}

	err = redisdb.OpenConnection(redisConnParams)
	if err != nil {
		logrus.Errorf("Error Open Redis Connection [%s]", err.Error())
		os.Exit(1)
	}

	natsConnParams, err := natsq.GetConnectionParams()
	if err != nil {
		logrus.Errorf("Error getting Nats Connection Params [%s]", err.Error())
		os.Exit(1)
	}

	err = natsq.OpenConnection(natsConnParams)
	if err != nil {
		logrus.Errorf("Error Open Nats Connection [%s]", err.Error())
		os.Exit(1)
	}

	initValidator()

	// server
	router := NewRouter()

	logrus.Info("Starting server...")
	err = http.ListenAndServe(":8081", router)
	if err != nil {
		logrus.Errorf("Error ListenAndServe [%s]", err.Error())
		os.Exit(1)
	}
}
