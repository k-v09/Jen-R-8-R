package main

import (
	"encoding/binary"
	"math"
	"os"
)

const (
	sampleRate  = 44100
	duration    = 2
	fundamental = 440
	numChannels = 1
	bitDepth    = 16
)

type Harmonic struct {
	Frequency float64
	Amplitude float64
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
	file, err := os.Create("out/harmonic.wav")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	writeWAVHeader(file, sampleRate, numChannels, bitDepth)

	harmonics := []Harmonic{
		{Frequency: fundamental, Amplitude: 0.5},
		{Frequency: fundamental * 2, Amplitude: 0.3},
		{Frequency: fundamental * 3, Amplitude: 0.15},
		{Frequency: fundamental * 4, Amplitude: 0.05},
	}

	generateHarmonicWave(file, sampleRate, duration, harmonics)
}
