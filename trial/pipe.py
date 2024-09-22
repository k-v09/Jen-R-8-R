import os
from pynput import keyboard
import termios
import sys

def disable_echo():
    fd = sys.stdin.fileno()
    old_settings = termios.tcgetattr(fd)
    new_settings = old_settings[:]
    new_settings[3] &= ~(termios.ECHO)  
    termios.tcsetattr(fd, termios.TCSADRAIN, new_settings)

def enable_echo():
    fd = sys.stdin.fileno()
    old_settings = termios.tcgetattr(fd)
    old_settings[3] |= termios.ECHO  
    termios.tcsetattr(fd, termios.TCSADRAIN, old_settings)

pipe_path = "/tmp/pipe_frequency"
if not os.path.exists(pipe_path):
    os.mkfifo(pipe_path)

def on_press(key):
    try:
        with open(pipe_path, 'w') as pipe:
            pipe.write(f"p:{key.char}\n")  
            pipe.flush()
        if key.char == 'q':
            return False
    except AttributeError:
        pass  

def on_release(key):
    try:
        with open(pipe_path, 'w') as pipe:
            pipe.write(f"r:{key.char}\n")  
            pipe.flush()
        if key.char == 'q':
            return False
    except AttributeError:
        pass  

if __name__ == "__main__":
    disable_echo()  
    try:
        with keyboard.Listener(on_press=on_press, on_release=on_release) as listener:
            print("Listener on (key presses/releases will not show in terminal)")
            listener.join()
    except KeyboardInterrupt:
        pass
    finally:
        enable_echo()  
