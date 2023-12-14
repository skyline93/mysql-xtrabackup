package main

import (
	"fmt"
	"os"

	"github.com/skyline93/mysql-xtrabackup/internal/mysql"
	"github.com/skyline93/mysql-xtrabackup/internal/repository"
	"github.com/spf13/cobra"
)

var cmdRestore = &cobra.Command{
	Use:   "restore",
	Short: "restore",
	Run: func(cmd *cobra.Command, args []string) {
		r, err := repository.Load(restoreOptions.RepoPath)
		if err != nil {
			fmt.Printf("load repo error: %s", err)
			os.Exit(1)
		}

		restorer := mysql.NewRestorer()

		err = restorer.Restore(r, restoreOptions.TargetPath, restoreOptions.BackupSetId)
		if err != nil {
			fmt.Printf("restore failed, err: %s", err.Error())
			os.Exit(1)
		}
	},
}

type RestoreOptions struct {
	BackupSetId string
	RepoPath    string
	TargetPath  string
}

var restoreOptions RestoreOptions

func init() {
	cmdRoot.AddCommand(cmdRestore)

	f := cmdRestore.Flags()
	f.StringVarP(&restoreOptions.BackupSetId, "backupset_id", "i", "", "backup set id")
	f.StringVarP(&restoreOptions.RepoPath, "repo_path", "p", "", "repo path")
	f.StringVarP(&restoreOptions.TargetPath, "target_path", "t", "", "target path")

	cmdRestore.MarkFlagRequired("backupset_id")
	cmdRestore.MarkFlagRequired("repo_path")
	cmdRestore.MarkFlagRequired("target_path")
}
