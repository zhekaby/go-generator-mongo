all: run_docker generate test build stop_docker

run_docker:
	docker-compose down --remove-orphans && docker-compose up -d

wait_docker:
	echo 'waiting replica is up...'
	#while ! docker ps | grep mongowrapper | grep '(healthy)'; do sleep 1; done
	sleep 8

stop_docker:
	docker-compose down --remove-orphans

generate:
	go run . --cs_var=MONGODB_CONNECTION_STRING --cs='mongodb://db1:31001,db2:31002/ipo?replicaSet=rs&readPreference=primaryPreferred' tests/

build:
	CGO_ENABLED=0 go build -a -ldflags="-w -s -X main.VERSION=`git rev-parse HEAD | cut -c1-8`" -o ~/go/bin/mongowrapper

test:
	go test -v -coverprofile cover.out ./tests && \
		go tool cover -func cover.out

