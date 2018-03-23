package main

import (
	"log"
	"os/exec"
)

func runDiscardOutput(dir string, name string, arg ...string) error {

	//err, stderr, stdout := run(cmd, false, false, false)
	_, _, err := run(false, false, false, dir, name, arg...)
	return err
}

func run(stderr, stdout, seperate bool, dir, name string, arg ...string) (string, string, error) {
	//cmd, err := exec.Command(name, arg...).Output()
	//if dir != "" {
	//cmd.Dir = dir
	//}
	//err := cmd.Run()
	//log.Println("--------------", cmd.ProcessState)

	//str, _ := cmd.Output()

	//return string(cmd), "", err

	//cmd := exec.Command("sh", "-c", "echo stdout; echo 1>&2 stderr")

	cmd := exec.Command(name, arg...)
	if dir != "" {
		cmd.Dir = dir
	}
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		log.Println("Error running:", name, arg)
		log.Println(err)
	}
	return string(stdoutStderr), "", err
}

func runCaptureStdout(dir string, name string, arg ...string) (string, error) {
	stdout, _, err := run(true, false, false, dir, name, arg...)
	return stdout, err
}

func runCaptureStderr(dir string, name string, arg ...string) (string, error) {
	_, stderr, err := run(false, true, false, dir, name, arg...)
	return stderr, err
}

// cmd := exec.Command("myCommand", "arg1", "arg2")
// cmd.Dir = "/path/to/work/dir"
// cmd.Run()

func runCaptureAll(dir string, name string, arg ...string) (string, string, error) {
	stdout, stderr, err := run(true, true, false, dir, name, arg...)
	return stdout, stderr, err
}
