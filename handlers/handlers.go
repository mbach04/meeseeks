package handlers

import (
	"github.com/labstack/echo"
	"github.com/mbach04/meeseeks/api"
	"log"
	"net/http"
)

/*
-----------------------------------------------------------------
		Request Structs
-----------------------------------------------------------------
*/

type LinuxCommand struct {
	Command string `json:"command"`
	Args    string `json:"args"`
}

type LsReq struct {
	Path string `json:"path"`
}

/*
-----------------------------------------------------------------
		Response Structs
-----------------------------------------------------------------
*/
type ApiResponse struct {
	Request  string              `json:"request"`
	Response api.LsCommandReturn `json:"response"`
}

/*
-----------------------------------------------------------------
		Handler Funcs
-----------------------------------------------------------------
*/


//LsCommand runs an `ls` on the provided `path`
func LsCmd(c echo.Context) (err error) {
	p := new(LsReq)
	if err = c.Bind(p); err != nil {
		log.Println("ERROR", err)
		return
	}
	r := &ApiResponse{
		Request:  "ls " + p.Path,
		Response: api.LsCommand(p.Path),
	}
	return c.JSONPretty(http.StatusOK, r, "  ")
}
