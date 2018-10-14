[![Go Report Card](https://goreportcard.com/badge/github.com/mbach04/meeseeks)](https://goreportcard.com/report/github.com/mbach04/meeseeks)

# Meeseeks
Linux system administration served over RESTful JSON API


### Compile for Linux 64 Bit
`GOOS=linux GOARCH=amd64 go build -gcflags=-trimpath=$GOPATH -asmflags=-trimpath=$GOPATH meeseeks.go`

### Dependent Packages
```
 github.com/gorilla/mux
 github.com/spf13/viper
 github.com/kabukky/httpscerts
```

### Run local (no compile)
Configure meeseeks.yml with the API port you'd like to use then:
`go run meeseeks.go`

### Curl the `ls` endpoint by `POST`ing a path
`curl -d '{"path": "/home/user"}' localhost:9191/api/v1/ls`

### Creating TLS certs for testing purposes
```openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout key.pem -out cert.pem```
Reference the location of `cert.pem` and `key.pem` in your `meeseeks.yml` config file.
