package repo

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSerialize(t *testing.T) {
	bs1 := NewBackupSet("/backup/set/path1", TypeBackupSetFull)
	bs2 := NewBackupSet("/backup/set/path2", "incr")
	bs3 := NewBackupSet("/backup/set/path3", "incr")

	bc1 := NewBackupCycle()

	bc1.Insert(bs1)
	bc1.Insert(bs2)
	bc1.Insert(bs3)

	bs4 := NewBackupSet("/backup/set/path4", TypeBackupSetFull)
	bs5 := NewBackupSet("/backup/set/path5", TypeBackupSetIncr)
	bs6 := NewBackupSet("/backup/set/path6", TypeBackupSetIncr)

	bc2 := NewBackupCycle()

	bc2.Insert(bs4)
	bc2.Insert(bs5)
	bc2.Insert(bs6)

	bs7 := NewBackupSet("/backup/set/path7", TypeBackupSetFull)
	bs8 := NewBackupSet("/backup/set/path8", TypeBackupSetIncr)
	bs9 := NewBackupSet("/backup/set/path9", TypeBackupSetIncr)
	bc3 := NewBackupCycle()

	bc3.Insert(bs7)
	bc3.Insert(bs8)
	bc3.Insert(bs9)

	repo := NewRepo("24680")
	repo.Insert(bc1)
	repo.Insert(bc2)
	repo.Insert(bc3)

	bs10 := NewBackupSet("/backup/set/path10", TypeBackupSetFull)
	bs11 := NewBackupSet("/backup/set/path11", TypeBackupSetIncr)
	repo.InsertBackupSet(bs10)
	repo.InsertBackupSet(bs11)

	p, _ := filepath.Abs("./")
	filePath := filepath.Join(p, fmt.Sprintf("%s.json", repo.Id))
	err := SerializeToJson(repo, filePath)
	assert.Nil(t, err)
}

func TestUnserialize(t *testing.T) {
	p, _ := filepath.Abs("./")
	filePath := filepath.Join(p, "24680.json")
	_, err := UnserializeFromJson(filePath)
	assert.Nil(t, err)
}
