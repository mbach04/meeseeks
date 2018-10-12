package handlers

import (
	"encoding/json"
	// "os"
	"net/http"
	"io/ioutil"
	"log"
	// "github.com/gorilla/mux"
	"github.com/mbach04/stunning-octo-garbanzo/utils"
)


type ApiResponse struct {
	Request  string `json:"request"`
	Response string `json:"response"`
}

type LinuxCommand struct {
	Command	string	`json:"command"`
	Args	string `json:"args"`
}

type LinuxCommandResponse struct {
	Stdout	string	`json:"stdout"`
	Stderr	string	`json:"stderr"`
	Exitcode int	`json:"exitcode"`
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

func RunCommand(w http.ResponseWriter, r *http.Request) {
	log.Println("POST: /command:", r.RemoteAddr)
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	lc := new(LinuxCommand)
	err = json.Unmarshal(body, lc)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	stdout, stderr, exitCode := utils.RunCommand(lc.Command, lc.Args)
	response := LinuxCommandResponse{Stdout: stdout, Stderr: stderr, Exitcode: exitCode}
	bytes, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	writeJsonResponse(w, bytes)
}