package mysql

import (
	"log"
	"os/exec"

	"github.com/skyline93/mysql-xtrabackup/internal/repository"
)

type Restorer struct {
}

func NewRestorer() *Restorer {
	return &Restorer{}
}

func (r *Restorer) Restore(repo *repository.Repository, targetPath string, backupSetId string) error {
	bs, err := repo.FindBackupSet(backupSetId)
	if err != nil {
		return err
	}

	bc, err := repo.FindBackupCycle(backupSetId)

	fullSet := bc.Head()


	var backupsets []repository.BackupSet

	backupset := *bs
	for {
		backupsets = append(backupsets, backupset)
		if backupset.Prev == nil {
			break
		}
		backupset = *backupset.Prev
	}

	fullSet := backupsets[len(backupsets)-1]
	cmd := exec.Command("cp", "-r", backupset.Path, targetPath)
	log.Printf("cmd: %s", cmd.String())
	if err := cmd.Run(); err != nil {
		return err
	}



	// for i := len(backupsets) - 1; i == 0; i-- {
	// 	bs := backupsets[i]

	// 	cmd := exec.Command("cp", "-r", backupset.Path, targetPath)
	// 	log.Printf("cmd: %s", cmd.String())
	// 	if err := cmd.Run(); err != nil {
	// 		return err
	// 	}
	// }

	return nil
}
