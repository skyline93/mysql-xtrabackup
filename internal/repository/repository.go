package repository

import (
	"encoding/json"
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
	repo.Name = conf.Identifer
	repo.Path = path

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
	v, err := json.Marshal(backupSet)
	if err != nil {
		return err
	}

	if backupSet.Type == TypeBackupSetFull {
		_, err := r.col.NewNode(backupSet.Id, v, true)
		if err != nil {
			return err
		}
	} else if backupSet.Type == TypeBackupSetIncr {
		_, err := r.col.NewNode(backupSet.Id, v, false)
		if err != nil {
			return err
		}
	}

	if err := stor.Serialize(r.col, filepath.Join(r.Path, "index")); err != nil {
		return err
	}

	return nil
}

func (r *Repository2) GetBackupSet(backupSetId string) (*BackupSet2, error) {
	n := r.col.GetNode(backupSetId)
	var backupSet BackupSet2
	if err := json.Unmarshal(n.Data, &backupSet); err != nil {
		return nil, err
	}

	return &backupSet, nil
}

func (r *Repository2) GetBeforeBackupSet(backupSetId string) ([]BackupSet2, error) {
	var backupSets []BackupSet2

	nodes := r.col.GetBeforeNodes(backupSetId)

	for _, n := range nodes {
		var backupSet BackupSet2
		if err := json.Unmarshal(n.Data, &backupSet); err != nil {
			return nil, err
		}

		backupSets = append(backupSets, backupSet)
	}

	return backupSets, nil
}

func (r *Repository2) GetLastBackupSet() (*BackupSet2, error) {
	n := r.col.GetLastNode()
	if n == nil {
		return nil, errors.New("last backupset is not found")
	}

	var backupSet BackupSet2
	if err := json.Unmarshal(n.Data, &backupSet); err != nil {
		return nil, err
	}
	return &backupSet, nil
}

func (r *Repository2) DataPath() string {
	return filepath.Join(r.Path, "data")
}

func (r *Repository2) ListBackupSets() ([]BackupSet2, error) {
	var backupSets []BackupSet2

	ns := r.col.GetAllNodes()

	for _, n := range ns {
		var backupSet BackupSet2
		if err := json.Unmarshal(n.Data, &backupSet); err != nil {
			return nil, err
		}

		backupSets = append(backupSets, backupSet)
	}

	return backupSets, nil
}
