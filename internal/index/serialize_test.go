package index

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSerialize(t *testing.T) {
	bs1 := NewBackupSet("/backup/set/path1", "full")
	bs2 := NewBackupSet("/backup/set/path2", "incr")
	bs3 := NewBackupSet("/backup/set/path3", "incr")

	bc1 := NewBackupCycle()

	bc1.Insert(bs1)
	bc1.Insert(bs2)
	bc1.Insert(bs3)

	bs4 := NewBackupSet("/backup/set/path4", "full")
	bs5 := NewBackupSet("/backup/set/path5", "incr")
	bs6 := NewBackupSet("/backup/set/path6", "incr")

	bc2 := NewBackupCycle()

	bc2.Insert(bs4)
	bc2.Insert(bs5)
	bc2.Insert(bs6)

	bs7 := NewBackupSet("/backup/set/path4", "full")
	bs8 := NewBackupSet("/backup/set/path5", "incr")
	bs9 := NewBackupSet("/backup/set/path6", "incr")
	bc3 := NewBackupCycle()

	bc3.Insert(bs7)
	bc3.Insert(bs8)
	bc3.Insert(bs9)

	repo := NewRepo("24680")
	repo.Insert(bc1)
	repo.Insert(bc2)
	repo.Insert(bc3)

	err := SaveToFile(repo, "./")
	assert.Nil(t, err)

	// os.Remove(fmt.Sprintf("%s.json", repo.Id))
}

func TestUnserialize(t *testing.T) {
	var repo *Repo
	err := LoadFromFile(repo, "24680.json")
	assert.Nil(t, err)
}
