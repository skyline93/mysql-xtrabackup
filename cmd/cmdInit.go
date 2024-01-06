package main

import (
	"fmt"
	"os"

	"github.com/skyline93/mysql-xtrabackup/internal/repository"
	"github.com/spf13/cobra"
)

var cmdInit = &cobra.Command{
	Use:   "init -n MYTEST1 -p /data/backup",
	Short: "init",
	Run: func(cmd *cobra.Command, args []string) {
		config := &repository.Config{
			Identifer:   "MYTEST1",
			Version:     "8.0.23",
			LoginPath:   "local",
			DbHostName:  "172.18.53.110",
			DbUser:      "mysql",
			Throttle:    400,
			TryCompress: true,

			BinPath:        "/usr/local/xtrabackup/8.0.28/bin",
			DataPath:       "/var/lib/mysql",
			BackupUser:     "skyline93",
			BackupHostName: "172.18.53.110",
		}

		r := repository.NewRepository(initOptions.RepoName, config)
		if err := r.Init(initOptions.Root); err != nil {
			fmt.Printf("err: %s\n", err)
			os.Exit(1)
		}
	},
}

type InitOptions struct {
	Root     string
	RepoName string
}

var initOptions InitOptions

func init() {
	cmdRoot.AddCommand(cmdInit)

	f := cmdInit.Flags()
	f.StringVarP(&initOptions.Root, "path", "p", "/data/repo/root", "repo root path")
	f.StringVarP(&initOptions.RepoName, "repo_name", "n", "", "repo name")

	cmdInit.MarkFlagRequired("path")
	cmdInit.MarkFlagRequired("repo_name")
}
