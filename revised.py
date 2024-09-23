import os
import sys
from pynput import keyboard

pipe = None
def setup_pipe():
    global pipe
    pipe_path = "/tmp/pipe_frequency"
    if not os.path.exists(pipe_path):
        os.mkfifo(pipe_path)
    pipe = open(pipe_path, 'w')

def cleanup_pipe():
    global pipe
    if pipe:
        pipe.close()

def write_to_pipe(message):
    global pipe
    if pipe:
        pipe.write(message)
        pipe.flush()

def on_press(key):
    try:
        char = key.char
    except AttributeError:
        char = str(key)
    write_to_pipe(f"p:{char}\n")
    return char != 'q'

def on_release(key):
    try:
        char = key.char
    except AttributeError:
        char = str(key)
    write_to_pipe(f"r:{char}\n")
    return char != 'q'

if __name__ == "__main__":
    try:
        setup_pipe()
        print("Keyboard listener starting...")
        with keyboard.Listener(on_press=on_press, on_release=on_release) as listener:
            listener.join()
    except Exception as e:
        print(f"An error occurred: {e}", file=sys.stderr)
        sys.exit(1)
    finally:
        cleanup_pipe()
        print("Keyboard listener stopped.")