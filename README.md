# Synage

## Overview
This project is a command-line tool written in Go that allows users to generate harmonic waves using various waveforms (sine, square, triangle, and sawtooth) and customize their amplitude and frequency. The tool outputs a `.wav` file that can be played using any audio software. Additionally, the program includes a listener mode that interacts with a Python script for real-time input, making it highly customizable for audio synthesis tasks.

## Features
- **Generate Harmonic Waves**: Generate up to 32 harmonics, each with configurable waveform and amplitude.
- **Supported Waveforms**: Choose between sine, square, triangle, and sawtooth waves for each harmonic.
- **Customizable Amplitudes**: Set the amplitude of each harmonic individually (0 to 100 scale).
- **Frequency Calculation**: Automatically calculate the frequency for different notes based on standard music notation.
- **Listener Mode**: Pipe listener mode interacts with a Python script to detect keyboard inputs for playing notes in real time.
- **WAV File Output**: Outputs the generated sound to a `.wav` file, ready for playback or further processing.

## How to Use

1. **Running the Program**:
   ```bash
   go run synage.go
   ```

2. **Available Commands**:
    - `exit`: Quit the program.
    - `generate [filename]`: Generates a wave file with the current harmonic settings. If no filename is provided, it defaults to harmonic_wave.wav.
    - `listen`: Starts the listener mode to capture real-time input from the Python script.

3. **Tuning Harmonics**: You can tune individual harmonics using the following format:
    ```
    [harmonic] [waveform] [amplitude]
    ```
    Example:
    ```
    1 square 88
    ```
    This sets the first harmonic (the fundamental) to a square wave with 88% amplitude.

4. **Waveform Options**:
    - `sine`
    - `square`
    - `triangle`
    - `saw`

## Example Usage
After starting the program, you can use commands like:
```
1 sine 100
2 square 80
generate super_awesome_wave
```
This will generate a .wav file with the first harmonic as a sine wave at full amplitude and the second harmonic as a square wave at 80% amplitude.

## File Structure
- `synage.go`: Main Go file for the application.
- `revised.py`: Python script used in listener mode for real-time input handling.
- `/generated`: Directory where .wav files are saved, but my sounds are too cool so you're not allowed to see them!! ...hear them that is!!

## Future Plans
- **Plugin Development**: Plan to convert this tool into a plugin for popular DAWs (Digital Audio Workstations), allowing more flexibility and easier integration for musicians and audio engineers.
- **MIDI Support**: Add MIDI input support for real-time harmonic control.
- **GUI Integration**: Develop a graphical interface for non-technical users to easily adjust harmonics and generate sound waves.
- **Expanded Waveform Options**: Incorporate additional waveforms and custom user-defined waveforms for more sound design possibilities.
- **Enhanced Real-Time Capabilities**: Improve real-time interaction and expand listener mode features.

## Requirements
- Go 1.16+
- Python 3.x (for listener mode)

## Installation
1. Clone the repository:
    ```bash
    git clone https://github.com/k-v09/harmonic-wave-generator.git
    ```

2. Navigate into the project directory:
    ```bash
    cd harmonic-wave-generator
    ```

3. Run the Go file:
    ```bash
    go run synage.go
    ```

For listener mode, ensure Python is installed and the `revised.py` script is in the same directory.

## Contributing
Feel free to fork the repository and submit pull requests for new features, bug fixes, or enhancements. I have no doubt that your code is better than mine.