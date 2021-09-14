build:
	go build -o=go-acess-control && ./go-acess-control start

docker-build:
	docker build \
		--rm \
		-f ./Dockerfile \
		-t serbanblebea/go-acess-control:0.1 \
		.

docker-run:
	docker run \
		-p 8081:8081 \
		--env-file ./.env \
		-d \
		--name go-acess-control \
		-v ${HOME}/Projects/Go/go-access-control/data:/app/data \
		serbanblebea/go-acess-control:0.1

docker: docker-build docker-run

docker-down:
	docker stop go-acess-control && docker rm go-acess-control

unit-test:
	go test -v .