package postgres

import "github.com/sirupsen/logrus"

func migrate() error {
	logrus.Info("migrating tables...")
	err := postgresDB.AutoMigrate(&Project{}, &Good{})
	if err != nil {
		logrus.Errorf("Error initial migraion [%s]", err.Error())
		return err
	}

	// TODO не сохранять если есть
	err = ProjectCreate("Первая запись")
	if err != nil {
		logrus.Errorf("Error creating project [%s]", err.Error())
		return err
	}

	logrus.Info("successfully migrated needed migrations")
	return nil
}
