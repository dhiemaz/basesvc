BINARY=basesvc
test: 
	go test -v -cover -covermode=atomic ./...

engine:
	go build -o ${BINARY} main.go

dependencies:
	@echo "> Installing the server dependencies ..."
	@go mod tidy -v
	@go install github.com/swaggo/swag/cmd/swag

unittest:
	go test -short  ./...

clean:
	if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi

lint-prepare:
	@echo "Installing golangci-lint" 
	curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh| sh -s latest

lint:
	./bin/golangci-lint run ./...

setup:
	@cp config/example/mysql.yml.example config/db/mysql.yml
	@cp config/example/rest.yml.example config/server/rest.yml
	@cp config/example/logger.yml.example config/logging/logger.yml

docs:
	@echo "Generating swagger"
	swag init -g infrastructure/router/router.go

.PHONY: clean install unittest lint-prepare lint docs engine test setup dependencies