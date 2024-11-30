# SETUP

create `.env` from env.examples, make necessary changes if required.

```
cp env.example .env
```

## RUN PROJECT WITH DOCKER
```
docker-compose up --build -d
```

## RUN PROJECT WITHOUT DOCKER

```
go mod tidy
go run main.go
```