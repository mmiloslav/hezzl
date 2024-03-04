package postgres

import (
	"fmt"
)

func ProjectCreate(name string) error {
	sql := fmt.Sprintf("INSERT INTO projects (name) VALUES ('%s')", name)
	return postgresDB.Exec(sql).Error
}
