# Используем образ Golang для сборки приложения
FROM golang:1.21.1 AS builder

# Установка рабочей директории внутри контейнера
WORKDIR /app

# Копируем файлы go.mod и go.sum для загрузки зависимостей
COPY go.mod . 
COPY go.sum .

# Загрузка зависимостей с помощью go mod download
RUN go mod download

# Копируем весь проект в контейнер
COPY . .

# Установка рабочей директории на папку с main.go
WORKDIR /app/app

# Запускаем main.go
CMD ["go", "run", "main.go"]
