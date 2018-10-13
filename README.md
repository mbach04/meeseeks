# Meeseeks
Linux system administration served over RESTful JSON API


### Compile for Linux 64 Bit
GOOS=linux GOARCH=amd64 go build -gcflags=-trimpath=$GOPATH -asmflags=-trimpath=$GOPATH meeseeks.go

### Run local
go run meeseeks.go