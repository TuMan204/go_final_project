FROM golang:1.26 AS build

WORKDIR /todo_list

COPY . .

RUN go mod tidy

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app main.go

FROM ubuntu:latest AS production

WORKDIR /root/todo_app
COPY --from=build ./app .
ADD ./web ./web
ADD ./scheduler.db .

ENV TODO_PORT=7540 \
    TODO_DBFILE=scheduler.db \
    TODO_PASSWORD=GolangForever

EXPOSE 7540

CMD [ "./app" ]