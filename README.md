# GOAPP

Golang stack for backend

## 命令说明

- docker build

`docker build -t goapp .`

- docker run

`docker run --name goapp goapp:latest` OR `docker run --name goapp goapp:latest goapp serve -p 6666`

- **查看帮助:**

`go run main.go -h`

- **编译:**

`GOOS=linux make` or `GOOS=windows make`



