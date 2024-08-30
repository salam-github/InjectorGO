package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"time"

	"github.com/google/uuid"
)

func main() {
	// Display splash screen and menu
	showSplashScreen()
	handleUserInput()
}

// showSplashScreen displays a splash screen.
func showSplashScreen() {
	// Red color ANSI escape code is \033[31m
	// Reset color ANSI escape code is \033[0m
	fmt.Println("\033[31m")
	fmt.Println(`
 ██▓ ███▄    █  ▄▄▄██▀▀▀▓█████  ▄████▄  ▄▄▄█████▓ ▒█████   ██▀███  
▓██▒ ██ ▀█   █    ▒██   ▓█   ▀ ▒██▀ ▀█  ▓  ██▒ ▓▒▒██▒  ██▒▓██ ▒ ██▒
▒██▒▓██  ▀█ ██▒   ░██   ▒███   ▒▓█    ▄ ▒ ▓██░ ▒░▒██░  ██▒▓██ ░▄█ ▒
░██░▓██▒  ▐▌██▒▓██▄██▓  ▒▓█  ▄ ▒▓▓▄ ▄██▒░ ▓██▓ ░ ▒██   ██░▒██▀▀█▄  
░██░▒██░   ▓██░ ▓███▒   ░▒████▒▒ ▓███▀ ░  ▒██▒ ░ ░ ████▓▒░░██▓ ▒██▒
░▓  ░ ▒░   ▒ ▒  ▒▓▒▒░   ░░ ▒░ ░░ ░▒ ▒  ░  ▒ ░░   ░ ▒░▒░▒░ ░ ▒▓ ░▒▓░
 ▒ ░░ ░░   ░ ▒░ ▒ ░▒░    ░ ░  ░  ░  ▒       ░      ░ ▒ ▒░   ░▒ ░ ▒░
 ▒ ░   ░   ░ ░  ░ ░ ░      ░   ░          ░      ░ ░ ░ ▒    ░░   ░ 
 ░           ░  ░   ░      ░  ░░ ░                   ░ ░     ░     
                               ░                                   
`)
	fmt.Println("\033[0m")
	fmt.Println("\033[32m")
	fmt.Println("====================================")
	fmt.Println("Welcome to the Injector!")
	fmt.Println("This tool combines multiple executables into a single bundle.")
	fmt.Println("You will be guided through a menu to specify your target and executables.")
	fmt.Println("====================================\n")
	fmt.Println("\033[0m")
	time.Sleep(2 * time.Second) // Pause for effect
}

// handleUserInput handles user interaction through a CLI menu.
func handleUserInput() {
	operatingSystem := checkOperatingSystem()

	for {
		fmt.Println("1. Create a new bundle")
		fmt.Println("2. Exit")
		fmt.Print("Select an option (1 or 2): ")

		var choice string
		fmt.Scanln(&choice)

		switch choice {
		case "1":
			targetFileName, executablePaths := getTargetAndExecutables()
			validateArguments(targetFileName, executablePaths, operatingSystem)
			executableData := collectExecutableData(executablePaths)
			tempFileName := createBundle(targetFileName, executableData)
			appendToBundle(targetFileName, executableData)
			logInfo("Successfully combined " + intToString(len(executablePaths)) + " executables.")

			// Ask user if they want to save the temporary .go file
			fmt.Println("Do you want to save the temporary .go file? (y/n)")
			var saveResponse string
			fmt.Scanln(&saveResponse)
			if saveResponse == "y" {
				saveTempFile(tempFileName)
			} else {
				removeFile(tempFileName)
			}
		case "2":
			fmt.Println("Exiting the program.")
			return
		default:
			logWarning("Invalid option selected. Please choose 1 or 2.")
		}
	}
}

// getTargetAndExecutables prompts the user for target name and executable paths.
func getTargetAndExecutables() (string, []string) {
	fmt.Print("Enter the target file name (e.g., target.exe): ")
	var targetFileName string
	fmt.Scanln(&targetFileName)

	fmt.Print("Enter the number of executables: ")
	var numExecutables int
	fmt.Scanln(&numExecutables)

	var executablePaths []string
	for i := 0; i < numExecutables; i++ {
		fmt.Printf("Enter path for executable %d: ", i+1)
		var path string
		fmt.Scanln(&path)
		executablePaths = append(executablePaths, path)
	}

	return targetFileName, executablePaths
}

// checkOperatingSystem determines the OS and returns its name.
func checkOperatingSystem() string {
	logInfo("Determining the operating system...")

	switch runtime.GOOS {
	case "windows":
		return "windows"
	case "linux":
		return "linux"
	case "darwin":
		return "macos"
	default:
		logFatal(errors.New("Unsupported OS. Supported OS's are Windows, Linux, and macOS."))
		return ""
	}
}

// validateArguments checks if files exist and confirms if overwriting is intended.
func validateArguments(targetFileName string, executablePaths []string, operatingSystem string) {
	logInfo("Validating arguments...")

	for _, path := range executablePaths {
		if !fileExists(path) {
			logFatal(errors.New("File " + path + " does not exist."))
		}
	}

	if fileExists(targetFileName + getOSFileExtension(operatingSystem)) {
		logWarning("File " + targetFileName + " already exists. It will be overwritten. Continue? (y/n)")
		var response string
		for response != "y" && response != "n" {
			fmt.Scanln(&response)
		}
		if response == "n" {
			logFatal(errors.New("Operation cancelled by the user."))
		}
	}
}

