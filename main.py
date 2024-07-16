import os
import sys
import subprocess
import tempfile

def check_args():
    if len(sys.argv) != 3:
        print(f"Usage: {sys.argv[0]} <file1> <file2>")
        sys.exit(1)
    return sys.argv[1], sys.argv[2]

def embed_and_create_executable(file1, file2):
    try:
        with open(file1, 'rb') as f:
            embedded_program1 = f.read()
        with open(file2, 'rb') as f:
            embedded_program2 = f.read()

        with open('important1.exe', 'wb') as f:
            f.write(embedded_program1)
        with open('important2.exe', 'wb') as f:
            f.write(embedded_program2)

        # Simulate creating a new executable
        result = subprocess.run(["go", "build", "-o", "newprogram.exe", "main.go"], capture_output=True, text=True)
        if result.returncode != 0:
            print(result.stderr)
            sys.exit(1)
        print("Successfully created important1.exe, important2.exe, and newprogram.exe.")
        
        return embedded_program1, embedded_program2
    except Exception as e:
        print(e)
        sys.exit(1)

def main():
    if len(sys.argv) == 3:
        file1, file2 = check_args()
        embedded_program1, embedded_program2 = embed_and_create_executable(file1, file2)
    else:
        embedded_program1 = open('important1.exe', 'rb').read()
        embedded_program2 = open('important2.exe', 'rb').read()

    try:
        tmpfile1_path = None
        tmpfile2_path = None
        
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
    except Exception as e:
        print(e)
        sys.exit(1)
    finally:
        try:
            if tmpfile1_path and os.path.exists(tmpfile1_path):
                os.remove(tmpfile1_path)
            if tmpfile2_path and os.path.exists(tmpfile2_path):
                os.remove(tmpfile2_path)
        except PermissionError as pe:
            print(f"Could not remove temporary file: {pe}")

if __name__ == "__main__":
    main()
