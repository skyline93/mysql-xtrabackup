```bash
/usr/local/xtrabackup/bin/xtrabackup --backup --compress --stream=xbstream --target-dir=./ | ssh root@backuper "/usr/local/xtrabackup/bin/xbstream -x"
```
