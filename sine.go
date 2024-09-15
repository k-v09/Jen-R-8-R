package main

import (
	"encoding/binary"
	"fmt"
	"math"
	"os"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

const (
	sampleRate   = 44100
	duration     = 2
	fundamental  = 440
	numChannels  = 1
	bitDepth     = 16
	numHarmonics = 16
	graphWidth   = 600
	graphHeight  = 200
	numPoints    = 200
)

type Harmonic struct {
	Frequency float64
	Amplitude float64
}

type WaveformGraph struct {
	widget.BaseWidget
}

func generateWaveFile(harmonics []Harmonic) {
	file, err := os.Create("out/harmonic_wave.wav")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	writeWAVHeader(file, sampleRate, numChannels, bitDepth)
	generateHarmonicWave(file, sampleRate, duration, harmonics)
	fmt.Println("Wave file generated successfully.")
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
	a := app.New()
	w := a.NewWindow("Harmonic Wave Generator")

	harmonics := make([]Harmonic, numHarmonics)
	sliders := make([]*widget.Slider, numHarmonics)
	labels := make([]*widget.Label, numHarmonics)

	for i := range harmonics {
		harmonics[i] = Harmonic{Frequency: fundamental * float64(i+1), Amplitude: 0}
		sliders[i] = widget.NewSlider(0, 100)
		sliders[i].Step = 1
		labels[i] = widget.NewLabel(fmt.Sprintf("Harmonic %d (%.0f Hz): 0.00", i+1, harmonics[i].Frequency))

		index := i
		sliders[i].OnChanged = func(value float64) {
			amplitude := value / 100
			harmonics[index].Amplitude = amplitude
			labels[index].SetText(fmt.Sprintf("Harmonic %d (%.0f Hz): %.2f", index+1, harmonics[index].Frequency, amplitude))
		}

		sliders[i].Orientation = widget.Vertical
	}

	generateButton := widget.NewButton("Generate Wave", func() {
		generateWaveFile(harmonics)
	})

	sliderBox := container.NewHBox()
	for i := range harmonics {
		vbox := container.NewVBox(labels[i], sliders[i])
		sliderBox.Add(vbox)
	}

	content := container.NewVBox(
		sliderBox,
		generateButton,
	)

	w.SetContent(content)
	w.ShowAndRun()
}
