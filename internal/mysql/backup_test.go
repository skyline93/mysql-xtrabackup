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

	bs1 := repository.NewBackupSet2(repository.TypeBackupSetFull)
	bs2 := repository.NewBackupSet2(repository.TypeBackupSetIncr)
	bs3 := repository.NewBackupSet2(repository.TypeBackupSetIncr)

	r := repository.NewRepository2("MYTEST1", config)
	err := r.Init("./")
	assert.Nil(t, err)
	assert.Equal(t, "MYTEST1", r.Name)

	r.AddBackupSet(bs1)
	r.AddBackupSet(bs2)
	r.AddBackupSet(bs3)
}
