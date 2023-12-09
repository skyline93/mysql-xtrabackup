package index

import (
	"encoding/gob"
	"encoding/json"
	"fmt"
	"os"
)

var load_helper map[string]interface{}

func addr(p interface{}) string {
	return fmt.Sprintf("%p", p)
}

type JsonBackupSet struct {
	Addr string `json:"addr"`
	Id   string `json:"id"`
	Path string `json:"path"`
	Type string `json:"type"`
	Prev string `json:"prev"`
	Next string `json:"next"`
}

func (bs *BackupSet) MarshalJSON() ([]byte, error) {
	prevId := ""
	nextId := ""

	if bs.Prev != nil {
		prevId = addr(bs.Prev)
	}

	if bs.Next != nil {
		nextId = addr(bs.Next)
	}

	return json.Marshal(&JsonBackupSet{
		Addr: addr(bs),
		Id:   bs.Id,
		Path: bs.Path,
		Type: bs.Type,
		Prev: prevId,
		Next: nextId,
	})
}

func (bs *BackupSet) UnmarshalJSON(data []byte) error {
	aux := &JsonBackupSet{}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	bs.Id = aux.Id
	bs.Path = aux.Path
	bs.Type = aux.Type

	load_helper[aux.Addr] = bs
	return nil
}

type JsonBackupCycle struct {
	Addr       string   `json:"addr"`
	Id         string   `json:"id"`
	BackupSets []string `json:"backupsets"`
	Prev       string   `json:"prev"`
	Next       string   `json:"next"`
}

func (bc *BackupCycle) MarshalJSON() ([]byte, error) {
	prevId := ""
	nextId := ""

	if bc.Prev != nil {
		prevId = addr(bc.Prev)
	}

	if bc.Next != nil {
		nextId = addr(bc.Next)
	}

	var record_ids []string
	for _, r := range bc.BackupSets {
		record_ids = append(record_ids, addr(r))
	}

	return json.Marshal(&JsonBackupCycle{
		Addr:       addr(bc),
		Id:         bc.Id,
		BackupSets: record_ids,
		Prev:       prevId,
		Next:       nextId,
	})
}

func (bc *BackupCycle) UnmarshalJSON(data []byte) error {
	aux := &JsonBackupCycle{}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	bc.Id = aux.Id
	for _, record_id := range aux.BackupSets {
		bc.BackupSets = append(bc.BackupSets, load_helper[record_id].(*BackupSet))
	}
	load_helper[aux.Addr] = bc

	return nil
}

type JsonRepo struct {
	Addr         string   `json:"addr"`
	Id           string   `json:"id"`
	BackupCycles []string `json:"backupcycles"`
}

func (r *Repo) MarshalJSON() ([]byte, error) {
	var record_ids []string
	for _, r := range r.BackupCycles {
		record_ids = append(record_ids, addr(r))
	}

	return json.Marshal(&JsonRepo{
		Addr:         addr(r),
		Id:           r.Id,
		BackupCycles: record_ids,
	})
}

func (r *Repo) UnmarshalJSON(data []byte) error {
	aux := &JsonRepo{}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	r.Id = aux.Id
	for _, record_id := range aux.BackupCycles {
		r.BackupCycles = append(r.BackupCycles, load_helper[record_id].(*BackupCycle))
	}
	load_helper[aux.Addr] = r
	return nil
}

type State struct {
	BackupSets   map[string]*BackupSet
	BackupCycles map[string]*BackupCycle
	Repos        map[string]*Repo
}

func NewState() *State {
	s := &State{}
	s.BackupSets = make(map[string]*BackupSet)
	s.BackupCycles = make(map[string]*BackupCycle)
	s.Repos = make(map[string]*Repo)
	return s
}

func (s *State) AddBackupSet(backupSet *BackupSet) {
	s.BackupSets[addr(backupSet)] = backupSet
}

func (s *State) AddBackupCycle(backupCycle *BackupCycle) {
	s.BackupCycles[addr(backupCycle)] = backupCycle
}

func (s *State) AddRepo(repo *Repo) {
	s.Repos[addr(repo)] = repo
}

// func SaveRepo(repo *Repo, path string) error {
// 	state := NewState()
// 	state.AddRepo(repo)

// 	for _, bc := range repo.BackupCycles {
// 		state.AddBackupCycle(bc)
// 		for _, bs := range bc.BackupSets {
// 			state.AddBackupSet(bs)
// 		}
// 	}

// 	d, err := json.MarshalIndent(state, "", "    ")
// 	if err != nil {
// 		return err
// 	}

// 	os.WriteFile(path, []byte(d), 0664)
// 	return nil
// }

// func LoadRepo(path string) (*Repo, error) {
// 	d, err := os.ReadFile(path)
// 	if err != nil {
// 		return nil, err
// 	}

// 	load_helper = make(map[string]interface{})
// 	state := NewState()
// 	if err := json.Unmarshal(d, state); err != nil {
// 		return nil, err
// 	}

// 	maps := make(map[string]interface{})
// 	if err = json.Unmarshal(d, &maps); err != nil {
// 		return nil, err
// 	}

// 	repo := &Repo{}
// 	for _, r := range state.Repos {
// 		repo = r
// 	}

// 	// for _, bc := range repo.BackupCycles{
// 	// maps["BackupCycles"]
// 	// if bc.Id
// 	// }

// 	return repo, nil
// }

func SaveRepo(repo *Repo, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)
	err = encoder.Encode(repo)
	if err != nil {
		return err
	}

	return nil
}

func LoadRepo(path string) (*Repo, error) {
	gob.Register(&BackupSet{})
	gob.Register(&BackupCycle{})
	gob.Register(&Repo{})

	var repo *Repo

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := gob.NewDecoder(file)
	err = decoder.Decode(&repo)
	if err != nil {
		return nil, err
	}

	fmt.Println("链表已从文件加载:", path)
	return repo, nil
}
