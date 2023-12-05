package index

type BackupSet struct {
	Path     string
	Type     string
	Previous *BackupSet
	Next     *BackupSet
}

func NewBackupSet(path string, backupSetType string) *BackupSet {
	return &BackupSet{
		Path: path,
		Type: backupSetType,
	}
}

type BackupCycle struct {
	FullBackupSet *BackupSet
	Previous      *BackupCycle
	Next          *BackupCycle
}

func NewBackupCycle() *BackupCycle {
	return &BackupCycle{}
}
