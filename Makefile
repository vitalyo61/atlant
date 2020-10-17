test.int: docker.start.test docker.stop
.PHONY: test.int

docker.start.components:
	@docker-compose up --build --abort-on-container-exit --remove-orphans mongo tests

docker.start.test:
	docker-compose up --build --abort-on-container-exit --remove-orphans mongo server1 server2 proxy csv tests_int

docker.stop:
	@docker-compose down

test: docker.start.components docker.stop

