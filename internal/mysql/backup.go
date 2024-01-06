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

func (b *Backuper) Backup(repo *repository.Repository2, backupType string) (err error) {
	bs := repository.NewBackupSet2(backupType)
	targetPath, err := filepath.Abs(filepath.Join(repo.DataPath(), bs.Id))
	if err != nil {
		return err
	}

	if _, err = os.Stat(targetPath); os.IsNotExist(err) {
		err = os.MkdirAll(targetPath, 0755)
		if err != nil {
			return err
		}
		log.Printf("create path: %s", targetPath)
	}

	defer func() {
		if err != nil {
			log.Printf("backup failed, err: %s", err)
			os.RemoveAll(targetPath)
		}
	}()

	backupArgs := []string{
		filepath.Join(repo.Config.BinPath, "xtrabackup"),
		"--backup",
		fmt.Sprintf("--throttle=%d", repo.Config.Throttle),
		fmt.Sprintf("--login-path=%s", repo.Config.LoginPath),
		fmt.Sprintf("--datadir=%s", repo.Config.DataPath),
		"--stream=xbstream",
	}

	if repo.Config.TryCompress {
		backupArgs = append(backupArgs, "--compress")
	}

	if backupType == repository.TypeBackupSetIncr {
		lastBackupSet, err := repo.GetLastBackupSet()
		if err != nil {
			return err
		}

		backupArgs = append(backupArgs, fmt.Sprintf("--incremental-lsn=%s", lastBackupSet.ToLSN))
	}

	streamArgs := []string{
		"ssh", fmt.Sprintf("%s@%s", repo.Config.BackupUser, repo.Config.BackupHostName),
		filepath.Join(repo.Config.BinPath, "xbstream"), "-x", "-C", targetPath,
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
	if err = cmd.Run(); err != nil {
		return err
	}

	content, err := os.ReadFile(filepath.Join(targetPath, "xtrabackup_checkpoints"))
	if err != nil {
		return err
	}

	checkpoints, err := b.parseCheckpoints(string(content))
	if err != nil {
		return err
	}

	size, err := b.getBackupSize(targetPath)
	if err != nil {
		return err
	}

	bs.Path = targetPath
	bs.FromLSN = checkpoints["from_lsn"]
	bs.ToLSN = checkpoints["to_lsn"]
	bs.Size = int64(size)

	if err = repo.AddBackupSet(bs); err != nil {
		return err
	}

	log.Printf("backup completed.\nbackupset: %s\npath: %s\nfrom_lsn: %s\nto_lsn: %s\nsize: %dbyte", bs.Id, bs.Path, bs.FromLSN, bs.ToLSN, bs.Size)
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
