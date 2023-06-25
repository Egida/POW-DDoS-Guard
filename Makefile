PORT=9990
ADDR ?= pow-server:${PORT}

build-server:
	docker build -t pow-server \
		--build-arg addr=${ADDR} \
		--build-arg application=server \
		.

build-client:
	docker build -t pow-client \
		--build-arg addr=${ADDR} \
		--build-arg application=client \
		.

create-network:
	docker network create pow || exit 0

start-server: create-network
	docker run \
		-p ${PORT}:${PORT} \
		--rm -d \
		--cpus 1 \
        --memory 300M \
		--name pow-server \
		--network pow \
		pow-server

start-client: create-network
	docker run \
		--rm \
		--network pow \
		pow-client

clean:
	docker stop pow-server pow-client || exit 0
	docker rmi -f pow-client pow-server || exit 0
	docker network rm pow || exit 0

build: clean create-network build-server build-client

run: create-network build-server build-client start-server start-client clean