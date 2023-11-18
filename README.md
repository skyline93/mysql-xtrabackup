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
ssh root@mysql '/usr/local/xtrabackup/bin/xtrabackup --backup --throttle=400 --login-path=local --datadir='/var/lib/mysql' --stream=xbstream --compress | ssh root@backuper "/usr/local/xtrabackup/bin/xbstream -x -C /data/f1"'
```

```bash
ssh root@mysql '/usr/local/xtrabackup/bin/xtrabackup --backup --throttle=400 --login-path=local --datadir='/var/lib/mysql' --stream=xbstream --compress --incremental-lsn=28648031 | ssh root@backuper "/usr/local/xtrabackup/bin/xbstream -x -C /data/i2-1"'
```

```bash
ssh root@mysql '/usr/local/xtrabackup/bin/xtrabackup --backup --throttle=400 --login-path=local --datadir='/var/lib/mysql' --stream=xbstream --compress --incremental-lsn=28648031 | ssh root@backuper "/usr/local/xtrabackup/bin/xbstream -x -C /data/i3-1"'
```

```bash
ssh root@mysql '/usr/local/xtrabackup/bin/xtrabackup --backup --throttle=400 --login-path=local --datadir='/var/lib/mysql' --stream=xbstream --compress --incremental-lsn=28648031 | ssh root@backuper "/usr/local/xtrabackup/bin/xbstream -x -C /data/i4-1"'
```

恢复
```bash
/usr/local/xtrabackup/bin/xtrabackup --decompress --remove-original --target-dir=/target/path
```

```bash
/usr/local/xtrabackup/bin/xtrabackup --prepare --apply-log-only --target-dir=/target/path 
/usr/local/xtrabackup/bin/xtrabackup --prepare --apply-log-only --target-dir=/target/path --incremental-dir=/incr/path
```

```vim
[mysqld]
basedir=
datadir=
socket=
pid-file=
port=
```

```bash
sudo chown -R mysql:mysql /target/path
```

```bash
sudo -u mysql mysqld --defaults-file=my.cnf
```
