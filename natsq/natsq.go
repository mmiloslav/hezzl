package natsq

import (
	"common"

	"github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"
)

var NatsConn *nats.Conn

type connectionParams struct {
	Url string
}

func OpenConnection(connParams connectionParams) (err error) {
	NatsConn, err = nats.Connect(connParams.Url)
	if err != nil {
		logrus.Errorf("error connectiong to nats [%s]", err.Error())
		return err
	}

	return nil
}

func GetConnectionParams() (connectionParams, error) {
	url, err := common.GetEnvVar("NATS_URL")
	if err != nil {
		logrus.Errorf("error getting NATS_URL var from env [%s]", err.Error())
		return connectionParams{}, err
	}

	return connectionParams{
		Url: url,
	}, nil
}
