package postgres

import (
	"encoding/json"
	"fmt"
	"redisdb"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var errHookNoRows = "sql: Scan error on column index 0, name \"max\": converting NULL to int is unsupported"

func (m *Good) Create() error {
	return postgresDB.Create(&m).Error
}

func (m *Good) Update(name, description string) error {
	err := postgresDB.Where(&m).First(&m).Error
	if err != nil {
		logrus.Errorf("error finding good by id=[%d], projectID=[%d] [%s]", m.ID, m.ProjectID, err.Error())
		return err
	}

	m.Name = name
	m.Description = description

	return postgresDB.Save(&m).Error
}

func (m *Good) Delete() error {
	tx := postgresDB.Begin()
	if tx.Error != nil {
		logrus.Errorf("error beginning transaction [%s]", tx.Error.Error())
		return tx.Error
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	err := tx.Exec("SET TRANSACTION ISOLATION LEVEL SERIALIZABLE").Error
	if err != nil {
		logrus.Errorf("error setting transaction level [%s]", err.Error())
		return err
	}

	err = tx.Where(&m).First(&m).Error
	if err != nil {
		logrus.Errorf("error finding good by id=[%d], projectID=[%d] [%s]", m.ID, m.ProjectID, err.Error())
		tx.Rollback()
		return err
	}

	m.Removed = true

	err = tx.Save(&m).Error
	if err != nil {
		logrus.Errorf("error deleting good with id=[%d], projectID=[%d] [%s]", m.ID, m.ProjectID, err.Error())
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (m *GoodSlice) Many(limit, offset int) (int, error) {
	allGoods := GoodSlice{}

	redisAllGoods, err := redisdb.RedisClient.Get(redisdb.AllGoodsKey).Bytes()
	if err != nil {
		// geting from postgres
		logrus.Info("all goods from postgres")

		err = postgresDB.Find(&allGoods).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			logrus.Errorf("error finding goods in db [%s]", err.Error())
			return -1, err
		}

		// caching to redis
		err = redisdb.Cache(redisdb.AllGoodsKey, allGoods)
		if err != nil {
			logrus.Errorf("error caching goods in redis [%s]", err.Error())
			return -1, err
		}
	} else {
		logrus.Info("all goods from redis")

		err := json.Unmarshal(redisAllGoods, &allGoods)
		if err != nil {
			logrus.Errorf("error unmarshal allGoods from redis [%s]", err.Error())
			return -1, err
		}
	}

	redisLimitOffsetGoods, err := redisdb.RedisClient.Get(fmt.Sprintf("goods_%d_%d", limit, offset)).Bytes()
	if err != nil {
		// geting from postgres
		logrus.Info("limit offset goods from postgres")

		err = postgresDB.Limit(limit).Offset(offset).Order("id").Find(&m).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			logrus.Errorf("error finding goods with limit and offset in db [%s]", err.Error())
			return -1, err
		}

		// caching to redis
		err = redisdb.Cache(fmt.Sprintf("goods_%d_%d", limit, offset), &m)
		if err != nil {
			logrus.Errorf("error caching goods in redis [%s]", err.Error())
			return -1, err
		}
	} else {
		logrus.Info("limit offset goods from redis")
		// getting from redis
		err := json.Unmarshal(redisLimitOffsetGoods, m)
		if err != nil {
			logrus.Errorf("error unmarshal limit offset goods from redis [%s]", err.Error())
			return -1, err
		}
	}

	return len(allGoods), nil
}

func (m *GoodSlice) ManyByPriority(priority int) error {
	return postgresDB.Where("priority >= ?", priority).Order("priority").Find(&m).Error
}

func (m *GoodSlice) Reprioritize(newPriority int, good *Good) error {
	err := postgresDB.Where(&good).First(&good).Error
	if err != nil {
		logrus.Errorf("error finding good by id=[%d], projectID=[%d] [%s]", good.ID, good.ProjectID, err.Error())
		return err
	}

	err = m.ManyByPriority(good.Priority)
	if err != nil {
		logrus.Errorf("error getting goods by priority [%s]", err.Error())
		return err
	}

	tx := postgresDB.Begin()
	if tx.Error != nil {
		logrus.Errorf("error beginning transaction [%s]", tx.Error.Error())
		return tx.Error
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	err = tx.Exec("SET TRANSACTION ISOLATION LEVEL SERIALIZABLE").Error
	if err != nil {
		logrus.Errorf("error setting transaction level [%s]", err.Error())
		return err
	}

	nextPriority := newPriority
	for _, elem := range *m {
		elem.Priority = nextPriority
		err := tx.Model(&elem).Updates(elem).Error
		if err != nil {
			logrus.Errorf("error saving good [%d] with new piority [%d] [%s]", elem.ID, elem.Priority, err.Error())
			tx.Rollback()
			return err
		}

		nextPriority++
	}

	err = tx.Commit().Error
	if err != nil {
		logrus.Errorf("error tx commit [%s]", err.Error())
		tx.Rollback()
		return err
	}

	return m.ManyByPriority(newPriority)
}

func (m *Good) BeforeCreate(tx *gorm.DB) error {
	var maxPriority int
	sql := "SELECT MAX(priority) FROM goods"
	err := postgresDB.Raw(sql).Scan(&maxPriority).Error
	if err != nil && err.Error() != errHookNoRows {
		logrus.Errorf("error hook create good [%s]", err.Error())
		return err
	}

	m.Priority = maxPriority + 1
	return nil
}
