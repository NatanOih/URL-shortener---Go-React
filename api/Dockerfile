


# Stage 1: Build the React app
FROM node:22-alpine AS frontend-build

WORKDIR /app

COPY ../frontend/package*.json ./frontend
RUN npm install
COPY ../frontend ./frontend
RUN npm run build


FROM golang:1.22-alpine 
RUN go install github.com/githubnemo/CompileDaemon@latest

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN chmod -R +x /app

EXPOSE 3000

ENTRYPOINT CompileDaemon --build="go build -o build/goapp" -command="./build/goapp" -build-dir=/app -polling=true