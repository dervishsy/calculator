# Переменные
BINARY_NAME=calculator
ORCHESTRATOR_BINARY=./$(BINARY_NAME)-orchestrator.exe
AGENT_BINARY=./$(BINARY_NAME)-agent.exe

# Команды
all: clean build

# Очистка
clean:
	rm -f $(ORCHESTRATOR_BINARY) $(AGENT_BINARY)

# Сборка
build: build-o build-a

build-o:
	go build -o $(ORCHESTRATOR_BINARY) ./cmd/orchestrator/main.go

build-a:
	go build -o $(AGENT_BINARY) ./cmd/agent/main.go

# Запуск
run-o:
	./$(ORCHESTRATOR_BINARY)

run-a:
	./$(AGENT_BINARY)

# Запуск с перезагрузкой после изменения исходников (требуется https://github.com/cosmtrek/air)
air-o:
	air -c "./.air/air-o.toml" 

air-a:
	air -c "./.air/air-a.toml" 

# Docker
docker-build:
	docker build -t calculator-orchestrator -f ./build/orchestrator/Dockerfile .
	docker build -t calculator-agent -f ./build/agent/Dockerfile .

docker-run:
	docker run -d -p 8080:8080 --name calculator-orchestrator calculator-orchestrator
	docker run -d --name calculator-agent calculator-agent

docker-compose:
	docker-compose up

# Тесты
test:
	curl --location 'localhost:8080/api/v1/calculate' --header 'Content-Type: application/json' --data '{"id":"100" ,"expression": "2 + 2 * 2"}'
	curl --location 'localhost:8080/api/v1/calculate' --header 'Content-Type: application/json' --data '{"id":"101" ,"expression": "2 * 2 * 2"}'
	curl --location 'localhost:8080/api/v1/calculate' --header 'Content-Type: application/json' --data '{"id":"102" ,"expression": "2 / 2 * 2"}'
	curl --location 'localhost:8080/api/v1/calculate' --header 'Content-Type: application/json' --data '{"id":"103" ,"expression": "2 - 2 * 2"}'
	curl --location 'localhost:8080/api/v1/calculate' --header 'Content-Type: application/json' --data '{"id":"104" ,"expression": "2 * 2 + 2"}'

# Покрытие кода тестами 
cover:
	go test -v -coverpkg=./... -coverprofile=./.tmp/.cover.out  ./...
	go tool cover -html=./tmp/.cover.out

# Покрытие кода тестами в виде svg файла (требуется https://github.com/nikolaydubina/go-cover-treemap)
cover-svg:
	go test -coverprofile ./.tmp/.cover.out ./...
	go-cover-treemap -coverprofile ./.tmp/.cover.out > ./.tmp/.out.svg

# Помощь
help:
	@echo "Available commands:"
	@echo "  make all             Build all binaries"
	@echo "  make clean           Remove all binaries"
	@echo "  make build           Build all binaries"
	@echo "  make build-a         Build agent binary"
	@echo "  make build-o         Build orchestrator binary"
	@echo "  make run-o           Run orchestrator binary"
	@echo "  make run-a           Run agent binary"
	@echo "  make air-o           Run reloaded after change orchestrator binary"
	@echo "  make air-a           Run reloaded after change agent binary"
	@echo "  make docker-build    Build Docker images"
	@echo "  make docker-run      Run Docker containers"
	@echo "  make test            Start test"
	@echo "  make cover           Measure test coverage"
	@echo "  make cover-svg       Measure test coverage to svg (https://github.com/nikolaydubina/go-cover-treemap)"
	@echo "  make help            Show this help"
