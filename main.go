package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

type BackupResult struct {
	LogStart string
	LogEnd   string
	Size     uint64
}

func (r *BackupResult) String() string {
	return fmt.Sprintf("\nlogStart: %s\nlogEnd: %s\nSize: %dbytes\n", r.LogStart, r.LogEnd, r.Size)
}

type Backuper struct {
	loginPath      string
	throttle       int
	tryCompress    bool
	binPath        string
	dbHostName     string
	dbUser         string
	backupHostName string
	backupUser     string
}

func NewBackuper(loginPath string, throttle int, tryCompress bool, binPath string, dbHostName, dbUser string, backupHostName, backupUser string) *Backuper {
	return &Backuper{
		loginPath:      loginPath,
		throttle:       throttle,
		tryCompress:    tryCompress,
		binPath:        binPath,
		dbHostName:     dbHostName,
		dbUser:         dbUser,
		backupHostName: backupHostName,
		backupUser:     backupUser,
	}
}

func (b *Backuper) Backup(dataPath string, targetPath string, logStart string) (*BackupResult, error) {
	if _, err := os.Stat(targetPath); os.IsNotExist(err) {
		err := os.MkdirAll(targetPath, 0755)
		if err != nil {
			return nil, err
		}
		log.Printf("create path: %s", targetPath)
	}

	backupArgs := []string{
		filepath.Join(b.binPath, "xtrabackup"), "--backup", fmt.Sprintf("--throttle=%d", b.throttle), fmt.Sprintf("--login-path=%s", b.loginPath), fmt.Sprintf("--datadir=%s", dataPath), "--stream=xbstream",
	}

	if b.tryCompress {
		backupArgs = append(backupArgs, "--compress")
	}

	if logStart != "" {
		backupArgs = append(backupArgs, fmt.Sprintf("--incremental-lsn=%s", logStart))
	}

	streamArgs := []string{
		"ssh", fmt.Sprintf("%s@%s", b.backupUser, b.backupHostName),
		filepath.Join(b.binPath, "xbstream"), "-x", "-C", targetPath,
	}

	args := append(append(backupArgs, []string{"|"}...), streamArgs...)

	err := b.runCmd("ssh", fmt.Sprintf("%s@%s", b.dbUser, b.dbHostName), strings.Join(args, " "))
	if err != nil {
		return nil, err
	}

	content, err := os.ReadFile(filepath.Join(targetPath, "xtrabackup_checkpoints"))
	if err != nil {
		return nil, err
	}

	checkpoints, err := b.parseCheckpoints(string(content))
	if err != nil {
		return nil, err
	}

	size, err := b.getBackupSize(targetPath)
	if err != nil {
		return nil, err
	}

	return &BackupResult{
		Size:     size,
		LogStart: checkpoints["from_lsn"],
		LogEnd:   checkpoints["to_lsn"],
	}, nil
}

func (b *Backuper) runCmd(command string, args ...string) error {
	cmd := exec.Command(command, args...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		log.Printf("run command failed, err: %s", err)
		return err
	}

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

func main() {
	loginPath := "local"
	throttle := 400
	tryCompress := true
	binPath := "/usr/local/xtrabackup/bin"
	dbHostName := "mysql"
	dbUser := "root"
	backupHostName := "backuper"
	backupUser := "root"

	dataPath := "/var/lib/mysql"
	targetPath := "/data/f1"
	logStart := ""

	backuper := NewBackuper(
		loginPath,
		throttle,
		tryCompress,
		binPath,
		dbHostName,
		dbUser,
		backupHostName,
		backupUser,
	)

	res, err := backuper.Backup(dataPath, targetPath, logStart)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s", res)
}
