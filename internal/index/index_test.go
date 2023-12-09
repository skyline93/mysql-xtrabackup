package index

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInsert(t *testing.T) {
	bs1 := NewBackupSet("/backup/set/path1", "full")
	bs2 := NewBackupSet("/backup/set/path2", "incr")
	bs3 := NewBackupSet("/backup/set/path3", "incr")

	bc1 := NewBackupCycle()

	bc1.Insert(bs1)
	assert.Equal(t, 1, len(bc1.BackupSets))

	bc1.Insert(bs2)
	bc1.Insert(bs3)
	assert.Equal(t, 3, len(bc1.BackupSets))

	assert.Equal(t, "full", bc1.Head().Type)
	assert.Equal(t, "/backup/set/path1", bc1.Head().Path)

	assert.Equal(t, "incr", bc1.Head().Next.Type)
	assert.Equal(t, "/backup/set/path2", bc1.Head().Next.Path)

	assert.Equal(t, "incr", bc1.Head().Next.Next.Type)
	assert.Equal(t, "/backup/set/path3", bc1.Head().Next.Next.Path)

	bs4 := NewBackupSet("/backup/set/path4", "full")
	bs5 := NewBackupSet("/backup/set/path5", "incr")
	bs6 := NewBackupSet("/backup/set/path6", "incr")

	bc2 := NewBackupCycle()
	bc2.Insert(bs4)
	bc2.Insert(bs5)
	bc2.Insert(bs6)

	repo := NewRepo("13579")
	repo.Insert(bc1)
	repo.Insert(bc2)
}
