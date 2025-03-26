# Используем официальный образ Go для сборки приложения
FROM golang:1.24.1 AS builder

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем все файлы в контейнер
COPY . .

# Загружаем зависимости и строим приложение
RUN go mod tidy
RUN go build -o mafia-bot

# Используем Alpine для финального образа
FROM alpine:latest

# Устанавливаем libc6-compat (для glibc)
RUN apk add --no-cache libc6-compat

# Копируем скомпилированный бинарник из стадии сборки
COPY --from=builder /app/mafia-bot .

# Открываем порт 8080
EXPOSE 8080

# Команда для запуска приложения
CMD ["/mafia-bot"]
