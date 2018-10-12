package utils

import (
	"io/ioutil"
    "bytes"
    "log"
    "os/exec"
    "os"
    "syscall"
    // "strings"
    "path/filepath"
)

const defaultFailedCode = 1
const (
	//B 1 byte
	B  = 1
	//KB 1 KiloByte
	KB = 1024 * B
	//MB 1 MegaByte
	MB = 1024 * KB
	//GB 1 GigaByte
	GB = 1024 * MB
)

type CommandReturn struct {
	Stdout      string  `json:"stdOut"`
	Stderr      string  `json:"stdErr"`
	ExitCode    int     `json:"exitCode"`
}

type LsCommandReturn struct {
    ExitCode    int         `json:"exitCode"`
    BasePath    string      `json:"basePath"`
    Files       Files       `json:"files"`
}

//Files is a list of File structs
type Files struct {
	Files []File	`json:"files"`
}

//File is 1 file name, and size in bytes
type File struct {
	Name    string	        `json:"name"`
    Bytes   int64		    `json:"bytes"`
    Mode    os.FileMode     `json:"mode"`
    IsDir   bool            `json:"isDir"`

}

func RunCommand(name string, args ...string) CommandReturn {
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

func LsCommand(path string) LsCommandReturn {
	flist, err := ioutil.ReadDir(path)
    var filesStruct Files
    exitCode := 0

    if err != nil {
        log.Println("Error reading path:", path, err)
        exitCode = defaultFailedCode
    }
    
	//TODO: Recurse into directory objects? Maybe a `tree` command style response would be useful and include an option to limit recurse level
	for _, f := range flist {
		fname := f.Name()
        fsize := f.Size()
        fmode := f.Mode()
        fisdir := f.IsDir()
		if f.IsDir() {
			fsize = int64(dirSizeBytes(filepath.Join(path, f.Name())))
		}
		filesStruct.Files = append(filesStruct.Files, File{Name: fname, Bytes: fsize, Mode: fmode, IsDir: fisdir})
    }
    
	return LsCommandReturn{exitCode, path, filesStruct}
}


//dirSizeBytes returns a directories size in bytes
//note: this func is not pre-checking if the path provided is
//a directory instead of a file
func dirSizeBytes(path string) float64 {
    var dirSize int64
    readSize := func(path string, file os.FileInfo, err error) error {
        if !file.IsDir() {
            dirSize += file.Size()
        }
        return nil
    }
    filepath.Walk(path, readSize)    
	sizeBytes := float64(dirSize)
    return sizeBytes
}

func pathExists(path string) bool {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			log.Println("UNABLE TO FIND", path)
			return false
		}
	}
	return true
}