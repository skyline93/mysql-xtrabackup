package repo

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

type JsonBackupSet struct {
	Id      string `json:"id"`
	Path    string `json:"path"`
	Type    string `json:"type"`
	CycleId string `json:"cycle_id"`
}

type JsonBackupCycle struct {
	Id     string `json:"id"`
	RepoId string `json:"repo_id"`
}

type JsonRepo struct {
	BackupSets   []JsonBackupSet   `json:"backupsets"`
	BackupCycles []JsonBackupCycle `json:"backupcycles"`
}

func SerializeToJson(r *Repo, path string) error {
	var jsonBackupSets []JsonBackupSet
	var jsonBackupCycles []JsonBackupCycle

	for _, bc := range r.BackupCycles {
		for _, bs := range bc.BackupSets {
			jsonBackupSets = append(jsonBackupSets, JsonBackupSet{
				Id:      bs.Id,
				Path:    bs.Path,
				Type:    bs.Type,
				CycleId: bc.Id,
			})
		}

		jsonBackupCycles = append(jsonBackupCycles, JsonBackupCycle{
			Id:     bc.Id,
			RepoId: r.Id,
		})
	}

	jsonRepo := &JsonRepo{
		BackupSets:   jsonBackupSets,
		BackupCycles: jsonBackupCycles,
	}

	d, err := json.Marshal(jsonRepo)
	if err != nil {
		return err
	}

	os.WriteFile(path, d, 0664)
	return nil
}

func UnserializeFromJson(path string) (*Repo, error) {
	d, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	jsonRepo := &JsonRepo{}
	if err := json.Unmarshal(d, jsonRepo); err != nil {
		return nil, err
	}

	fileName := filepath.Base(path)
	repoId := strings.TrimSuffix(fileName, filepath.Ext(fileName))
	repo := NewRepo(repoId)
	for _, bc := range jsonRepo.BackupCycles {
		backupCycle := &BackupCycle{Id: bc.Id}
		for _, bs := range jsonRepo.BackupSets {
			if bs.CycleId == bc.Id {
				backupCycle.Insert(&BackupSet{Id: bs.Id, Path: bs.Path, Type: bs.Type})
			}
		}
		repo.Insert(backupCycle)
	}

	return repo, nil
}
