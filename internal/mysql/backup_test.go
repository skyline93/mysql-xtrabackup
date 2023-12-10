package mysql

import (
	"testing"

	"github.com/skyline93/mysql-xtrabackup/internal/repo"
	"github.com/stretchr/testify/assert"
)

func TestInitRepo(t *testing.T) {
	config := &repo.Config{
		Identifer:   "MYTEST1",
		Version:     "8.0.23",
		LoginPath:   "MYTEST1",
		DbHostName:  "localhost",
		DbUser:      "root",
		Throttle:    400,
		TryCompress: true,
	}

	bs1 := repo.NewBackupSet(repo.TypeBackupSetFull)
	bs2 := repo.NewBackupSet(repo.TypeBackupSetIncr)
	bs3 := repo.NewBackupSet(repo.TypeBackupSetIncr)

	r := repo.NewRepo("MYTEST1", config)
	err := r.Init("./")
	assert.Nil(t, err)
	assert.Equal(t, "MYTEST1", r.Id)

	r.AddBackupSet(bs1)
	r.AddBackupSet(bs2)
	r.AddBackupSet(bs3)

	r.Commit()
}
