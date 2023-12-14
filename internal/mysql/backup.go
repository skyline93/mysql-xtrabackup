package mysql

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/skyline93/mysql-xtrabackup/internal/repository"
)

type Backuper struct {
}

func NewBackuper() *Backuper {
	return &Backuper{}
}

func (b *Backuper) Backup(repo *repository.Repository, backupType string) error {
	var lastBackupSet *repository.BackupSet

	// TODO 封装到repo中
	if repo.Tail() != nil {
		lastBackupSet = repo.Tail().Tail()
	}

	bs := repository.NewBackupSet(backupType)
	repo.AddBackupSet(bs)

	if _, err := os.Stat(bs.Path); os.IsNotExist(err) {
		err := os.MkdirAll(bs.Path, 0755)
		if err != nil {
			return err
		}
		log.Printf("create path: %s", bs.Path)
	}

	backupArgs := []string{
		filepath.Join(repo.Config.BinPath, "xtrabackup"), "--backup", fmt.Sprintf("--throttle=%d", repo.Config.Throttle), fmt.Sprintf("--login-path=%s", repo.Config.LoginPath), fmt.Sprintf("--datadir=%s", repo.Config.DataPath), "--stream=xbstream",
	}

	if repo.Config.TryCompress {
		backupArgs = append(backupArgs, "--compress")
	}

	if backupType == repository.TypeBackupSetIncr {
		backupArgs = append(backupArgs, fmt.Sprintf("--incremental-lsn=%s", lastBackupSet.ToLSN))
	}

	streamArgs := []string{
		"ssh", fmt.Sprintf("%s@%s", repo.Config.BackupUser, repo.Config.BackupHostName),
		filepath.Join(repo.Config.BinPath, "xbstream"), "-x", "-C", bs.Path,
	}

	args := append(append(backupArgs, []string{"|"}...), streamArgs...)

	xtraLogPath, err := filepath.Abs(fmt.Sprintf("logs/xtrabackup-%s.log", time.Now().Format("20060102150405")))
	if err != nil {
		return err
	}

	log.Printf("log path: %s", xtraLogPath)
	logFile, err := os.Create(xtraLogPath)
	if err != nil {
		return err
	}
	defer logFile.Close()

	cmd := exec.Command("ssh", fmt.Sprintf("%s@%s", repo.Config.DbUser, repo.Config.DbHostName), strings.Join(args, " "))
	cmd.Stdout = logFile
	cmd.Stderr = logFile

	log.Printf("cmd: %s", cmd.String())
	if err := cmd.Run(); err != nil {
		return err
	}

	content, err := os.ReadFile(filepath.Join(bs.Path, "xtrabackup_checkpoints"))
	if err != nil {
		return err
	}

	checkpoints, err := b.parseCheckpoints(string(content))
	if err != nil {
		return err
	}

	size, err := b.getBackupSize(bs.Path)
	if err != nil {
		return err
	}

	bs.FromLSN = checkpoints["from_lsn"]
	bs.ToLSN = checkpoints["to_lsn"]
	bs.Size = int64(size)

	if err = repo.Commit(); err != nil {
		return err
	}

	log.Printf("backup completed. \n\nbackupset: %s\npath: %s\nfrom_lsn: %s\nto_lsn: %s", bs.Id, bs.Path, bs.FromLSN, bs.ToLSN)
	return nil
}

func (b *Backuper) parseCheckpoints(content string) (map[string]string, error) {
	checkpointsMap := make(map[string]string)
	scanner := bufio.NewScanner(strings.NewReader(content))

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			checkpointsMap[key] = value
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return checkpointsMap, nil
}

func (b *Backuper) getBackupSize(targetPath string) (uint64, error) {
	cmd := exec.Command("du", "-sb", targetPath)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return 0, err
	}

	// 解析 du 输出获取备份数据量
	fields := strings.Fields(string(output))
	size, err := strconv.ParseUint(fields[0], 10, 64)
	if err != nil {
		return 0, err
	}
	return size, nil
}
