package main

import (
	"github.com/skyline93/mysql-xtrabackup/internal/mysql"
	"github.com/skyline93/mysql-xtrabackup/internal/repo"
	"github.com/spf13/cobra"
)

var cmdBackup = &cobra.Command{
	Use:   "backup",
	Short: "backup",
	Run: func(cmd *cobra.Command, args []string) {
		// loginPath := "local"
		// throttle := 400
		// tryCompress := true
		// binPath := "/usr/local/xtrabackup/bin"
		// dbHostName := "mysql"
		// dbUser := "root"
		// backupHostName := "backuper"
		// backupUser := "root"

		// dataPath := "/var/lib/mysql"

		config := &repo.Config{
			Identifer:   "MYTEST1",
			Version:     "8.0.23",
			LoginPath:   "local",
			DbHostName:  "mysql",
			DbUser:      "root",
			Throttle:    400,
			TryCompress: true,

			BinPath:        "/usr/local/xtrabackup/bin",
			DataPath:       "/var/lib/mysql",
			BackupUser:     "root",
			BackupHostName: "backuper",
		}

		r := repo.NewRepo("MYTEST1", config)
		r.Init("./")

		backuper := mysql.NewBackuper()

		err := backuper.Backup(r, "full")
		if err != nil {
			panic(err)
		}

		// fmt.Printf("%s", res)
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
