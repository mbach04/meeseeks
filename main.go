package main

import (
	"log"
	"net/http"
	"os"
	"time"
	// "strings"
	"github.com/gorilla/mux"
	"github.com/mbach04/stunning-octo-garbanzo/handlers"
)
func main() {
	log.SetFlags(log.LstdFlags)
	wwwdir := getEnv("HTTP_DIR", "/var/q7_www")
	serverPort := getEnv("api_PORT", "8080")
	router := mux.NewRouter().StrictSlash(true)

	srv := &http.Server{
		Handler: router,
		Addr:    ":" + serverPort,
		// Good practice: enforce timeouts for servers
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	sub := router.PathPrefix("/api/v1").Subrouter()

	sub.Methods("GET").Path("/hello").HandlerFunc(handlers.GetHello)

	//Serve the archive files
	// router.PathPrefix("/download/files/archive").Handler(http.StripPrefix("/download/files/archive", http.FileServer(http.Dir("/media/archive"))))
	// router.PathPrefix("/download/files/collector").Handler(http.StripPrefix("/download/files/collector", http.FileServer(http.Dir("/media/1879-AACB"))))

	// important!: this route must be last in order to prevent overriding all subroutes
	// matched by /* (aka the entire api)
	// This will serve files under http://localhost:8080/<filename>
	router.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir(wwwdir))))


	log.Println("Start API Server")
	log.Println("Listening on:", serverPort)
	log.Fatal(srv.ListenAndServe())
}



func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}
