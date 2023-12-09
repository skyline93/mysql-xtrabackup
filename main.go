package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"os"

	"github.com/google/uuid"
)

type BackupSet struct {
	Id   string
	Path string
	Type string
	Prev *BackupSet
	Next *BackupSet
}

type BackupCycle struct {
	Id         string
	BackupSets []*BackupSet
	Prev       *BackupCycle
	Next       *BackupCycle
}

type Repo struct {
	id           string
	backupCycles []*BackupCycle
}

func NewBackupSet(path string, backupSetType string) *BackupSet {
	return &BackupSet{
		Id:   uuid.New().String(),
		Path: path,
		Type: backupSetType,
	}
}

func NewBackupCycle() *BackupCycle {
	return &BackupCycle{}
}

func (bc *BackupCycle) Insert(backupSet *BackupSet) {
	backupSet.Prev = bc.Tail()
	if bc.Tail() != nil {
		bc.Tail().Next = backupSet
	}
	bc.BackupSets = append(bc.BackupSets, backupSet)
}

func (bc *BackupCycle) Head() *BackupSet {
	if len(bc.BackupSets) > 0 {
		return bc.BackupSets[0]
	}
	return nil
}

func (bc *BackupCycle) Tail() *BackupSet {
	if len(bc.BackupSets) > 0 {
		return bc.BackupSets[len(bc.BackupSets)-1]
	}
	return nil
}

func NewRepo(Id string) *Repo {
	return &Repo{id: Id}
}

func (r *Repo) Insert(backupCycle *BackupCycle) {
	backupCycle.Prev = r.Tail()
	if r.Tail() != nil {
		r.Tail().Next = backupCycle
	}
	r.backupCycles = append(r.backupCycles, backupCycle)
}

func (r *Repo) Head() *BackupCycle {
	if len(r.backupCycles) > 0 {
		return r.backupCycles[0]
	}
	return nil
}

func (r *Repo) Tail() *BackupCycle {
	if len(r.backupCycles) > 0 {
		return r.backupCycles[len(r.backupCycles)-1]
	}
	return nil
}

func init() {
	gob.Register(&BackupSet{})
	gob.Register(&BackupCycle{})
	gob.Register(&Repo{})
}

// Custom serialization for BackupSet
func (bs *BackupSet) GobEncode() ([]byte, error) {
	var bsCopy BackupSet
	bsCopy.Id = bs.Id
	bsCopy.Path = bs.Path
	bsCopy.Type = bs.Type
	// Exclude Prev and Next to avoid circular references
	return gobEncode(bsCopy)
}

// Custom deserialization for BackupSet
func (bs *BackupSet) GobDecode(data []byte) error {
	var bsCopy BackupSet
	if err := gobDecode(data, &bsCopy); err != nil {
		return err
	}
	bs.Id = bsCopy.Id
	bs.Path = bsCopy.Path
	bs.Type = bsCopy.Type
	return nil
}

// Custom serialization for BackupCycle
func (bc *BackupCycle) GobEncode() ([]byte, error) {
	var bcCopy BackupCycle
	bcCopy.Id = bc.Id
	// Exclude Prev and Next to avoid circular references
	return gobEncode(bcCopy)
}

// Custom deserialization for BackupCycle
func (bc *BackupCycle) GobDecode(data []byte) error {
	var bcCopy BackupCycle
	if err := gobDecode(data, &bcCopy); err != nil {
		return err
	}
	bc.Id = bcCopy.Id
	return nil
}

// Custom serialization for Repo
func (r *Repo) GobEncode() ([]byte, error) {
	var rCopy Repo
	rCopy.id = r.id
	// Exclude Prev and Next to avoid circular references
	return gobEncode(rCopy)
}

// Custom deserialization for Repo
func (r *Repo) GobDecode(data []byte) error {
	var rCopy Repo
	if err := gobDecode(data, &rCopy); err != nil {
		return err
	}
	r.id = rCopy.id
	return nil
}

