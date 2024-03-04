package postgres

import (
	"common"
	"fmt"
	"natsq"
	"time"

	"github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	postgresDB *gorm.DB
	listener   *pq.Listener
)

type connectionParams struct {
	User     string
	Password string
	Host     string
	Port     string
	DBName   string
}

func OpenConnection(connParams connectionParams) (err error) {
	logrus.Info("opening postgres connection...")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", connParams.Host, connParams.User, connParams.Password, connParams.DBName, connParams.Port)
	postgresDB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logrus.Errorf("error opening postgres gorm connection [%s]", err.Error())
		return err
	}
	logrus.Info("successfully opened postgres connection")

	err = makeListener(dsn)
	if err != nil {
		logrus.Errorf("error making listener [%s]", err.Error())
		return err
	}

	err = migrate()
	if err != nil {
		logrus.Errorf("Error migrating postgres tables [%s]", err.Error())
		return err
	}

	return nil
}

func makeListener(dsn string) error {
	listener = pq.NewListener(dsn, 10*time.Second, time.Minute, func(ev pq.ListenerEventType, err error) {
		if err != nil {
			logrus.Infof("pq.ListenerEventType error: [%s]", err.Error())
		}
	})
	err := listener.Listen("event")
	if err != nil {
		logrus.Errorf("error listening postgres events [%s]", err.Error())
		return err
	}

	go listen()

	return nil
}

func listen() {
	for {
		select {
		case <-time.After(5 * time.Second):
			logrus.Info("Waiting for notifications...")
		case n := <-listener.Notify:
			logrus.Infof("Received notification: [%s]", n.Extra)

			err := natsq.NatsConn.Publish("log-events", []byte(n.Extra))
			if err != nil {
				logrus.Errorf("Error publishing to NATS [%s] [%s]", n.Extra, err.Error())
			}
		}
	}
}

func GetConnectionParams() (connectionParams, error) {
	user, err := common.GetEnvVar("POSTGRES_USER")
	if err != nil {
		logrus.Errorf("error getting POSTGRES_USER var from env [%s]", err.Error())
		return connectionParams{}, err
	}

	password, err := common.GetEnvVar("POSTGRES_PASSWORD")
	if err != nil {
		logrus.Errorf("error getting POSTGRES_PASSWORD var from env [%s]", err.Error())
		return connectionParams{}, err
	}

	host, err := common.GetEnvVar("POSTGRES_HOST")
	if err != nil {
		logrus.Errorf("error getting POSTGRES_HOST var from env [%s]", err.Error())
		return connectionParams{}, err
	}

	port, err := common.GetEnvVar("POSTGRES_PORT")
	if err != nil {
		logrus.Errorf("error getting POSTGRES_PORT var from env [%s]", err.Error())
		return connectionParams{}, err
	}

	dbName, err := common.GetEnvVar("POSTGRES_DB")
	if err != nil {
		logrus.Errorf("error getting POSTGRES_DB var from env [%s]", err.Error())
		return connectionParams{}, err
	}

	return connectionParams{
		User:     user,
		Password: password,
		Host:     host,
		Port:     port,
		DBName:   dbName,
	}, nil
}
