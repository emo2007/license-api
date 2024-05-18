PROJECT_DIR = $(CURDIR)
PROJECT_BIN = ${PROJECT_DIR}/bin
TOOLS_BIN = ${PROJECT_BIN}/tools

bin.build:
	mkdir -p ${PROJECT_DIR}/build
	rm -f ${PROJECT_DIR}/build/license-api
	go build -ldflags="-s -w" -o ${PROJECT_DIR}/build/license-api ${PROJECT_DIR}/cmd/main.go

d.build:
	sudo docker buildx build . -t license-api:latest

d.net:
	sudo docker network create --driver bridge --subnet=192.168.2.0/24 --attachable blockd-net

d.drop-net:
	sudo docker network rm blockd-net

up: 
	sudo docker compose --profile license up -d

.PHONY: run.local
run.local: bin.build
	${PROJECT_DIR}/build/license-api \
		-log-level=debug \
		-log-local=true \
		-log-add-source=true \
		-rest-address=localhost:8081 \
		-db-host=localhost:8433 \
		-db-database=blockd \
		-db-user=blockd \
		-db-secret=blockd \
		-db-enable-tls=false \
		-jwt-secret=local_jwt_secret \ 
		-cache-host=localhost:6379 

.PHONY: run.debug
run.debug: bin.build
	${PROJECT_DIR}/build/license-api \
		-log-level=debug \
		-log-local=false \
		-log-add-source=true \
		-rest-address=localhost:8081 \
		-db-host=localhost:8433 \
		-db-database=blockd \
		-db-user=blockd \
		-db-secret=blockd \
		-db-enable-tls=false 


start.d:
	sudo systemctl start docker