INJECTOR
========

INJECTOR is a Go-based command-line tool designed for combining multiple executable files into a single bundled executable. This tool is of particular interest in the field of cybersecurity, where understanding and managing executable files and their interactions is crucial. This README provides a comprehensive overview of the program's functionality, its applications in cybersecurity, and best practices for using it.

Table of Contents
-----------------

- [Features](#features)
- [How It Works](#how-it-works)
- [Installation](#installation)
- [Usage](#usage)
  - [Running the Program](#running-the-program)
  - [Interacting with the Menu](#interacting-with-the-menu)
  - [Arguments and Inputs](#arguments-and-inputs)
- [Program Workflow](#program-workflow)
- [File Structure and Bundling](#file-structure-and-bundling)
- [Security Implications](#security-implications)
- [Testing and Validation](#testing-and-validation)
- [Troubleshooting](#troubleshooting)
- [Contributing](#contributing)
- [License](#license)

Features
--------

- **Splash Screen and Menu**: Provides an interactive user interface for bundling executables, making it user-friendly and accessible.
- **Executable Bundling**: Combines multiple executable files into a single target executable, which can streamline distribution and management.
- **User Input Handling**: Guides users through the process with prompts for necessary inputs, including file paths and target names.
- **Comprehensive Logging**: Detailed logs for different stages of the program, including informational, warning, and error messages.

How It Works
------------

INJECTOR combines multiple executable files into a single executable bundle using the following process:

1. **User Interaction**: The program starts with a splash screen and interactive menu, allowing users to specify the target file name and executable paths.
2. **Executable Collection**: The program reads the contents of each specified executable file.
3. **Bundle Creation**: A temporary Go source file is generated, containing code to handle the embedded executables.
4. **Compilation**: The temporary Go source file is compiled into a single executable file that includes the bundled content.
5. **Appending Data**: The contents of the specified executables are appended to the compiled bundle.
6. **Execution**: When the final bundle is run, it extracts and executes the embedded executables sequentially.

Installation
------------

1. **Clone the Repository**:

```bash
    git clone <repository_url>
    cd injector
```

2. **Install Dependencies**: Ensure you have Go installed on your system. Install necessary dependencies using:

```bash
    go mod tidy
```

3. **Build the Program** (optional): If you prefer to build the program instead of running it directly:

```bash
 go build -o injector main.go
```

`go build -o injector main.go

Usage
-----

### Running the Program

To run INJECTOR, use the following command:

```bash
 go run main.go
```

go run main.go

or, if you built the executable:

```bash
 ./injector
```

./injector

### Interacting with the Menu

1. **Splash Screen**: Displays initial information and options.
2. **Menu Options**: Choose to create a new bundle or exit the program.
3. **Prompt for Target and Executables**: Provide the target file name and paths to the executable files to be combined.

### Arguments and Inputs

When prompted, you will need to provide:

1. **Target File Name**: The name of the final bundled executable (e.g., `mybundle.exe`).
2. **Executable Paths**: Paths to each of the executable files to include in the bundle.

Technical Breakdown of Bundling Process
---------------------------------------

### Overview

The INJECTOR program bundles multiple executable files into a single target executable. This process involves several key steps, including reading executable contents, generating a Go source file to handle these contents, compiling the source file into an executable, and appending the original executables to the compiled binary.

Here's a detailed technical breakdown:

### 1\. **Reading Executable Contents**

The first step involves reading the binary data from the specified executable files. This is done using the `ioutil.ReadFile` function in Go, which reads the entire content of each executable file into a byte slice.

```go
func readFileContents(name string) []byte {
 contents, err := ioutil.ReadFile(name)
 if err != nil {
  logFatal(errors.New("Failed to read file " + name + ": " + err.Error()))
 }
 return contents
}
```

### 2\. **Generating the Temporary Go Source File**

A temporary Go source file is generated to create a program that can handle the embedded executables. This source file contains code that reads the bundled executable contents and executes them.

#### Temporary Source File Generation

The source file is created with a unique name to avoid conflicts. The generated code includes the following key parts:

- **Reading Embedded Data**: The program reads the contents of the bundled file (which includes the original executables) using `ioutil.ReadFile`.

- **Parsing and Executing**: The data is parsed using a predefined delimiter (`"h$lL@ w@rLd"`). Each segment is treated as a separate executable, and each segment is executed in sequence using `github.com/amenzhinsky/go-memexec`.

```go
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
 bytes, _:= ioutil.ReadFile(os.Args[0])
 executables := strings.Split(string(bytes), "h$lL@ w@rLd")
 for i := 2; i < len(executables); i++ {
  exe,_ := memexec.New([]byte(executables[i]))
  defer exe.Close()
  cmd := exe.Command()
  res, _ := cmd.Output()
  fmt.Print(string(res))
 }
}
`)
}
```

### 3\. **Compiling the Temporary Source File**

The generated Go source file is compiled into an executable using the `go build` command. This compiled executable is a standalone file capable of handling the embedded data.

```go

func createBundle(name string, data [][]byte) {
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

 removeFile(tempFileName)
}
```

### 4\. **Appending Executable Data**

After the compilation, the binary data of the original executables is appended to the newly created bundle executable. This is done by opening the bundle file in append mode and writing the binary data of each executable.

```go
func appendToBundle(name string, data [][]byte) {
 logInfo("Appending executable data to the bundle...")

 bundleFile := openFile(name)
 defer bundleFile.Close()

 for _, content := range data {
  writeToFile(name, bundleFile, content)
 }
}
```

### 5\. **Running the Bundle**

When the bundled executable is run, it performs the following steps:

- **Reads Its Own Binary**: The program reads its own binary content using `ioutil.ReadFile`.

- **Extracts Embedded Executables**: It uses the predefined delimiter to split the binary content into separate sections, each representing an embedded executable.

- **Executes Each Embedded Executable**: For each embedded executable, a new process is created using the `memexec` package, which provides an in-memory execution environment. This allows the program to run the embedded executables without needing to extract them to the filesystem.

```go
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
```

### Summary of Technical Aspects

- **Delimiter Usage**: A unique delimiter (`"h$lL@ w@rLd"`) is used to separate executable sections within the final bundle. This delimiter is essential for identifying and extracting each embedded executable.

- **In-Memory Execution**: The use of `memexec` allows for executing embedded binaries directly from memory, which can be advantageous for avoiding file system modifications.

- **Modular Design**: The bundling process is modular, involving separate functions for reading files, generating source code, compiling, and appending data. This design enhances maintainability and readability.

#### `go-memexec` Overview

The `go-memexec` package is used for running executables directly from memory, without writing them to the file system. This is particularly useful for scenarios where you want to execute embedded binaries within a larger application, often used in software bundling or self-extracting applications.

#### Key Functions and Usage

1. **Creating an Executable from Bytes**

    The `memexec.New` function is used to create a new in-memory executable from a byte slice containing the binary data of an executable file. This function takes a byte slice and returns a `*memexec.Exec` instance, which represents the in-memory executable.

```go
    exe, err := memexec.New([]byte(executableData))
    if err != nil {
        logFatal(errors.New("Failed to create in-memory executable: " + err.Error()))
    }
```

2. **Executing the In-Memory Executable**

    Once the in-memory executable is created, you can use the `Command` method to obtain an `exec.Cmd` instance. This allows you to execute the in-memory binary using standard Go `exec` commands. You can run the executable and capture its output.

```go
    cmd := exe.Command()
    output, err := cmd.Output()
    if err != nil {
        logFatal(errors.New("Failed to execute in-memory executable: " + err.Error()))
    }
    fmt.Print(string(output))
```

3. **Closing the In-Memory Executable**

    It is important to close the in-memory executable once you are done with it to free up resources. This is done using the `Close` method.

```go
    defer exe.Close()
```

#### Example in the Bundling Program

In the context of the INJECTOR program, `go-memexec` is used to execute each embedded executable from within the final bundled executable. The final bundle reads its own binary data, splits it based on a delimiter, and executes each segment as if it were a separate executable.

Here's a brief example of how `go-memexec` is used in the `generateBundleProgram` function:

```go
`func generateBundleProgram() []byte {
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
```

In this example, `memexec.New` creates an in-memory executable from each extracted segment of the bundled file. Each executable is then run, and its output is captured and displayed.

This approach allows for a clean and efficient method of bundling and executing multiple binaries, leveraging in-memory execution to avoid file system artifacts.

File Structure and Bundling
---------------------------

1. **Temporary Source File**: A Go source file is generated to manage embedded executables. This file is created with a unique name to prevent conflicts.
2. **Bundle Creation**: The temporary source file is compiled into the final target executable.
3. **Data Append**: Executable contents are appended to the final bundle, making it a self-contained file.

Security Implications
---------------------

### Potential Uses in Cybersecurity

- **Educational Purposes**: Understanding file bundling and executable manipulation is valuable for learning about software distribution and management.
- **Red Team Exercises**: In controlled environments, bundling can simulate the distribution of multiple tools or payloads.
- **Software Packaging**: Bundling multiple executables into a single file can streamline distribution and deployment in secure environments.

### Risks and Considerations

- **Malicious Use**: Bundling executables can be exploited by malicious actors to combine malware or unwanted software, disguising their presence and making detection more challenging.
- **Integrity and Trust**: Ensuring the integrity and authenticity of bundled executables is crucial. Unauthorized modifications can compromise security.
- **Anti-Virus Detection**: Security software might flag bundled executables if they exhibit unusual behaviors or if the content is known to be malicious.

### Best Practices

- **Use in Controlled Environments**: Conduct testing and educational exercises in isolated and controlled environments to prevent unintended consequences.
- **Verify Executables**: Ensure all executables included in the bundle are from trusted sources.
- **Educate Users**: Provide clear instructions and warnings to users about the potential risks and how to handle bundled executables securely.

Testing and Validation
----------------------

1. **Testing Bundling**:

    - Create test executables with known outputs.
    - Use INJECTOR to bundle these executables.
    - Run the resulting bundle to verify that all embedded executables execute correctly.
2. **Validation**:

    - Confirm all specified file paths are accurate.
    - Check that the final bundled executable operates as expected.

Troubleshooting
---------------

1. **Compilation Errors**: Ensure there are no syntax errors in the temporary Go source file. Check for missing dependencies.
2. **File Not Found**: Verify that all specified executable paths are correct and that files exist.
3. **Execution Issues**: Confirm that the bundled executable runs correctly and that all embedded files are intact.

Contributing
------------

Contributions are welcome. Please submit issues or pull requests with clear descriptions of changes or improvements. Ensure your code follows the existing style and includes tests where applicable.

License
-------

This project is licensed under the MIT License. See the LICENSE file for details.
