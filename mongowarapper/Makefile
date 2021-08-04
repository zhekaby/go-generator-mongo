all: run_docker generate wait_docker test stop_docker

run_docker:
	docker-compose down --remove-orphans && docker-compose up -d

wait_docker:
	echo 'waiting replica is up...'
	while ! docker ps | grep mongowarapper | grep '(healthy)'; do sleep 1; done
	sleep 5

stop_docker:
	docker-compose down --remove-orphans

generate:
	go run . --cs_var=MONGODB_CONNECTION_STRING --cs='mongodb://db1:33001,db2:33002/ipo?replicaSet=mongowarapper-tests&readPreference=primaryPreferred' tests/

test:
	go test -v -coverprofile cover.out ./tests && \
		go tool cover -func cover.out