// fileExists checks if a file exists.
func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !errors.Is(err, os.ErrNotExist)
}

// getOSFileExtension returns the file extension based on the operating system.
func getOSFileExtension(osName string) string {
	if osName == "windows" {
		return ".exe"
	}
	return ""
}

// collectExecutableData reads the content of the executable files.
func collectExecutableData(paths []string) [][]byte {
	logInfo("Collecting data from executable files...")

	var data [][]byte
	for _, path := range paths {
		data = append(data, []byte("h$lL@ w@rLd"))
		data = append(data, readFileContents(path))
	}

	return data
}

// createBundle creates a bundle executable file with a unique name.
func createBundle(name string, data [][]byte) string {
	logInfo("Creating bundle executable...")

	tempFileName := uuid.NewString() + ".go"
	tempFile := openFile(tempFileName)
	writeToFile(tempFileName, tempFile, generateBundleProgram())
	tempFile.Close()

	cmd := exec.Command("go", "build", "-o", "./"+name, tempFileName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		logFatal(errors.New("Failed to compile bundle: " + err.Error() + "\nOutput:\n" + string(output)))
	}

	return tempFileName
}

// appendToBundle appends executable data to the bundle file.
func appendToBundle(name string, data [][]byte) {
	logInfo("Appending executable data to the bundle...")

	bundleFile := openFile(name)
	defer bundleFile.Close()

	for _, content := range data {
		writeToFile(name, bundleFile, content)
	}
}

// readFileContents reads the content of a file.
func readFileContents(name string) []byte {
	contents, err := ioutil.ReadFile(name)
	if err != nil {
		logFatal(errors.New("Failed to read file " + name + ": " + err.Error()))
	}
	return contents
}

// openFile opens a file with append and write permissions.
func openFile(name string) *os.File {
	file, err := os.OpenFile(name, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		logFatal(errors.New("Failed to open file " + name + ": " + err.Error()))
	}
	return file
}

// writeToFile writes content to a file.
func writeToFile(name string, file *os.File, content []byte) {
	_, err := file.Write(content)
	if err != nil {
		logFatal(errors.New("Failed to write to file " + name + ": " + err.Error()))
	}
}

// removeFile deletes a file.
func removeFile(name string) {
	err := os.Remove(name)
	if err != nil {
		logFatal(errors.New("Failed to delete file " + name + ": " + err.Error()))
	}
}

// saveTempFile saves the temporary Go source file.
func saveTempFile(fileName string) {
	destName := fileName + ".saved"
	err := os.Rename(fileName, destName)
	if err != nil {
		logError("Failed to save the temporary .go file: " + err.Error())
	} else {
		logInfo("Temporary .go file saved as " + destName)
	}
}

// generateBundleProgram generates the source code for the bundle executable.
func generateBundleProgram() []byte {
	return []byte(
		`package main
	
import (
	"os"
	"strings"
	"io/ioutil"
	"github.com/amenzhinsky/go-memexec"
	"fmt"
)

func main() {
	bytes, _ := ioutil.ReadFile(os.Args[0])
	executables := strings.Split(string(bytes), "h$lL@ w@rLd")
	for i := 2; i < len(executables); i++ {
		exe, _ := memexec.New([]byte(executables[i]))
		defer exe.Close()
		cmd := exe.Command()
		res, _ := cmd.Output()
		fmt.Print(string(res))
	}
}
`)
}

// logInfo prints an informational message to stdout.
func logInfo(msg string) {
	fmt.Println("\033[32m[INFO]\033[0m [" + currentTime() + "][" + callerLocation() + ":" + callerLine() + "]: " + msg)
}

// logWarning prints a warning message to stdout.
func logWarning(msg string) {
	fmt.Println("\033[33m[WARNING]\033[0m [" + currentTime() + "][" + callerLocation() + ":" + callerLine() + "]: " + msg)
}

// logError prints an error message to stdout.
func logError(msg string) {
	fmt.Println("\033[31m[ERROR]\033[0m [" + currentTime() + "][" + callerLocation() + ":" + callerLine() + "]: " + msg)
}

// logFatal prints a fatal error message and exits the program.
func logFatal(err error) {
	fmt.Println("\033[41m[FATAL ERROR]\033[0m [" + currentTime() + "][" + callerLocation() + ":" + callerLine() + "]: " + err.Error())
	fmt.Println("\033[41m[EXITING]\033[0m")
	os.Exit(1)
}

// currentTime gets the current time in "01-02-2006 15:04:05.00000" format.
func currentTime() string {
	dt := time.Now()
	return dt.Format("01-02-2006 15:04:05.00000")
}

// callerLine returns the line number of the caller.
func callerLine() string {
	_, _, line, _ := runtime.Caller(2)
	return strconv.Itoa(line)
}

// callerLocation returns the file location of the caller.
func callerLocation() string {
	_, location, _, _ := runtime.Caller(2)
	return location
}

// intToString converts an integer to a string.
func intToString(n int) string {
	return strconv.Itoa(n)
}
