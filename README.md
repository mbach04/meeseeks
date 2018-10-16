[![Go Report Card](https://goreportcard.com/badge/github.com/mbach04/meeseeks)](https://goreportcard.com/report/github.com/mbach04/meeseeks)
[![Build Status](https://travis-ci.com/mbach04/meeseeks.svg?branch=master)](https://travis-ci.com/mbach04/meeseeks)

# Meeseeks
Linux system administration served over RESTful JSON API

### Whats with the name Meeseeks?
Meeseeks are creatures first introduced in Rick and Morty in the fifth episode of the first season. Meeseeks perform a task and upon completion are promptly removed from existence. This is in line with how Linux system administration could operate at scale. Instead of having large overhead with increasing numbers of concurrent SSH tunnels or resource waste with agent to master communication setups (ie Puppet), we can streamline tasks with an SSL/TLS based implementation that keeps the connections light and short lived. Most systems talk over HTTP(S) and this is an effort to make linux system administration do the same. Much can be said for the flexibility offered by being able to auth to a fleet of linux hosts and then proceed to administer them with something as simple as a curl statement including a json web token Bearer. This opens up the possibilities for the managing system implementation signficantly. You could control a fleet of hosts with nothing more than the Ansible URL module for example. So, much like a meeseeks, you make a request, the system does a thing, the response is sent, and the transactoin is complete.

### Project Kanban Board
Get a very early idea of where this is headed on the project board or look for ways you can contribute.
[Meeseeks Board](https://github.com/mbach04/meeseeks/projects/1)

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

### Login to the server and test a secure endpoint
`export TOKEN=$(curl -k -X POST -d 'username=admin' -d 'password=secret' https://localhost:9191/login | jq -r ".token")`  
`curl localhost:1323/restricted -H "Authorization: Bearer ${TOKEN}"`

### Curl the `ls` endpoint by `POST`ing a path
```curl -k -X POST  https://localhost:9191/api/v1/ls -H "Authorization: Bearer ${TOKEN}" -H 'Content-Type: application/json' -d '{"path": "/Users"}'```  

### Creating TLS certs for testing purposes
```openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout key.pem -out cert.pem```
Reference the location of `cert.pem` and `key.pem` in your `meeseeks.yml` config file.  
Or use Go's built in crypto package to create certs:
`go run $GOROOT/src/crypto/tls/generate_cert.go --host localhost`
