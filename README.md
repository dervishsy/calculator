# Распределенный калькулятор арифметических выражений

Это распределенный калькулятор арифметических выражений, написанный на Go.

Он состоит из оркестратора и нескольких агентов, которые могут выполнять арифметические операции параллельно.

## Возможности

- Оркестратор предоставляет RESTful API для отправки арифметических выражений для вычисления
- Агенты запрашивают задачи у оркестратора, выполняют вычисления и отправляют результаты обратно
- Параллельное выполнение арифметических операций с настраиваемой вычислительной мощностью
- Настраиваемое время операций для имитации долгосрочных вычислений
- Веб-интерфейс для ввода выражений и просмотра результатов

## Требования

- Go 1.22 и выше

## Установка

1. Клонируйте репозиторий:

```
git clone https://github.com/dervishsy/calculator.git
```

2. Перейдите в директорию проекта:

```
cd calculator
```

3. Соберите проект:

```
go build -o ./orchestrator.exe ./cmd/orchestrator/main.go
go build -o ./agent.exe ./cmd/agent/main.go
```

или

```
make build
```

## Конфигурация

Приложение может быть настроено с помощью YAML-файлов. Файлы конфигурации по умолчанию:
 - `configs/agent.yml`
 - `configs/orchestrator.yml`

Так же можно настроить приложение с помощью переменных окружения. Если они заданы, то они будут использованы вместо значений из файлов конфигурации.


### Конфигурация агента

- `computingPower`: Количество воркер-горутин для имитации параллельных вычислений
- `orchestratorURL`: URL-адрес оркестратора

или с помощью следующих переменных окружения:

- `COMPUTING_POWER` - Количество воркер-горутин для имитации параллельных вычислений
- `ORCHESTRATOR_URL` - URL-адрес оркестратора

### Конфигурация оркестратора

- `server.port`: Номер порта для HTTP-сервера оркестратора
- `timeAdditionMS`: Симулируемое время (в миллисекундах) для операций сложения
- `timeSubtractionMS`: Симулируемое время (в миллисекундах) для операций вычитания
- `timeMultiplicationMS`: Симулируемое время (в миллисекундах) для операций умножения
- `timeDivisionMS`: Симулируемое время (в миллисекундах) для операций деления

или с помощью следующих переменных окружения:
- `SERVER_PORT` : Номер порта для HTTP-сервера оркестратора
- `TIME_ADDITION_MS` : Симулируемое время (в миллисекундах) для операций сложения
- `TIME_SUBTRACTION_MS`: Симулируемое время (в миллисекундах) для операций вычитания
- `TIME_MULTIPLICATIONS_MS`: Симулируемое время (в миллисекундах) для операций умножения
- `TIME_DIVISIONS_MS`: Симулируемое время (в миллисекундах) для операций деления


## Использование

1. Запустите оркестратор:

```
go run cmd/orchestrator/main.go
```

2. Запустите одного или несколько агентов:

```
go run cmd/agent/main.go
```

или

```
make docker-compose
```

3. Отправить арифметическое выражение с помощью API оркестратора можно так:

```
curl --location 'http://localhost:8080/api/v1/calculate' --header 'Content-Type: application/json' --data '{"id":"100" ,"expression": "2 + 2 * 2"}'
```

4. Проверить статус выражения можно так:

```
curl --location 'http://localhost:8080/api/v1/expressions'
```

```
curl --location 'http://localhost:8080/api/v1/expressions/:id'
```

5. Веб-интерфейс находится по адресу `http://localhost:8080` .