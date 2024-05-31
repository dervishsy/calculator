# Distributed Calculator for Arithmetic Expressions

This is a distributed calculator for arithmetic expressions, written in Go.

It consists of an orchestrator and multiple agents that can perform arithmetic operations in parallel.

## Features

- The orchestrator provides a RESTful API for submitting arithmetic expressions for calculation
- Agents request tasks from the orchestrator, perform calculations, and send the results back
- Parallel execution of arithmetic operations with configurable computing power
- Configurable operation times to simulate long-running computations
- Web interface for entering expressions and viewing results

## Requirements

- Go 1.22 or later

## Installation

1. Clone the repository:
```

git clone https://github.com/dervishsy/calculator.git

```

2. Change to the project directory:

```

cd calculator

```

3. Build the project:

```

go build -o ./orchestrator.exe ./cmd/orchestrator/main.go

go build -o ./agent.exe ./cmd/agent/main.go

```

or

```

make build

```

## Configuration

The application can be configured using YAML files. The default configuration files are:

- `configs/agent.yml`

- `configs/orchestrator.yml`

You can also configure the application using environment variables. If set, they will be used instead of the values from the configuration files.

### Agent Configuration

- `computingPower`: The number of worker goroutines for simulating parallel computations
- `orchestratorURL`: The URL of the orchestrator

or using the following environment variables:

- `COMPUTING_POWER`: The number of worker goroutines for simulating parallel computations
- `ORCHESTRATOR_URL`: The URL of the orchestrator

### Orchestrator Configuration

- `server.port`: The port number for the orchestrator's HTTP server
- `timeAdditionMS`: The simulated time (in milliseconds) for addition operations
- `timeSubtractionMS`: The simulated time (in milliseconds) for subtraction operations
- `timeMultiplicationMS`: The simulated time (in milliseconds) for multiplication operations
- `timeDivisionMS`: The simulated time (in milliseconds) for division operations

or using the following environment variables:

- `SERVER_PORT`: The port number for the orchestrator's HTTP server
- `TIME_ADDITION_MS`: The simulated time (in milliseconds) for addition operations
- `TIME_SUBTRACTION_MS`: The simulated time (in milliseconds) for subtraction operations
- `TIME_MULTIPLICATIONS_MS`: The simulated time (in milliseconds) for multiplication operations
- `TIME_DIVISIONS_MS`: The simulated time (in milliseconds) for division operations

## Usage

1. Start the orchestrator:

```

go run cmd/orchestrator/main.go

```

2. Start one or more agents:

```

go run cmd/agent/main.go

```

or

```

make docker-compose

```

3. Send an arithmetic expression using the orchestrator's API:

```

curl --location 'http://localhost:8080/api/v1/calculate' --header 'Content-Type: application/json' --data '{"id":"100" ,"expression": "2 + 2 * 2"}'

```

4. Check the status of an expression:

```

curl --location 'http://localhost:8080/api/v1/expressions'

```

```

curl --location 'http://localhost:8080/api/v1/expressions/:id'

```

5. The web interface is available at `http://localhost:8080`.