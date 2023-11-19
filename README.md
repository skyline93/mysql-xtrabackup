```bash
/usr/local/xtrabackup/bin/xtrabackup --backup --compress --stream=xbstream --target-dir=./ | ssh root@backuper "/usr/local/xtrabackup/bin/xbstream -x"
```

```bash
mysql_config_editor set --login-path=local --host=mysql --user=root --password
```

```bash
mysql --login-path=local
```

备份
```bash
ssh root@mysql '/usr/local/xtrabackup/bin/xtrabackup --backup --throttle=400 --login-path=local --datadir='/var/lib/mysql' --stream=xbstream --compress | ssh root@backuper "/usr/local/xtrabackup/bin/xbstream -x -C /data/backup/f1"'
```

```bash
ssh root@mysql '/usr/local/xtrabackup/bin/xtrabackup --backup --throttle=400 --login-path=local --datadir='/var/lib/mysql' --stream=xbstream --compress --incremental-lsn=28798178 | ssh root@backuper "/usr/local/xtrabackup/bin/xbstream -x -C /data/i2-1"'
```

```bash
ssh root@mysql '/usr/local/xtrabackup/bin/xtrabackup --backup --throttle=400 --login-path=local --datadir='/var/lib/mysql' --stream=xbstream --compress --incremental-lsn=28798178 | ssh root@backuper "/usr/local/xtrabackup/bin/xbstream -x -C /data/i3-1"'
```

```bash
ssh root@mysql '/usr/local/xtrabackup/bin/xtrabackup --backup --throttle=400 --login-path=local --datadir='/var/lib/mysql' --stream=xbstream --compress --incremental-lsn=28798178 | ssh root@backuper "/usr/local/xtrabackup/bin/xbstream -x -C /data/i4-1"'
```

```bash
ssh root@mysql '/usr/local/xtrabackup/bin/xtrabackup --backup --throttle=400 --login-path=local --datadir='/var/lib/mysql' --stream=xbstream --compress --incremental-lsn=29249757 | ssh root@backuper "/usr/local/xtrabackup/bin/xbstream -x -C /data/i5-1"'
```

恢复
```bash
/usr/local/xtrabackup/bin/xtrabackup --decompress --remove-original --target-dir=/data/restore/f1
/usr/local/xtrabackup/bin/xtrabackup --decompress --remove-original --target-dir=/data/restore/i2-1
```

```bash
/usr/local/xtrabackup/bin/xtrabackup --prepare --apply-log-only --target-dir=/data/restore/f1 
/usr/local/xtrabackup/bin/xtrabackup --prepare --apply-log-only --target-dir=/data/restore/f1 --incremental-dir=/data/restore/i2-1
```

```vim
cat << EOF > my.cnf
[mysqld]
basedir=/usr/bin/mysql
datadir=/data/restore/f1
socket=/data/restore/mysql.sock
pid-file=/data/restore/mysql.pid
port=3308
EOF
```

```bash
sudo chown -R mysql:mysql /data/restore
```

```bash
sudo -u mysql mysqld_safe --defaults-file=/data/restore/my.cnf
```
