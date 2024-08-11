package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

const (
	sampleRate   = 44100
	duration     = 2
	fundamental  = 440
	numChannels  = 1
	bitDepth     = 16
	numHarmonics = 8
)

type Harmonic struct {
	Frequency float64
	Amplitude float64
}

func generateWaveFile(harmonics []Harmonic) {
	file, err := os.Create("out/alt_harm2.wav")
	if err != nil {
		panic(err)
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

func generateHarmonicWave(file *os.File, sampleRate, duration int, harmonics []Harmonic) {
	for i := 0; i < sampleRate*duration; i++ {
		t := float64(i) / float64(sampleRate)
		sample := 0.0
		for _, h := range harmonics {
			sample += h.Amplitude * math.Sin(2*math.Pi*h.Frequency*t)
		}
		sample /= float64(len(harmonics))
		intSample := int16(sample * 32767)
		binary.Write(file, binary.LittleEndian, intSample)
	}
}

func main() {
	harmonics := make([]Harmonic, numHarmonics)
	for i := 0; i < numHarmonics; i++ {
		harmonics[i] = Harmonic{Frequency: fundamental * float64(i+1)}
	}

	for {
		fmt.Println("\nCurrent harmonic amplitudes:")
		for i, h := range harmonics {
			fmt.Printf("%d. Frequency: %.2f Hz, Amplitude: %.2f\n", i+1, h.Frequency, h.Amplitude)
		}

		fmt.Println("\nEnter new amplitudes (0-1) for each harmonic, separated by spaces:")
		fmt.Println("Or type 'q' to quit, 'g' to generate the wave with current settings")

		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "q" {
			fmt.Println("Exiting program.")
			return
		} else if input == "g" {
			generateWaveFile(harmonics)
			fmt.Println("Wave file generated with current settings.")
		} else {
			amplitudes := strings.Split(input, " ")
			if len(amplitudes) != numHarmonics {
				fmt.Printf("Please enter exactly %d amplitudes.\n", numHarmonics)
				continue
			}

			for i, amp := range amplitudes {
				value, err := strconv.ParseFloat(amp, 64)
				if err != nil || value < 0 || value > 1 {
					fmt.Printf("Invalid amplitude for harmonic %d. Please enter a number between 0 and 1.\n", i+1)
					continue
				}
				harmonics[i].Amplitude = value
			}
		}
	}
}
