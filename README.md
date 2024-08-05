Executable Binder
=================

Table of Contents
-----------------

- [Overview](#overview)
- [Features](#features)
- [Requirements](#requirements)
- [Installation](#installation)
- [Usage](#usage)
- [How It Works](#how-it-works)
- [Challenges and Learnings](#challenges-and-learnings)
- [Security Considerations](#security-considerations)
- [Conclusion](#conclusion)

Overview
--------

The Executable Binder is a tool designed to merge two executable programs into a single executable. This tool is implemented in Go and allows users to embed two executables into a new executable file, which can then execute both embedded programs concurrently.

Features
--------

- **Embed Executables**: Combine two executable files into a single executable.
- **Custom Output Path**: Specify a custom output path for the new executable.
- **Concurrency**: Execute the embedded programs concurrently.
- **Logging**: Detailed logging to help with debugging and tracking execution.

Requirements
------------

- **Go 1.19 or later**: The tool is implemented in Go and requires Go 1.19 or later.
- **MinGW-w64**: For Windows users, MinGW-w64 is required to ensure 64-bit mode compatibility.

Installation
------------

1. **Install Go**: Download and install Go 1.19 or later from the [official Go website](https://golang.org/dl/).
2. **Install MinGW-w64**: For Windows users, download and install MinGW-w64 from SourceForge.
3. **Set Up Environment Variables**:
    - Add the Go `bin` directory to your `PATH`.
    - Add the MinGW-w64 `bin` directory to your `PATH`.

Usage
-----

### Building the Program

To build the `injector.exe` executable, run the following command:

```sh
go build -o injector.exe main.go
```

### Merging Two Programs

To merge two executable programs, use the following command:

```sh
injector.exe <file1> <file2> [outputPath]
```

Where:

- `file1`: Path to the first executable.
- `file2`: Path to the second executable.
- `outputPath` (optional): Custom output path for the new executable. Defaults to `newprogram.exe`.

### Running the Merged Program

To run the merged executable, simply execute the output file created by the binder:

```sh
newprogram.exe
```

How It Works
------------

1. **Argument Validation**: The program validates the provided arguments to ensure two input files are specified.
2. **Reading Executables**: The specified executable files are read and stored as byte slices.
3. **Embedding**: The byte slices are written to temporary files (`important1.exe` and `important2.exe`).
4. **Building New Executable**: A new executable is built using the Go `exec.Command` to run `go build`.
5. **Executing Embedded Programs**: When the new executable is run, the embedded executables are extracted to temporary files and executed concurrently.

Challenges and Learnings
------------------------

### Initial Challenges

We initially attempted to create this tool using Python. However, embedding binary data and executing it securely and efficiently proved challenging in Python. We faced issues with file handling, concurrency, and cross-platform compatibility.

### Transition to Go

Switching to Go provided several benefits:

- **Static Typing**: Improved type safety and reduced runtime errors.
- **Concurrency**: Go's goroutines made it easier to handle concurrent execution.
- **Efficiency**: Go's performance and compilation to native binaries improved execution speed and ease of distribution.

### Technical Challenges

- **64-bit Compatibility**: Ensuring the compiler supported 64-bit mode required installing MinGW-w64 and configuring environment variables.
- **Logging and Error Handling**: Implementing comprehensive logging and detailed error messages helped in debugging and improving the user experience.

Security Considerations
-----------------------

### Dangers of File Binding

Binding multiple executable files can pose several security risks:

- **Malware Distribution**: Malicious actors can use file binders to distribute malware hidden within legitimate software.
- **Unintended Behavior**: Executing multiple programs together can lead to unintended interactions and behavior.
- **Integrity Risks**: Ensuring the integrity and authenticity of the embedded executables is crucial.

### Learnings

- **Awareness**: Understanding the potential dangers of file binding emphasized the need for secure development practices.
- **User Education**: Educating users on the risks associated with running merged executables is important for promoting safe usage.

Conclusion
----------

The Executable Binder tool allows users to merge two executables into a single file, providing custom output paths, concurrent execution, and detailed logging. While developing this tool, we learned valuable lessons about security, file handling, and concurrency. We chose Go over Python for its performance, concurrency model, and static typing, which significantly improved the tool's reliability and efficiency.

