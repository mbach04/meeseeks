package handlers

import (
	"encoding/json"
	// "os"
	"net/http"
	"io/ioutil"
	"log"
	// "github.com/gorilla/mux"
	"github.com/mbach04/meeseeks/utils"
)

/* 
-----------------------------------------------------------------
		Request Structs
-----------------------------------------------------------------
*/

type LinuxCommand struct {
	Command	string	`json:"command"`
	Args	string `json:"args"`
}

type LsReq struct {
	Path	string	`json:"path"`
}

/* 
-----------------------------------------------------------------
		Response Structs
-----------------------------------------------------------------
*/
type ApiResponse struct {
	Request  string `json:"request"`
	Response string `json:"response"`
}

/* 
-----------------------------------------------------------------
		Handler Funcs
-----------------------------------------------------------------
*/

//GetHello returns a simple `hello world` string as the `response` json key
func GetHello(w http.ResponseWriter, r *http.Request) {
	response := ApiResponse{Response: "Hello World"}
	bytes, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	writeJsonResponse(w, bytes)
}

//RunCommand executes a command on the localhost
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
	response := utils.RunCommand(lc.Command, lc.Args)
	bytes, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	writeJsonResponse(w, bytes)
}

//RunCommand executes a command on the localhost
func LsCmd(w http.ResponseWriter, r *http.Request) {
	log.Println("POST: /ls:", r.RemoteAddr)
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	ls := new(LsReq)
	err = json.Unmarshal(body, ls)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	response := utils.LsCommand(ls.Path)
	bytes, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	writeJsonResponse(w, bytes)
}

/* 
-----------------------------------------------------------------
		Helper funcs to the handler funcs go here 
-----------------------------------------------------------------
*/

func writeJsonResponse(w http.ResponseWriter, bytes []byte) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(bytes)
}