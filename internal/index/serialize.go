package index

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Node struct {
	Id   string `json:"id"`
	Prev string `json:"prev"`
	Next string `json:"next"`
	Path string `json:"path"`
	Type string `json:"type"`
}

type LinkedList struct {
	Id    string           `json:"id"`
	Prev  string           `json:"prev"`
	Next  string           `json:"next"`
	Nodes map[string]*Node `json:"nodes"`
}

type RepoMarshal struct {
	Id    string                 `json:"id"`
	Start string                 `json:"start"`
	Index map[string]*LinkedList `json:"index"`
}

func (bs *BackupSet) convertNode() *Node {
	prevId := ""
	nextId := ""

	if bs.Prev != nil {
		prevId = bs.Prev.Id
	}

	if bs.Next != nil {
		nextId = bs.Next.Id
	}

	return &Node{
		Id:   bs.Id,
		Prev: prevId,
		Next: nextId,
		Path: bs.Path,
		Type: bs.Type,
	}
}

func (bs *BackupSet) MarshalJSON() ([]byte, error) {
	return json.Marshal(bs.convertNode())
}

func (bc *BackupCycle) convertLinkedList() *LinkedList {
	prevId := ""
	nextId := ""

	if bc.Prev != nil {
		prevId = bc.Prev.Id
	}

	if bc.Next != nil {
		nextId = bc.Next.Id
	}

	nodes := map[string]*Node{}

	for _, bs := range bc.BackupSets {
		nodes[bs.Id] = bs.convertNode()
	}

	return &LinkedList{
		Id:    bc.Id,
		Prev:  prevId,
		Next:  nextId,
		Nodes: nodes,
	}
}

func (bc *BackupCycle) MarshalJSON() ([]byte, error) {
	return json.Marshal(bc.convertLinkedList())
}

func (r *Repo) MarshalJSON() ([]byte, error) {
	var bcs = make(map[string]*LinkedList)

	for _, bc := range r.BackupCycles {
		bcs[bc.Id] = bc.convertLinkedList()
	}

	repoMarshal := RepoMarshal{Id: r.Id, Start: r.BackupCycles[0].Id, Index: bcs}
	return json.Marshal(repoMarshal)
}

func (r *Repo) UnmarshalJSON(data []byte) error {
	var rm = &RepoMarshal{}
	err := json.Unmarshal(data, rm)
	if err != nil {
		return err
	}

	r.Id = rm.Id

	for {
		l := rm.Index[rm.Start]
		bc := &BackupCycle{Id: l.Id, BackupSets: make([]*BackupSet, 0)}
		for {
			n := l.Nodes[l.Id]
			bs := &BackupSet{Id: n.Id, Path: n.Path, Type: n.Type}
			bc.Insert(bs)
			if n.Next == "" {
				break
			}
		}

		r.BackupCycles = append(r.BackupCycles, bc)

		if l.Next == "" {
			break
		}
	}

	return nil
}

func SaveToFile(repo *Repo, path string) error {
	p, err := filepath.Abs(path)
	if err != nil {
		return err
	}

	filePath := filepath.Join(p, fmt.Sprintf("%s.json", repo.Id))

	v, err := json.Marshal(repo)
	if err != nil {
		return err
	}

	err = os.WriteFile(filePath, v, 0664)
	if err != nil {
		return err
	}

	return nil
}

func LoadFromFile(repo *Repo, path string) error {
	p, err := filepath.Abs(path)
	if err != nil {
		return err
	}

	data, err := os.ReadFile(p)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, repo)
}
