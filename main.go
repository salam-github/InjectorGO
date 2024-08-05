package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
)

var embeddedProgram1 []byte
var embeddedProgram2 []byte

func main() {
	var file1, file2 string
	var err error

	if filepath.Base(os.Args[0]) == "injector.exe" {
		file1, file2, err = checkArgs()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		embeddedProgram1, err = os.ReadFile(file1)
		if err != nil {
			panic(err)
		}
		embeddedProgram2, err = os.ReadFile(file2)
		if err != nil {
			panic(err)
		}

		cmd := exec.Command("go", "build", "-o", "newprogram.exe", "main.go")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Run()
		if err != nil {
			panic(err)
		}
		fmt.Println("Successfully created newprogram.exe with embedded executables")
		return
	}

	writeAndExecuteEmbeddedFiles()
}

func checkArgs() (string, string, error) {
	args := os.Args[1:]
	if len(args) != 2 {
		return "", "", fmt.Errorf("Usage: %s <file1> <file2>", os.Args[0])
	}
	return args[0], args[1], nil
}

func writeAndExecuteEmbeddedFiles() {
	tmpfile1, err := ioutil.TempFile("", "embedded_program1_*.exe")
	if err != nil {
		panic(err)
	}
	defer os.Remove(tmpfile1.Name())

	_, err = tmpfile1.Write(embeddedProgram1)
	if err != nil {
		panic(err)
	}
	err = tmpfile1.Close()
	if err != nil {
		panic(err)
	}

	tmpfile2, err := ioutil.TempFile("", "embedded_program2_*.exe")
	if err != nil {
		panic(err)
	}
	defer os.Remove(tmpfile2.Name())

	_, err = tmpfile2.Write(embeddedProgram2)
	if err != nil {
		panic(err)
	}
	err = tmpfile2.Close()
	if err != nil {
		panic(err)
	}

	err = os.Chmod(tmpfile1.Name(), 0755)
	if err != nil {
		panic(err)
	}

	err = os.Chmod(tmpfile2.Name(), 0755)
	if err != nil {
		panic(err)
	}

	cmd1 := exec.Command(tmpfile1.Name())
	cmd1.Stdout = os.Stdout
	cmd1.Stderr = os.Stderr
	err = cmd1.Start()
	if err != nil {
		panic(err)
	}

	cmd2 := exec.Command(tmpfile2.Name())
	cmd2.Stdout = os.Stdout
	cmd2.Stderr = os.Stderr
	err = cmd2.Start()
	if err != nil {
		panic(err)
	}

	err = cmd1.Wait()
	if err != nil {
		panic(err)
	}

	err = cmd2.Wait()
	if err != nil {
		panic(err)
	}
}
