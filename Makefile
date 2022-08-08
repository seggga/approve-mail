USERNAME := seggga
APP_NAME := mail-srv
VERSION := 1.0.1

run_container:
	docker run -ti docker.io/$(USERNAME)/$(APP_NAME):$(VERSION) sh

run_app:
	AUTH_PORT_3000_TCP_PORT=3000 AUTH_PORT_4000_TCP_PORT=4000 go run ./cmd/server/main.go -c ./configs/config.yaml


gen_proto:
	mkdir -p pkg/proto && \
	protoc  proto/*.proto --go-grpc_out=pkg --go_out=pkg

build_cpu_profile:
	go tool pprof -svg http://172.31.193.58:3000/debug/pprof/profile\?seconds\=5 > ./pprof/pprf.svg

build_memory_profile:
	go tool pprof http://172.31.193.58:3000/debug/pprof/heap

kafka/compose/up:
	docker-compose -f stack_kafka.yaml up -d

kafka/compose/down:
	docker-compose -f stack_kafka.yaml down