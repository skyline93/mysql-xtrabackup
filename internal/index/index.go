package index

import (
	"github.com/google/uuid"
)

type BackupSet struct {
	Id   string
	Path string
	Type string
	Prev *BackupSet
	Next *BackupSet
}

func NewBackupSet(path string, backupSetType string) *BackupSet {
	return &BackupSet{
		Id:   uuid.New().String(),
		Path: path,
		Type: backupSetType,
	}
}

type BackupCycle struct {
	Id         string
	BackupSets []*BackupSet
	Prev       *BackupCycle
	Next       *BackupCycle
}

func NewBackupCycle() *BackupCycle {
	return &BackupCycle{}
}

func (bc *BackupCycle) Head() *BackupSet {
	return bc.BackupSets[0]
}

func (bc *BackupCycle) Tail() *BackupSet {
	return bc.BackupSets[len(bc.BackupSets)-1]
}

func (bc *BackupCycle) Insert(backupSet *BackupSet) {
	if len(bc.BackupSets) == 0 {
		bc.Id = backupSet.Id
		bc.BackupSets = append(bc.BackupSets, backupSet)
	} else {
		backupSet.Prev = bc.Tail()
		bc.Tail().Next = backupSet

		bc.BackupSets = append(bc.BackupSets, backupSet)
	}
}

type Repo struct {
	Id           string         `json:"id"`
	BackupCycles []*BackupCycle `json:"backup_cycles"`
}

func NewRepo(Id string) *Repo {
	return &Repo{Id: Id}
}

func (r *Repo) Head() *BackupCycle {
	return r.BackupCycles[0]
}

func (r *Repo) Tail() *BackupCycle {
	return r.BackupCycles[len(r.BackupCycles)-1]
}

func (r *Repo) Insert(backupCycle *BackupCycle) {
	if len(r.BackupCycles) == 0 {
		r.BackupCycles = append(r.BackupCycles, backupCycle)
	} else {
		backupCycle.Prev = r.Tail()
		r.Tail().Next = backupCycle

		r.BackupCycles = append(r.BackupCycles, backupCycle)
	}
}
