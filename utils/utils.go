package osutils

import (
    "bytes"
    "log"
    "os/exec"
    "syscall"
)

const defaultFailedCode = 1

type CommandReturn struct {
	Stdout      string
	Stderr      string
	ExitCode    int
}

func RunCommand(name string, args ...string) (result CommandReturn) {
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
    log.Printf("command result, stdout: %v, stderr: %v, exitCode: %v", stdout, stderr, exitCode)
    result.Stdout = stdout
    result.Stderr = stderr
    result.ExitCode = exitCode
    return
}
