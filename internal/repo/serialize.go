package repo

import (
	"encoding/json"
	"os"
	"path/filepath"
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

func (r *Repo) serialize() error {
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

	return os.WriteFile(filepath.Join(r.Path, "index"), d, 0664)
}

func Load(repoPath string) (*Repo, error) {
	d, err := os.ReadFile(filepath.Join(repoPath, "index"))
	if err != nil {
		return nil, err
	}

	jsonRepo := &JsonRepo{}
	if err := json.Unmarshal(d, jsonRepo); err != nil {
		return nil, err
	}

	config, err := loadConfigFromRepo(repoPath)
	if err != nil {
		return nil, err
	}

	repoId := filepath.Base(repoPath)
	repo := NewRepo(repoId, config)
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

// func SerializeToJson(r *Repo, path string) error {
// 	var jsonBackupSets []JsonBackupSet
// 	var jsonBackupCycles []JsonBackupCycle

// 	for _, bc := range r.BackupCycles {
// 		for _, bs := range bc.BackupSets {
// 			jsonBackupSets = append(jsonBackupSets, JsonBackupSet{
// 				Id:      bs.Id,
// 				Path:    bs.Path,
// 				Type:    bs.Type,
// 				CycleId: bc.Id,
// 			})
// 		}

// 		jsonBackupCycles = append(jsonBackupCycles, JsonBackupCycle{
// 			Id:     bc.Id,
// 			RepoId: r.Id,
// 		})
// 	}

// 	jsonRepo := &JsonRepo{
// 		BackupSets:   jsonBackupSets,
// 		BackupCycles: jsonBackupCycles,
// 	}

// 	d, err := json.Marshal(jsonRepo)
// 	if err != nil {
// 		return err
// 	}

// 	return os.WriteFile(path, d, 0664)
// }

// func UnserializeFromJson(path string) (*Repo, error) {
// 	d, err := os.ReadFile(path)
// 	if err != nil {
// 		return nil, err
// 	}

// 	jsonRepo := &JsonRepo{}
// 	if err := json.Unmarshal(d, jsonRepo); err != nil {
// 		return nil, err
// 	}

// 	fileName := filepath.Base(path)
// 	repoId := strings.TrimSuffix(fileName, filepath.Ext(fileName))
// 	repo := NewRepo(repoId)
// 	for _, bc := range jsonRepo.BackupCycles {
// 		backupCycle := &BackupCycle{Id: bc.Id}
// 		for _, bs := range jsonRepo.BackupSets {
// 			if bs.CycleId == bc.Id {
// 				backupCycle.Insert(&BackupSet{Id: bs.Id, Path: bs.Path, Type: bs.Type})
// 			}
// 		}
// 		repo.Insert(backupCycle)
// 	}

// 	return repo, nil
// }
