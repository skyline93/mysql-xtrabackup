package main

import (
	"fmt"

	"github.com/skyline93/mysql-xtrabackup/internal/mysql"
	"github.com/spf13/cobra"
)

var cmdBackup = &cobra.Command{
	Use:   "backup",
	Short: "backup",
	Run: func(cmd *cobra.Command, args []string) {
		loginPath := "local"
		throttle := 400
		tryCompress := true
		binPath := "/usr/local/xtrabackup/bin"
		dbHostName := "mysql"
		dbUser := "root"
		backupHostName := "backuper"
		backupUser := "root"

		dataPath := "/var/lib/mysql"

		backuper := mysql.NewBackuper(
			loginPath,
			throttle,
			tryCompress,
			binPath,
			dbHostName,
			dbUser,
			backupHostName,
			backupUser,
		)

		res, err := backuper.Backup(dataPath, backupOptions.TargetPath, backupOptions.LogStart)
		if err != nil {
			panic(err)
		}

		fmt.Printf("%s", res)
	},
}

type BackupOptions struct {
	TargetPath string
	LogStart   string
}

var backupOptions BackupOptions

func init() {
	cmdRoot.AddCommand(cmdBackup)

	f := cmdBackup.Flags()
	f.StringVarP(&backupOptions.TargetPath, "target_path", "p", "/data/backup/f1", "target path")
	f.StringVarP(&backupOptions.LogStart, "log_start", "s", "", "server port")

	cmdBackup.MarkFlagRequired("target_path")
}
