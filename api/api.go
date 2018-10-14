package api

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
)

const defaultFailedCode = 1

/*
-----------------------------------------------------------------
		Structs that build the json response of a command call
-----------------------------------------------------------------
*/
type CommandReturn struct {
	Stdout   string `json:"stdOut"`
	Stderr   string `json:"stdErr"`
	ExitCode int    `json:"exitCode"`
}

type LsCommandReturn struct {
	ExitCode int    `json:"exitCode"`
	BasePath string `json:"basePath"`
	Files    []File `json:"files"`
	Count    int    `json:"count"`
	Stderr   string `json:"stderr"`
}

/*
-----------------------------------------------------------------
		Helper structs that are shared
-----------------------------------------------------------------
*/

//Files is a list of File structs
type Files struct {
	Files []File `json:"files"`
}

//File is 1 file
type File struct {
	Name    string `json:"name"`
	Bytes   int64  `json:"bytes"`
	Type    string `json:"type"`
	ModTime string `json:"modTime"`
	Perms   string `json:"perms"`
}

/*
-----------------------------------------------------------------
        Commands that do the ~things~
        Naming convention:
            xCommand where `x` is what this thing does
            Think of it as a parallel to the binary you would
            normally use on a command line
-----------------------------------------------------------------
*/
//Bash is a feature in flux as this poses infinite security risks
//but resides currently as a proof of concept to show how the wiring
//of the various parts work together
func Bash(name string, args ...string) CommandReturn {
	log.Println("RUN COMMAND:", name, args)
	var outbuf, errbuf bytes.Buffer
	var exitCode int
	cmd := exec.Command(name, args...)
	cmd.Stdout = &outbuf
	cmd.Stderr = &errbuf

	err := cmd.Run()
	stdout := outbuf.String()
	stderr := errbuf.String()

	if err != nil {
		// try to get the exit code
		if exitError, ok := err.(*exec.ExitError); ok {
			ws := exitError.Sys().(syscall.WaitStatus)
			exitCode = ws.ExitStatus()
		} else {
			// This will happen (in OSX) if `name` is not available in $PATH,
			// in this situation, exit code could not be gotten, and stderr will be
			// empty string very likely, so we use the default fail code, and format err
			// to string and set to stderr
			log.Printf("Could not get exit code for failed program: %v, %v", name, args)
			exitCode = defaultFailedCode
			if stderr == "" {
				stderr = err.Error()
			}
		}
	} else {
		// success, exitCode should be 0 if go is ok
		ws := cmd.ProcessState.Sys().(syscall.WaitStatus)
		exitCode = ws.ExitStatus()
	}

	log.Printf("stdout: %v", stdout)
	log.Printf("stderr: %v", stderr)
	log.Printf("exitcode: %v", exitCode)

	result := new(CommandReturn)
	result.Stdout = stdout
	result.Stderr = stderr
	result.ExitCode = exitCode
	return *result
}

//LsCommand implements an api abstraction of a standard ls command with the args
// -lsah on a given path with the following exceptions:
//size is in bytes and let to the caller to format into a human readable context
//currently not implemented: [perms, owner]
func LsCommand(path string) LsCommandReturn {
	flist, err := ioutil.ReadDir(path)
	// var filesStruct Files
	exitCode := 0
	count := 1
	stdErr := ""
	lsReturn := new(LsCommandReturn)

	if err != nil {
		log.Println("Error reading path:", path, err)
		stdErr = fmt.Sprintf("%v", err)
		exitCode = defaultFailedCode
	}

	//TODO: Recurse into directory objects?
	//Maybe even a `tree` command style response would be useful
	//and include an option to limit recurse level
	for _, f := range flist {
		ftype := ""
		fi, err := os.Lstat(filepath.Join(path, f.Name()))
		if err != nil {
			log.Println(err)
			continue //skip this file on error
		}
		// p := fi.Mode()
		switch mode := fi.Mode(); {
		case mode.IsRegular():
			ftype = "regular"
		case mode.IsDir():
			ftype = "dir"
		case mode&os.ModeSymlink != 0:
			ftype = "sym"
		case mode&os.ModeNamedPipe != 0:
			ftype = "namePipe"
		}

		lsReturn.Files = append(lsReturn.Files,
			File{Name: f.Name(),
				Bytes:   f.Size(),
				Type:    ftype,
				ModTime: f.ModTime().String(),
				Perms:   fmt.Sprintf("%04o", f.Mode().Perm()),
			})
		count++
	}
	lsReturn.ExitCode = exitCode
	lsReturn.Stderr = stdErr
	lsReturn.Count = count
	lsReturn.BasePath = path
	return *lsReturn
}

/*
-----------------------------------------------------------------
		Helper funcs to the commands go here
-----------------------------------------------------------------
*/

func pathExists(path string) bool {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			log.Println("Does not exist", path)
			return false
		}
	}
	return true
}
