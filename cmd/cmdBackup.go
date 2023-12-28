package main

import (
	"fmt"
	"os"

	"github.com/skyline93/mysql-xtrabackup/internal/mysql"
	"github.com/skyline93/mysql-xtrabackup/internal/repository"
	"github.com/spf13/cobra"
)

var cmdBackup = &cobra.Command{
	Use:   "backup -p /data/backup/MYTEST1 -t full",
	Short: "backup",
	Run: func(cmd *cobra.Command, args []string) {
		repo := repository.Repository2{}
		if err := repository.LoadRepository2(&repo, backupOptions.RepoPath); err != nil {
			fmt.Printf("load repo error: %s", err)
			os.Exit(1)
		}

		backuper := mysql.NewBackuper()
		if err := backuper.Backup(&repo, backupOptions.BackupType); err != nil {
			fmt.Printf("backup failed error: %s", err)
			os.Exit(1)
		}
	},
}

type BackupOptions struct {
	BackupType string
	RepoPath   string
}

var backupOptions BackupOptions

func init() {
	cmdRoot.AddCommand(cmdBackup)

	f := cmdBackup.Flags()
	f.StringVarP(&backupOptions.BackupType, "backup_type", "t", "full", "backup type")
	f.StringVarP(&backupOptions.RepoPath, "repo_path", "p", "", "repo path")

	cmdBackup.MarkFlagRequired("backup_type")
	cmdBackup.MarkFlagRequired("repo_path")
}
