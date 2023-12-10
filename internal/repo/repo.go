package repo

import (
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

const (
	TypeBackupSetFull = "full"
	TypeBackupSetIncr = "incr"
)

type BackupSet struct {
	Id      string
	Path    string
	Type    string
	FromLSN string
	ToLSN   string
	Prev    *BackupSet
	Next    *BackupSet
}

type BackupCycle struct {
	Id         string
	BackupSets []*BackupSet
	Prev       *BackupCycle
	Next       *BackupCycle
}

type Repo struct {
	Id           string         `json:"id"`
	Path         string         `json:"path"`
	BackupCycles []*BackupCycle `json:"backup_cycles"`
	Config       *Config        `json:"-"`
}

func NewBackupSet(backupSetType string) *BackupSet {
	return &BackupSet{
		Id:   uuid.New().String(),
		Type: backupSetType,
	}
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

func NewRepo(Id string, config *Config) *Repo {
	return &Repo{Id: Id, Config: config}
}

func (r *Repo) Init(path string) error {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return err
	}

	repoPath := filepath.Join(absPath, r.Id)
	r.Path = repoPath

	if err = os.MkdirAll(repoPath, 0764); err != nil {
		return err
	}

	if err = os.MkdirAll(filepath.Join(r.Path, "data"), 0764); err != nil {
		return err
	}

	if err = saveConfigToRepo(r.Config, r.Path); err != nil {
		return err
	}

	return nil
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

func (r *Repo) DataPath() string {
	return filepath.Join(r.Path, "data")
}

func (r *Repo) AddBackupSet(backupSet *BackupSet) {
	path := filepath.Join(r.DataPath(), backupSet.Id)
	backupSet.Path = path

	if backupSet.Type == TypeBackupSetFull {
		bc := NewBackupCycle()
		bc.Insert(backupSet)
		r.Insert(bc)
	} else {
		r.Tail().Insert(backupSet)
	}
}

func (r *Repo) Commit() error {
	return r.serialize()
}
