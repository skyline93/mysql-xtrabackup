package mysql

import (
	"testing"

	"github.com/skyline93/mysql-xtrabackup/internal/repository"
	"github.com/stretchr/testify/assert"
)

func TestInitrepository(t *testing.T) {
	config := &repository.Config{
		Identifer:   "MYTEST1",
		Version:     "8.0.23",
		LoginPath:   "MYTEST1",
		DbHostName:  "localhost",
		DbUser:      "root",
		Throttle:    400,
		TryCompress: true,
	}

	bs1 := repository.NewBackupSet(repository.TypeBackupSetFull)
	bs2 := repository.NewBackupSet(repository.TypeBackupSetIncr)
	bs3 := repository.NewBackupSet(repository.TypeBackupSetIncr)

	r := repository.NewRepository("MYTEST1", config)
	err := r.Init("./")
	assert.Nil(t, err)
	assert.Equal(t, "MYTEST1", r.Id)

	r.AddBackupSet(bs1)
	r.AddBackupSet(bs2)
	r.AddBackupSet(bs3)

	r.Commit()
}
