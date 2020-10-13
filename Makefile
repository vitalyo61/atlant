docker.start.components:
	docker-compose up -d --remove-orphans mongo

docker.stop:
	docker-compose down

test.unit:
	go test -v ./...

test: docker.start.components test.unit docker.stop
