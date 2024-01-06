package mysql

import (
	"github.com/skyline93/mysql-xtrabackup/internal/repository"
)

type Restorer struct {
}

func NewRestorer() *Restorer {
	return &Restorer{}
}

func (r *Restorer) Restore(repo *repository.Repository, targetPath string, backupSetId string) error {
	return nil
}
