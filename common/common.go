package common

import (
	"errors"
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
)

func GetEnvVar(name string) (string, error) {
	val, ok := os.LookupEnv(name)
	if !ok {
		errorText := fmt.Sprintf("[%s] not found in env", name)
		logrus.Error(errorText)
		return "", errors.New(errorText)
	}

	return val, nil
}
