package repository

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/skyline93/mysql-xtrabackup/internal/stor"
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
	Size    int64
	Prev    *BackupSet
	Next    *BackupSet
}

type BackupCycle struct {
	Id         string
	BackupSets []*BackupSet
	Prev       *BackupCycle
	Next       *BackupCycle
}

type Repository struct {
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

func (bc *BackupCycle) Find(backupSetId string) (*BackupSet, error) {
	bs := bc.Head()

	for {
		if bs.Id == backupSetId {
			return bs, nil
		}

		if bs.Next == nil {
			return nil, errors.New("not found")
		}

		bs = bs.Next
	}
}

func NewRepository(Id string, config *Config) *Repository {
	return &Repository{Id: Id, Config: config}
}

func (r *Repository) Init(path string) error {
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

func (r *Repository) Head() *BackupCycle {
	return r.BackupCycles[0]
}

func (r *Repository) Tail() *BackupCycle {
	// TODO 初始没有index文件的问题
	if r.BackupCycles == nil {
		return nil
	}
	return r.BackupCycles[len(r.BackupCycles)-1]
}

func (r *Repository) Insert(backupCycle *BackupCycle) {
	if len(r.BackupCycles) == 0 {
		r.BackupCycles = append(r.BackupCycles, backupCycle)
	} else {
		backupCycle.Prev = r.Tail()
		r.Tail().Next = backupCycle

		r.BackupCycles = append(r.BackupCycles, backupCycle)
	}
}

func (r *Repository) DataPath() string {
	return filepath.Join(r.Path, "data")
}

func (r *Repository) AddBackupSet(backupSet *BackupSet) {
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

func (r *Repository) Commit() error {
	return r.serialize()
}

func (r *Repository) FindBackupSet(backupSetId string) (*BackupSet, error) {
	bc := r.Head()

	for {
		bs, err := bc.Find(backupSetId)
		if err != nil {
			if bc.Next != nil {
				bc = bc.Next
				continue
			}
			return nil, errors.New("not found")
		}

		if bs.Id == backupSetId {
			return bs, nil
		}
	}
}

func (r *Repository) FindBackupCycle(backupSetId string) (*BackupCycle, error) {
	bc := r.Head()
	for {
		_, err := bc.Find(backupSetId)
		if err != nil {
			bc = bc.Next
			continue
		}

		return bc, nil
	}
}

type BackupSet2 struct {
	Id      string
	Path    string
	Type    string
	FromLSN string
	ToLSN   string
	Size    int64
}

type Repository2 struct {
	col    *stor.Collection
	Config *Config
	Path   string
}

func NewBackupSet2(backupSetType string) *BackupSet2 {
	return &BackupSet2{
		Id:   uuid.New().String(),
		Type: backupSetType,
	}
}

func NewRepository2() *Repository2 {
	return &Repository2{col: stor.NewCollection()}
}

func (r *Repository2) AddBackupSet(backupSet *BackupSet) error {
	if backupSet.Type == TypeBackupSetFull {
		_, err := r.col.NewNode(backupSet.Id, backupSet, true)
		if err != nil {
			return err
		}
	} else if backupSet.Type == TypeBackupSetIncr {
		_, err := r.col.NewNode(backupSet.Id, backupSet, false)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *Repository2) GetBackupSet(backupSetId string) (*BackupSet2, error) {
	n := r.col.GetNode(backupSetId)
	backupSet, ok := n.Data.(*BackupSet2)
	if !ok {
		return nil, errors.New("the backup set is not found")
	}

	return backupSet, nil
}

func (r *Repository2) GetBeforeBackupSet(backupSetId string) ([]BackupSet2, error) {
	var backupSets []BackupSet2

	nodes := r.col.GetBeforeNodes(backupSetId)

	for _, n := range nodes {
		backupSet, ok := n.Data.(*BackupSet2)
		if !ok {
			return nil, errors.New("unknow error")
		}

		backupSets = append(backupSets, *backupSet)
	}

	return backupSets, nil
}

func (r *Repository2) GetLastBackupSet() (*BackupSet2, error) {
	n := r.col.GetLastNode()
	if n == nil {
		return nil, errors.New("last backupset is not found")
	}

	return n.Data.(*BackupSet2), nil
}

func (r *Repository2) DataPath() string {
	return filepath.Join(r.Path, "data")
}