// gobEncode encodes the object to bytes
func gobEncode(obj interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(obj)
	return buf.Bytes(), err
}

// gobDecode decodes bytes to the object
func gobDecode(data []byte, obj interface{}) error {
	dec := gob.NewDecoder(bytes.NewReader(data))
	return dec.Decode(obj)
}

// IDToObject 映射 ID 到对象的映射
var IDToObject = make(map[string]interface{})

// 序列化Repo到文件
func serializeRepo(filename string, repo *Repo) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)

	// 清空映射关系
	IDToObject = make(map[string]interface{})

	// 序列化Repo
	err = encoder.Encode(repo)
	if err != nil {
		return err
	}

	fmt.Println("Repo 已序列化到文件:", filename)
	return nil
}

// 从文件反序列化Repo
func deserializeRepo(filename string) (*Repo, error) {
	var repo Repo

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := gob.NewDecoder(file)

	// 清空映射关系
	IDToObject = make(map[string]interface{})

	// 反序列化Repo
	err = decoder.Decode(&repo)
	if err != nil {
		return nil, err
	}

	// 根据ID恢复引用关系
	resolveReferences(&repo)

	fmt.Println("Repo 已从文件反序列化:", filename)
	return &repo, nil
}

// 根据ID恢复引用关系
func resolveReferences(repo *Repo) {
	// 根据ID恢复BackupCycles
	for i, cycle := range repo.backupCycles {
		repo.backupCycles[i] = restoreBackupCycle(cycle)
	}

	// 根据ID恢复BackupSets
	for _, cycle := range repo.backupCycles {
		for i, set := range cycle.BackupSets {
			cycle.BackupSets[i] = restoreBackupSet(set)
		}
	}
}

// 根据ID恢复BackupCycle对象
func restoreBackupCycle(cycle *BackupCycle) *BackupCycle {
	return restoreObject(cycle.Id, cycle).(*BackupCycle)
}

// 根据ID恢复BackupSet对象
func restoreBackupSet(set *BackupSet) *BackupSet {
	return restoreObject(set.Id, set).(*BackupSet)
}

// 根据ID恢复对象
func restoreObject(id string, obj interface{}) interface{} {
	IDToObject[id] = obj
	return obj
}

func main() {
	// 省略其他部分...
	bs1 := NewBackupSet("/backup/set/path1", "full")
	bs2 := NewBackupSet("/backup/set/path2", "incr")
	bs3 := NewBackupSet("/backup/set/path3", "incr")

	bc1 := NewBackupCycle()

	bc1.Insert(bs1)
	bc1.Insert(bs2)
	bc1.Insert(bs3)

	bs4 := NewBackupSet("/backup/set/path4", "full")
	bs5 := NewBackupSet("/backup/set/path5", "incr")
	bs6 := NewBackupSet("/backup/set/path6", "incr")

	bc2 := NewBackupCycle()

	bc2.Insert(bs4)
	bc2.Insert(bs5)
	bc2.Insert(bs6)

	bs7 := NewBackupSet("/backup/set/path4", "full")
	bs8 := NewBackupSet("/backup/set/path5", "incr")
	bs9 := NewBackupSet("/backup/set/path6", "incr")
	bc3 := NewBackupCycle()

	bc3.Insert(bs7)
	bc3.Insert(bs8)
	bc3.Insert(bs9)

	repo := NewRepo("24680")
	repo.Insert(bc1)
	repo.Insert(bc2)
	repo.Insert(bc3)

	// 序列化Repo到文件
	err := serializeRepo("repo.gob", repo)
	if err != nil {
		fmt.Println("序列化Repo时出错:", err)
		return
	}

	// 从文件反序列化Repo
	loadedRepo, err := deserializeRepo("repo.gob")
	if err != nil {
		fmt.Println("从文件反序列化Repo时出错:", err)
		return
	}

	// 打印反序列化后的Repo
	fmt.Printf("反序列化后的Repo: %+v\n", loadedRepo)
}
