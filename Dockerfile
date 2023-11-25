FROM golang:latest
LABEL authors="cha2ned"

WORKDIR /app
COPY . .

RUN go install github.com/go-task/task/v3/cmd/task@latest
ENTRYPOINT ["task", "run"]