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
	Name   string
}

func NewBackupSet2(backupSetType string) *BackupSet2 {
	return &BackupSet2{
		Id:   uuid.New().String(),
		Type: backupSetType,
	}
}

func NewRepository2(name string, config *Config) *Repository2 {
	return &Repository2{
		col:    stor.NewCollection(),
		Name:   name,
		Config: config,
	}
}

func LoadRepository2(repo *Repository2, path string) error {
	indexPath := filepath.Join(path, "index")
	col := stor.Collection{}

	if err := stor.Deserialize(&col, indexPath); err != nil {
		return err
	}

	confPath := filepath.Join(path, "config")
	conf := Config{}
	if err := loadConfigFromRepo(&conf, confPath); err != nil {
		return err
	}

	repo.col = &col
	repo.Config = &conf

	return nil
}

func (r *Repository2) Init(path string) error {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return err
	}

	repoPath := filepath.Join(absPath, r.Name)
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

	if err = stor.Serialize(r.col, filepath.Join(r.Path, "index")); err != nil {
		return err
	}

	return nil
}

func (r *Repository2) AddBackupSet(backupSet *BackupSet2) error {
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

	if err := stor.Serialize(r.col, r.Path); err != nil {
		return err
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
