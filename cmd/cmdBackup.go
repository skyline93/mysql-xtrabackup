package main

import (
	"github.com/spf13/cobra"
)

var cmdBackup = &cobra.Command{
	Use:   "backup -p /data/backup/MYTEST1 -t full",
	Short: "backup",
	Run: func(cmd *cobra.Command, args []string) {
		// r, err := repository.Load(backupOptions.RepoPath)
		// if err != nil {
		// 	fmt.Printf("load repo error: %s", err)
		// 	os.Exit(1)
		// }

		// backuper := mysql.NewBackuper()

		// err = backuper.Backup(r, backupOptions.BackupType)
		// if err != nil {
		// 	panic(err)
		// }
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
