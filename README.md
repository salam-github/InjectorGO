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
