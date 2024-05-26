# Используем официальный образ Go
FROM golang:1.22-alpine

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем все файлы проекта в контейнер
COPY . .

# Собираем приложение
RUN go build -o agent ./cmd/agent/main.go

# Запускаем приложение
CMD ["./agent"]