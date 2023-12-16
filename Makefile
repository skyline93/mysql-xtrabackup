build: build-mysql build-backuper

build-mysql:
	docker build -t glf9832/mysql-xtrabackup:8.0.28-3 ./mysql

build-backuper:
	docker build -t glf9832/xtrabackup:8.0.28-4 ./backuper

up:
	mkdir -p .testdata/mysqldata .testdata/backupdata .testdata/restoredata
	docker-compose up -d

down:
	docker-compose down

exec-backuper:
	docker-compose exec -it backuper bash

exec-restorer:
	docker-compose exec -it restorer bash

exec-mysql:
	docker-compose exec -it mysql bash

clean: down
	sudo rm -rf .testdata
	docker volume rm mysql-xtrabackup_backup_data mysql-xtrabackup_restore_data mysql-xtrabackup_mysql_data
