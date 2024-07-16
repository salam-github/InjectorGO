# Embedded Executables Program

This program is designed to embed two executable files (`important1.exe` and `important2.exe`) into a main Go program. It provides functionality to update these embedded executables with new versions and then run them.

## How the Program Works

The program operates in two modes:

1. **Injector Mode**:
   - When the program is compiled and run as `injector.exe`, it reads two new executable files provided as command-line arguments.
   - It replaces the embedded `important1.exe` and `important2.exe` with these new files.
   - It recompiles the program with the updated embedded executables and saves it as `newprogram.exe`.

2. **Normal Mode**:
   - When the program is run with any name other than `injector.exe`, it extracts the embedded executables to temporary files and executes them.

## How to Embed Files

### Step-by-Step Process

1. **Ensure Files Exist**:
   - Make sure `important1.exe` and `important2.exe` exist in your directory. If not, create or copy them there.

2. **Compile the Injector**:
   - Compile the Go code as `injector.exe`:
     ```sh
     go build -o injector.exe main.go
     ```

3. **Run the Injector**:
   - Use `injector.exe` to embed the new executables:
     ```sh
     .\injector.exe .\important1.exe .\important2.exe
     ```
   - This updates the embedded executables in the program and recompiles it, saving the new version as `newprogram.exe`.

### Code Explanation

The key parts of the program are:

1. **Embedding Executables Using `go:embed`**:
   - The `//go:embed` directive is used to embed the executables into the Go program at compile time.
     ```go
     //go:embed important1.exe
     var embeddedProgram1 []byte

     //go:embed important2.exe
     var embeddedProgram2 []byte
     ```

2. **Injector Mode**:
   - When the program is named `injector.exe`, it reads new executable files and replaces the embedded executables.
     ```go
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
     ```

3. **Normal Mode**:
   - When the program is not named `injector.exe`, it writes the embedded executables to temporary files and executes them.
     ```go
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
     ```

4. **Argument Checking Function**:
   - This function ensures that exactly two command-line arguments are provided.
     ```go
     func check_args() (string, string, error) {
         args := os.Args[1:]
         if len(args) != 2 {
             return "", "", fmt.Errorf("Usage: %s <file1> <file2>", os.Args[0])
         }
         return args[0], args[1], nil
     }
     ```

## How to Run the Program

1. **Run the New Program**:
   - Execute the newly compiled program (`newprogram.exe`) to run the embedded executables:
     ```sh
     .\newprogram.exe
     ```

When executed, `newprogram.exe` will extract the embedded executables to temporary files and run them.

## Example Commands

### Initial Setup

1. **Initialize Go Module (if not already initialized)**:
    ```sh
    go mod init test
    ```

2. **Compile `injector.exe`**:
    ```sh
    go build -o injector.exe main.go
    ```

3. **Embed New Executables**:
    ```sh
    .\injector.exe .\important1.exe .\important2.exe
    ```

4. **Run the Program**:
    ```sh
    .\newprogram.exe
    ```

## Notes

- Ensure `main.go` is in the same directory when running these commands.
- Ensure `important1.exe` and `important2.exe` files exist and are accessible.
- The `injector.exe` creates a new executable named `newprogram.exe` with the embedded executables.
- The program must not be named `injector.exe` when running in normal mode to execute the embedded executables.

By following these steps, you can embed and run executables directly within your main Go program.
