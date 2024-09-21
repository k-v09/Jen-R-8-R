import os
from pynput import keyboard

note_frequencies = {
    'n': 440.0,  # A4
    'z': 261.63, # C4
}

# Create pipe
pipe_path = "/tmp/pipe_frequency"
if not os.path.exists(pipe_path):
    os.mkfifo(pipe_path)

# handle keys
def on_press(key):
    try:
        if key.char in note_frequencies:
            frequency = note_frequencies[key.char]
            with open(pipe_path, 'w') as pipe:
                pipe.write(str(frequency))
                pipe.flush()
    except AttributeError:
        pass  # Special keys (like ctrl, shift) are ignored

# Start listening for key presses
with keyboard.Listener(on_press=on_press) as listener:
    print("Press 'n' for A4 (440Hz), 'z' for C4 (261.63Hz)...")
    listener.join()
