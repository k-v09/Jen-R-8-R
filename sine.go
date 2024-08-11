package main

import (
	"encoding/binary"
	"math"
	"os"
)

const (
	sampleRate  = 44100
	duration    = 2
	frequency   = 440
	amplitude   = 0.5
	numChannels = 1
	bitDepth    = 16
)

func main() {
	file, err := os.Create("out/sine_wave.wav")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	writeWAVHeader(file, sampleRate, numChannels, bitDepth)
	generateSineWave(file, sampleRate, duration, frequency, amplitude)
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

func generateSineWave(file *os.File, sampleRate, duration int, frequency, amplitude float64) {
	for i := 0; i < sampleRate*duration; i++ {
		t := float64(i) / float64(sampleRate)
		sample := int16(amplitude * math.Sin(2*math.Pi*frequency*t) * 32767)
		binary.Write(file, binary.LittleEndian, sample)
	}
}
