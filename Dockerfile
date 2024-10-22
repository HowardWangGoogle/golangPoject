#Build stage
FROM golang:1.23.2-alpine3.20 AS builder
WORKDIR /app
COPY . .

RUN go build -o main main.go
RUN apk add --no-cache curl
RUN apk add --no-cache jq
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.14.1/migrate.linux-amd64.tar.gz | tar xvz


FROM alpine:3.20
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/migrate.linux-amd64 ./migrate
COPY app.env .
COPY start.sh .
COPY wait-for.sh .
COPY db/migration ./migration

RUN apk add --no-cache jq

RUN chmod +x /app/start.sh /app/wait-for.sh


EXPOSE 8080
ENTRYPOINT [ "/app/start.sh" ]
CMD [ "/app/main" ]
