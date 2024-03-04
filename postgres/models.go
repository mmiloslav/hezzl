package postgres

import (
	"time"
)

type Project struct {
	ID        int       `gorm:"primaryKey"`
	Name      string    `gorm:"type:varchar(100)"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}

type GoodSlice []Good
type Good struct {
	ID          int     `gorm:"primaryKey"`
	ProjectID   int     `gorm:"not null"`
	Project     Project `gorm:"foreignKey:ProjectID"`
	Name        string  `gorm:"type:varchar(100)"`
	Description string  `gorm:"type:varchar(255)"`
	Priority    int
	Removed     bool      `gorm:"default:false"`
	CreatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}
