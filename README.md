# Test task for Golang developer in iContext

[The task - REST API](Task.md)

## Clone this repo
```
git clone https://github.com/richtheme/iContext
```
and go to project's folder:
```
cd iContext
```
Configure ```/configs/api.toml```

## To run this application you can choose:
## 1. Using Make
```
make run host=<redis_host> port=<redis_port>
```
Example:
```
make run host=localhost port=6379
```

## 2. Using go run
```
go run ./cmd/api/main.go -host <redis_host> -port <redis_port>
```
Example:
```
go run ./cmd/api/main.go -host localhost -port 6379
```
## 3. Using Docker

Build container:
```
docker build -t i_context . 
```
Open bash in this container:
```
docker run -p 8080:8080 -ti i_context /bin/bash
```

And run app:
```
./iContext -host <redis_host> -port <redis_port>
```
Example if redis in docker:
```
./iContext -host host.docker.internal -port 6379
```


## To run tests
Change redisAddr in file ```/internal/api/api_test.go```
## 1. Using Make
```
make test
```
## 2. Using go test
```
go test ./...
```
or
```
go test -v ./...
```


## If first launch you can create table
Change credentials to your postgresql server and run:
```
make migrate
```

