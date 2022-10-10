run-gym: 
	cd ./services/gym/app && wire
	cd ./services/gym/app/ && go run .

run-base:
	docker-compose -f deployments/docker/base-docker-compose.yml up

migrateup:
	migrate -path deployments/migration -database "postgres://postgres:postgres@localhost:5432/gym?sslmode=disable" up

migratedown:
	migrate -path deployments/migration -database "postgres://postgres:postgres@localhost:5432/gym?sslmode=disable" down

critic:
	go get github.com/go-critic/go-critic/cmd/gocritic
	go install github.com/go-critic/go-critic/cmd/gocritic
	gocritic check ./...

migrate:
	bash deployments/scripts/sample.sh

