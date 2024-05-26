# Используем официальный образ Go
FROM golang:1.22-alpine

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем все файлы проекта в контейнер
COPY . .

# Собираем приложение
RUN go build -o orchestrator ./cmd/orchestrator/main.go

# Открываем порт для оркестратора
EXPOSE 8080

# Запускаем приложение
CMD ["./orchestrator"]