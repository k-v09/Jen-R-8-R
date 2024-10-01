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
	Waveform  int
}

var (
	octave          = 4
	currentHarmonic = 0 // Default to first harmonic
	relations       = map[string]int{
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
		s2 := 0.0
		for _, h := range harmonics {
			if h.Waveform <= 25 {
				sample += sinCalc(t, h.Frequency, h.Amplitude)
				s2 += squareCalc(t, h.Frequency, h.Amplitude)
				sample = float64(1-h.Waveform/25)*sample + float64(h.Waveform/25)*s2
			} else if h.Waveform <= 50 {
				sample += squareCalc(t, h.Frequency, h.Amplitude)
				s2 += triangleCalc(t, h.Frequency, h.Amplitude)
				sample = float64(1-h.Waveform/50)*sample + float64(h.Waveform/50)*s2
			} else if h.Waveform <= 75 {
				sample += triangleCalc(t, h.Frequency, h.Amplitude)
				s2 = sawCalc(t, h.Frequency, h.Amplitude)
				sample = float64(1-h.Waveform/75)*sample + float64(h.Waveform/75)*s2
			} else if h.Waveform <= 100 {
				sample += sawCalc(t, h.Frequency, h.Amplitude)
				s2 = sinCalc(t, h.Frequency, h.Amplitude)
				sample = float64(1-h.Waveform/100)*sample + float64(h.Waveform/100)*s2
			} else {
				sample += sinCalc(t, h.Frequency, h.Amplitude)
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

func pipeListener(harmonics []Harmonic) {
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

			if strings.HasPrefix(line, "sel:") {
				value := strings.TrimPrefix(line, "sel:")
				harmonicIndex, err := strconv.Atoi(value)
				if err == nil && harmonicIndex > 0 && harmonicIndex <= numHarmonics {
					currentHarmonic = harmonicIndex - 1
					fmt.Printf("Selector value: Harmonic %d selected\n", harmonicIndex)
				} else {
					fmt.Println("Invalid harmonic selection")
				}
			} else if strings.HasPrefix(line, "pot:") {
				value := strings.TrimPrefix(line, "pot:")
				amplitude, err := strconv.ParseFloat(value, 64)
				if err == nil && amplitude >= 0 && amplitude <= 100 {
					harmonics[currentHarmonic].Amplitude = amplitude / 100
					fmt.Printf("Potentiometer value: Amplitude of harmonic %d set to %.2f%%\n", currentHarmonic+1, amplitude)
				} else {
					fmt.Println("Invalid amplitude value")
				}
			} else if strings.HasPrefix(line, "z:") {
				key := strings.TrimPrefix(line, "z:")
				if k, ok := stdKeys[key]; ok {
					freq := calculateFrequency(relations[k], octave)
					fmt.Printf("Playing note: %s at frequency %.2f Hz\n", k, freq)
				}
			} else if strings.HasPrefix(line, "r:") {
				key := strings.TrimPrefix(line, "r:")
				fmt.Printf("Key released: %s\n", key)
			} else if strings.HasPrefix(line, "w:") {
				value := strings.TrimPrefix(line, "w:")
				fmt.Printf("Waveform of %d changed to value %s\n", currentHarmonic+1, value)
			} else if strings.HasPrefix(line, "generate_wave:") {
				fmt.Println("Received command to generate .wav file: ", line)
				generateWaveFile(harmonics, "weird")
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
	harmonics := make([]Harmonic, numHarmonics)
	for i := range harmonics {
		harmonics[i] = Harmonic{
			Frequency: fundamental * float64(i+1),
			Amplitude: 0,
			Waveform:  0, // Default to sine wave
		}
	}
	go pipeListener(harmonics)
	select {}
}
