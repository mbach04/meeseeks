package main

import (
	"log"
	"net/http"
	"time"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
	"github.com/mbach04/stunning-octo-garbanzo/handlers"
)


func main() {
	log.SetFlags(log.LstdFlags)
	router := mux.NewRouter().StrictSlash(true)

	v1, err := readConfig("api", map[string]interface{}{
		"debug": true,
		"api__port": 8080})
	if err != nil {
		log.Println("Error reading config file:", err)
	}

	// Store CONFIG info
	debug := v1.GetBool("debug")
	apiPort := v1.GetString("API_PORT")

	// Dump CONFIG info to Log
	log.Println("DEBUG:", debug)
	log.Println("API_PORT:", apiPort)

	srv := &http.Server{
		Handler: router,
		Addr:    ":" + apiPort,
		// Good practice: enforce timeouts for servers
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	//namespace all api calls with base of <host>:<port>/api/v1
	sub := router.PathPrefix("/api/v1").Subrouter()

	sub.Methods("GET").Path("/hello").HandlerFunc(handlers.GetHello)
	//curl localhost:9191/api/v1/hello
	
	sub.Methods("POST").Path("/command").HandlerFunc(handlers.RunCommand)
	//curl -d '{"command": "/bin/sleep", "args": "10"}' localhost:9191/api/v1/command
	
	sub.Methods("POST").Path("/ls").HandlerFunc(handlers.LsCmd)
	//curl -d '{"path": "/Users/go/src/github.com/"}' localhost:9191/api/v1/ls | jq "."

	log.Println("Listening on:", apiPort)
	log.Fatal(srv.ListenAndServe())
}


func readConfig(filename string, defaults map[string]interface{}) (*viper.Viper, error) {
	v := viper.New()
	for key, value := range defaults {
		v.SetDefault(key, value)
	}
	v.SetConfigName(filename)
	v.AddConfigPath(".")
	v.AutomaticEnv()
	err := v.ReadInConfig()
	return v, err
}