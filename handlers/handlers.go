package handlers

import (
	"encoding/json"
	// "os"
	"net/http"
	// "io/ioutil"
	// "log"
	// "github.com/gorilla/mux"
)


type ApiResponse struct {
	Request  string `json:"request"`
	Response string `json:"response"`
}


// type Files struct {
// 	Files []File
// }

// type File struct {
// 	Id    int
// 	Name  string
// 	Bytes int64
// }

// type FilesResponse struct {
// 	Request  string `json:"request"`
// 	Response Files  `json:"response"`
// }

func writeJsonResponse(w http.ResponseWriter, bytes []byte) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(bytes)
}

func GetHello(w http.ResponseWriter, r *http.Request) {
	response := ApiResponse{Request: "HELLO", Response: "WORLD"}
	bytes, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	writeJsonResponse(w, bytes)
}