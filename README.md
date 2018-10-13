# Meeseeks
Linux system administration served over RESTful JSON API


### Compile for Linux 64 Bit
GOOS=linux GOARCH=amd64 go build -gcflags=-trimpath=$GOPATH -asmflags=-trimpath=$GOPATH meeseeks.go

### Run local (no compile)
Configure api.yml with the API port you'd like to use then:
`go run meeseeks.go`

### Curl the `ls` endpoint by POSTing a path
`curl -d '{"path": "/home/user"}' localhost:9191/api/v1/ls`