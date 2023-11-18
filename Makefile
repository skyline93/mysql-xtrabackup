build: build-mysql build-backuper

build-mysql:
	docker build -t glf9832/mysql-xtrabackup:8.0.28 ./mysql

build-backuper:
	docker build -t glf9832/xtrabackup:8.0.28 ./backuper

up:
	mkdir -p .testdata/mysqldata .testdata/backupdata
	docker-compose up -d

down:
	docker-compose down

exec-backuper:
	docker-compose exec -it backuper bash

exec-mysql:
	docker-compose exec -it mysql bash

clean: down
	rm -rf .testdata
