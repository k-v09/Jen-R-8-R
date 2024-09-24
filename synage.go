package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	sampleRate   = 44100
	duration     = 2
	fundamental  = 440
	numChannels  = 1
	bitDepth     = 16
	numHarmonics = 32
)

type Harmonic struct {
	Frequency float64
	Amplitude float64
	Waveform  string
}

var (
	octave    = 4
	relations = map[string]int{
		"c":  -9,
		"c#": -8,
		"d":  -7,
		"d#": -6,
		"e":  -5,
		"f":  -4,
		"f#": -3,
		"g":  -2,
		"g#": -1,
		"a":  0,
		"a#": 1,
		"b":  2,
	}
	stdKeys = map[string]string{
		"z": "c",
		"s": "c#",
		"x": "d",
		"d": "d#",
		"c": "e",
		"v": "f",
		"g": "f#",
		"b": "g",
		"h": "g#",
		"n": "a",
		"j": "a#",
		"m": "b",
	}
)

func generateWaveFile(harmonics []Harmonic, fileName string) {
	file, err := os.Create("generated/" + fileName + ".wav")
	if err != nil {
		log.Fatalf("Error creating file: %v", err)
	}
	defer file.Close()
	writeWAVHeader(file, sampleRate, numChannels, bitDepth)
	generateHarmonicWave(file, sampleRate, duration, harmonics)
}

func writeWAVHeader(file *os.File, sampleRate, numChannels, bitDepth int) {
	file.WriteString("RIFF")
	binary.Write(file, binary.LittleEndian, int32(36+duration*sampleRate*numChannels*bitDepth/8))
	file.WriteString("WAVE")

	file.WriteString("fmt ")
	binary.Write(file, binary.LittleEndian, int32(16))
	binary.Write(file, binary.LittleEndian, int16(1))
	binary.Write(file, binary.LittleEndian, int16(numChannels))
	binary.Write(file, binary.LittleEndian, int32(sampleRate))
	binary.Write(file, binary.LittleEndian, int32(sampleRate*numChannels*bitDepth/8))
	binary.Write(file, binary.LittleEndian, int16(numChannels*bitDepth/8))
	binary.Write(file, binary.LittleEndian, int16(bitDepth))

	file.WriteString("data")
	binary.Write(file, binary.LittleEndian, int32(duration*sampleRate*numChannels*bitDepth/8))
}

func sinCalc(x float64, freq float64, amp float64) float64 {
	return amp * math.Sin(2*math.Pi*freq*x)
}
func squareCalc(x float64, freq float64, amp float64) float64 {
	t := math.Sin(2 * math.Pi * freq * x)
	return amp * t / math.Abs(t)
}
func triangleCalc(x float64, freq float64, amp float64) float64 {
	t := math.Sin(2 * math.Pi * freq * x)
	return 2 * amp * math.Asin(t) / math.Pi
}
func sawCalc(x float64, freq float64, amp float64) float64 {
	t := 2 * amp * math.Asin(math.Sin(math.Pi*freq*x)) / math.Pi
	d := (2 * freq * math.Cos(math.Pi*freq*t)) / (math.Sqrt(1 - math.Pow(math.Sin(math.Pi*freq*t), 2)))
	return amp * t * d / math.Abs(d)
}

func generateHarmonicWave(file *os.File, sampleRate, duration int, harmonics []Harmonic) {
	for i := 0; i < sampleRate*duration; i++ {
		t := float64(i) / float64(sampleRate)
		sample := 0.0
		for _, h := range harmonics {
			switch h.Waveform {
			case "sine":
				sample += sinCalc(t, h.Frequency, h.Amplitude)
			case "square":
				sample += squareCalc(t, h.Frequency, h.Amplitude)
			case "triangle":
				sample += triangleCalc(t, h.Frequency, h.Amplitude)
			case "saw":
				sample += sawCalc(t, h.Frequency, h.Amplitude)
			}
		}
		sample /= float64(len(harmonics))
		intSample := int16(sample * 32767)
		binary.Write(file, binary.LittleEndian, intSample)
	}
	fmt.Println("Wave file generated successfully.")
}

func calculateFrequency(dist int, oct int) float64 {
	distance := float64(dist+(oct-4)*13) / 12
	return 440 * (math.Pow(2.0, distance))
}

func pipeListener() {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error getting current directory: %v", err)
	}
	pythonScriptPath := filepath.Join(dir, "revised.py")
	if _, err := os.Stat(pythonScriptPath); os.IsNotExist(err) {
		log.Fatalf("Python script does not exist at path: %s", pythonScriptPath)
	}
	cmd := exec.Command("python3", "-u", pythonScriptPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		log.Fatalf("Error starting Python script: %v", err)
	}
	fmt.Println("Python script started. Listening for key events...")

	pipePath := "/tmp/pipe_frequency"
	for {
		pipe, err := os.Open(pipePath)
		if err != nil {
			log.Fatalf("Error opening pipe: %v", err)
		}
		reader := bufio.NewReader(pipe)
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				pipe.Close()
				break
			}

			line = strings.TrimSpace(line)
			if strings.HasPrefix(line, "p:") {
				key := strings.TrimPrefix(line, "p:")
				if k, ok := stdKeys[key]; ok {
					freq := calculateFrequency(relations[k], octave)
					fmt.Printf("Playing note: %s at frequency %.2f Hz\n", k, freq)
				}
			} else if strings.HasPrefix(line, "r:") {
				key := strings.TrimPrefix(line, "r:")
				fmt.Printf("Key released: %s\n", key)
			}

			if line == "q" {
				fmt.Println("Exiting pipe program...")
				pipe.Close()
				os.Exit(0)
			}
		}
	}
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	harmonics := make([]Harmonic, numHarmonics)
	for i := range harmonics {
		harmonics[i] = Harmonic{Frequency: fundamental * float64(i+1), Amplitude: 0, Waveform: "sine"} // Default to sine wave
	}

	fmt.Println("Enter 'exit' to quit, 'generate [filename]' to create a wave file, 'listen' to start listening, or tune harmonics like '1 square 88':")

	for {
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "exit" {
			break
		} else if strings.HasPrefix(input, "generate") {
			parts := strings.Split(input, " ")
			fileName := "harmonic_wave.wav"
			if len(parts) > 1 {
				fileName = parts[1]
			}
			generateWaveFile(harmonics, fileName)
		} else if input == "listen" {
			pipeListener()
		} else {
			parts := strings.Split(input, " ")
			if len(parts) == 3 {
				harmonicIndex, err := strconv.Atoi(parts[0])
				if err != nil || harmonicIndex < 1 || harmonicIndex > len(harmonics) {
					fmt.Println("Invalid harmonic index.")
					continue
				}

				waveform := parts[1]
				if waveform != "sine" && waveform != "square" && waveform != "triangle" && waveform != "saw" {
					fmt.Println("Invalid waveform. Choose from 'sine', 'square', 'triangle', 'saw'.")
					continue
				}

				amplitude, err := strconv.ParseFloat(parts[2], 64)
				if err != nil || amplitude < 0 || amplitude > 100 {
					fmt.Println("Invalid amplitude. Must be between 0 and 100.")
					continue
				}

				harmonics[harmonicIndex-1].Waveform = waveform
				harmonics[harmonicIndex-1].Amplitude = amplitude / 100
				fmt.Printf("Harmonic %d set to %s wave with amplitude %.2f\n", harmonicIndex, waveform, amplitude)
			} else {
				fmt.Println("Invalid input. Use format: [harmonic] [waveform] [amplitude]")
			}
		}
	}
}
