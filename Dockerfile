# Build stage
FROM golang:1.24-alpine AS build-env
WORKDIR /app
COPY . .
RUN go build -o main main.go
RUN apk add --no-cache curl
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.16.2/migrate.linux-amd64.tar.gz | tar xvz

# Run stage
FROM alpine:3.18
WORKDIR /app
COPY --from=build-env /app/main .
COPY --from=build-env /app/migrate ./migrate
COPY app.env .
COPY start.sh .
COPY wait-for.sh .
COPY db/migration ./migration

EXPOSE 8080
CMD [ "/app/main" ]
ENTRYPOINT [ "/app/start.sh" ]