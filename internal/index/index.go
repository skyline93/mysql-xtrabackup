package index

import (
	"encoding/json"
	"os"

	"github.com/google/uuid"
)

type BackupSet struct {
	ID   string     `json:"id"`
	Path string     `json:"path"`
	Type string     `json:"type"`
	Prev *BackupSet `json:"-"`
	Next *BackupSet `json:"-"`
}

func (bs *BackupSet) MarshalJSON() ([]byte, error) {
	prevID := ""
	nextID := ""

	if bs.Prev != nil {
		prevID = bs.Prev.ID
	}

	if bs.Next != nil {
		nextID = bs.Next.ID
	}

	return json.Marshal(&struct {
		ID   string `json:"id"`
		Path string `json:"path"`
		Type string `json:"type"`
		Prev string `json:"prev"`
		Next string `json:"next"`
	}{
		ID:   bs.ID,
		Path: bs.Path,
		Type: bs.Type,
		Prev: prevID,
		Next: nextID,
	})
}

func (bs *BackupSet) UnmarshalJSON(data []byte) error {
	aux := struct {
		ID   string `json:"id"`
		Path string `json:"path"`
		Type string `json:"type"`
		Prev string `json:"prev"`
		Next string `json:"next"`
	}{}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	bs.ID = aux.ID
	bs.Path = aux.Path
	bs.Type = aux.Type

	// Assuming you have a function to retrieve BackupSet by ID
	bs.Prev = GetBackupSetByID(aux.Prev)
	bs.Next = GetBackupSetByID(aux.Next)

	return nil
}

func NewBackupSet(path string, backupSetType string) *BackupSet {
	return &BackupSet{
		ID:   uuid.New().String(),
		Path: path,
		Type: backupSetType,
	}
}

type BackupCycle struct {
	BackupSets []*BackupSet
	Previous   *BackupCycle
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
		bc.BackupSets = append(bc.BackupSets, backupSet)
	} else {
		backupSet.Prev = bc.Tail()
		bc.Tail().Next = backupSet

		bc.BackupSets = append(bc.BackupSets, backupSet)
	}
}

func (bc *BackupCycle) UnmarshalJSON(data []byte) error {
	var aux struct {
		BackupSets []*BackupSet `json:"backupSets"`
		Previous   string       `json:"previous"`
		Next       string       `json:"next"`
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	bc.BackupSets = aux.BackupSets

	// Assuming you have a function to retrieve BackupCycle by ID
	bc.Previous = GetBackupCycleByID(aux.Previous)
	bc.Next = GetBackupCycleByID(aux.Next)

	return nil
}

// Function to retrieve BackupCycle by ID
func GetBackupCycleByID(id string) *BackupCycle {
	// Implement this function based on your data structure.
	// It should return the BackupCycle instance with the given ID.
	return nil
}

func LoadBackupCycleFromJSONFile(filename string) (*BackupCycle, error) {
	fileData, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var bc *BackupCycle
	err = json.Unmarshal(fileData, &bc)
	if err != nil {
		return nil, err
	}

	return bc, nil
}

func SaveBackupCycleToJSONFile(bc *BackupCycle, filename string) error {
	data, err := json.MarshalIndent(bc, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filename, data, 0644)
}
