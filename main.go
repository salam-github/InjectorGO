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
	if filepath.Base(os.Args[0]) == "injector.exe" {
		if len(os.Args) == 3 {
			file1, file2, err := checkArgs()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			err = embedAndCreateExecutable(file1, file2)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			fmt.Println("Successfully created newprogram.exe.")
		} else {
			fmt.Printf("Usage: %s <file1> <file2>\n", os.Args[0])
			os.Exit(1)
		}
	} else {
		// When running newprogram.exe
		embeddedProgram1, err := os.ReadFile("important1.exe")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		embeddedProgram2, err := os.ReadFile("important2.exe")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		writeAndExecuteTemporaryFiles(embeddedProgram1, embeddedProgram2)
	}
}

func checkArgs() (string, string, error) {
	args := os.Args[1:]
	if len(args) != 2 {
		return "", "", fmt.Errorf("Usage: %s <file1> <file2>", os.Args[0])
	}
	return args[0], args[1], nil
}

func embedAndCreateExecutable(file1, file2 string) error {
	f1, err := os.ReadFile(file1)
	if err != nil {
		return err
	}
	f2, err := os.ReadFile(file2)
	if err != nil {
		return err
	}

	err = os.WriteFile("important1.exe", f1, 0755)
	if err != nil {
		return err
	}
	err = os.WriteFile("important2.exe", f2, 0755)
	if err != nil {
		return err
	}

	cmd := exec.Command("go", "build", "-o", "newprogram.exe", "main.go")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

func writeAndExecuteTemporaryFiles(embeddedProgram1, embeddedProgram2 []byte) {
	tmpfile1, err := ioutil.TempFile("", "embedded_program1_*.exe")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer os.Remove(tmpfile1.Name())

	_, err = tmpfile1.Write(embeddedProgram1)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = tmpfile1.Close()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	tmpfile2, err := ioutil.TempFile("", "embedded_program2_*.exe")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer os.Remove(tmpfile2.Name())

	_, err = tmpfile2.Write(embeddedProgram2)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = tmpfile2.Close()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = os.Chmod(tmpfile1.Name(), 0755)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = os.Chmod(tmpfile2.Name(), 0755)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	cmd1 := exec.Command(tmpfile1.Name())
	cmd1.Stdout = os.Stdout
	cmd1.Stderr = os.Stderr
	err = cmd1.Run()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	cmd2 := exec.Command(tmpfile2.Name())
	cmd2.Stdout = os.Stdout
	cmd2.Stderr = os.Stderr
	err = cmd2.Run()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
