package main

import (
	"github.com/gorilla/mux"
	"github.com/kabukky/httpscerts"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"time"

	"github.com/mbach04/meeseeks/handlers"
)

func main() {
	log.SetFlags(log.LstdFlags)
	router := mux.NewRouter().StrictSlash(true)


	v1, err := readConfig("meeseeks", map[string]interface{}{
		"api_port": 8080,
		"debug":    true,
		"tls_cert": "cert.pem",
		"tls_key":  "key.pem",
	})
	if err != nil {
		log.Println("Error reading config file:", err)
	}

	// Store CONFIG info
	debug := v1.GetBool("debug")
	apiPort := v1.GetString("API_PORT")
	cert := v1.GetString("tls_cert")
	key := v1.GetString("tls_key")

	// Dump CONFIG info to Log
	log.Println("DEBUG:", debug)
	log.Println("API_PORT:", apiPort)
	log.Println("CERT", cert)
	log.Println("KEY", key)

	// Check if the cert files are available.
	err = httpscerts.Check(cert, key)
	// If they are not available, generate new ones.
	if err != nil {
		err = httpscerts.Generate(cert, key, "localhost:"+apiPort)
		if err != nil {
			log.Fatal("Error: Couldn't create https certs.")
		}
	}

	srv := &http.Server{
		Handler: router,
		Addr:    ":" + apiPort,
		// Good practice: enforce timeouts for servers
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	//namespace all api calls with base of <host>:<port>/api/v1
	sub := router.PathPrefix("/api/v1").Subrouter()

	//API Endpoints
	sub.Methods("GET").Path("/hello").HandlerFunc(handlers.GetHello)
	sub.Methods("POST").Path("/bash").HandlerFunc(handlers.Bash)
	sub.Methods("POST").Path("/ls").HandlerFunc(handlers.LsCmd)

	log.Println("Listening on:", apiPort)
	log.Fatal(srv.ListenAndServeTLS(cert, key))
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
