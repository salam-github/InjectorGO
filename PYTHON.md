# Embedded Executables Program

This program is designed to embed two executable files (`important1.exe` and `important2.exe`) and manage their execution using a Python script. It provides functionality to update these embedded executables with new versions and then run them.

## How the Program Works

The program operates in two main phases:

1. **Embedding Phase**:
   - The script reads two new executable files provided as command-line arguments.
   - It replaces the embedded `important1.exe` and `important2.exe` with these new files.
   - It creates new executable files (`important1.exe` and `important2.exe`) and updates them.

2. **Execution Phase**:
   - The script extracts the embedded executables to temporary files and executes them.

## How to Embed Files

### Step-by-Step Process

1. **Ensure Files Exist**:
   - Make sure `important1.exe` and `important2.exe` exist in your directory. If not, create or copy them there.

2. **Run the Script**:
   - Use the Python script to embed the new executables and create the updated files:
     ```sh
     python main.py important1.exe important2.exe
     ```
   - This updates the embedded executables in the script and creates new versions of `important1.exe` and `important2.exe`.

### Code Explanation

The key parts of the program are:

1. **Embedding Executables**:
   - The script reads the contents of `important1.exe` and `important2.exe` and writes them to new files.
     ```python
     with open(file1, 'rb') as f:
         embedded_program1 = f.read()
     with open(file2, 'rb') as f:
         embedded_program2 = f.read()

     with open('important1.exe', 'wb') as f:
         f.write(embedded_program1)
     with open('important2.exe', 'wb') as f:
         f.write(embedded_program2)
     ```

2. **Execution Phase**:
   - The script writes the embedded executables to temporary files and executes them.
     ```python
     with tempfile.NamedTemporaryFile(delete=False, suffix=".exe") as tmpfile1:
         tmpfile1_path = tmpfile1.name
         tmpfile1.write(embedded_program1)
         tmpfile1.close()

     with tempfile.NamedTemporaryFile(delete=False, suffix=".exe") as tmpfile2:
         tmpfile2_path = tmpfile2.name
         tmpfile2.write(embedded_program2)
         tmpfile2.close()

     os.chmod(tmpfile1_path, 0o755)
     os.chmod(tmpfile2_path, 0o755)

     result1 = subprocess.run([tmpfile1_path], capture_output=True, text=True)
     if result1.returncode != 0:
         print(result1.stderr)
         sys.exit(1)

     result2 = subprocess.run([tmpfile2_path], capture_output=True, text=True)
     if result2.returncode != 0:
         print(result2.stderr)
         sys.exit(1)
     ```

3. **Argument Checking Function**:
   - This function ensures that exactly two command-line arguments are provided.
     ```python
     def check_args():
         if len(sys.argv) != 3:
             print(f"Usage: {sys.argv[0]} <file1> <file2>")
             sys.exit(1)
         return sys.argv[1], sys.argv[2]
     ```

## How to Run the Program

1. **Run the Embedding and Creation Process**:
   - Execute the script with the paths to `important1.exe` and `important2.exe` to embed them and create new versions:
     ```sh
     python main.py important1.exe important2.exe
     ```

2. **Run the Program**:
   - Execute the script without arguments to run the embedded executables:
     ```sh
     python main.py
     ```

When executed, the script will extract the embedded executables to temporary files and run them.

## Example Commands

### Initial Setup

1. **Ensure `important1.exe` and `important2.exe` Exist**:
    - Place `important1.exe` and `important2.exe` in the same directory as the script.

2. **Embed New Executables and Create Updated Files**:
    ```sh
    python main.py important1.exe important2.exe
    ```

3. **Run the Embedded Executables**:
    ```sh
    python main.py
    ```

## Notes

- Ensure the Python script (`main.py`) is in the same directory when running these commands.
- Ensure `important1.exe` and `important2.exe` files exist and are accessible.
- The script updates and creates new versions of `important1.exe` and `important2.exe` with the embedded executables.

By following these steps, you can embed and run executables directly using the provided Python script.
