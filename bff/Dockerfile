FROM golang:latest

EXPOSE 3000

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download
RUN go install github.com/air-verse/air@latest

ENTRYPOINT ["air"]