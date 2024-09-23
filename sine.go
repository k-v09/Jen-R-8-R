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
	"strings"

	"fyne.io/fyne/v2"
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
	numHarmonics = 32
	graphWidth   = 600
	graphHeight  = 200
	numPoints    = 200
)

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

func createGeneratorContainer() *fyne.Container {
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

	sb1 := container.NewHBox()
	sb2 := container.NewHBox()
	sb3 := container.NewHBox()
	sb4 := container.NewHBox()
	for i := range harmonics {
		vbox := container.NewVBox(labels[i], sliders[i])
		if i < 8 {
			sb1.Add(vbox)
		} else if i < 16 {
			sb2.Add(vbox)
		} else if i < 24 {
			sb3.Add(vbox)
		} else if i < 32 {
			sb4.Add(vbox)
		}
	}

	return (container.NewVBox(
		widget.NewLabel("Generator Mode"),
		sb1,
		sb2,
		sb3,
		sb4,
		generateButton,
	))
}

func createLiveContainer() *fyne.Container {
	psButton := widget.NewButton("Start Listener", func() {
		pipeListener()
	})
	return (container.NewVBox(
		widget.NewLabel("Live Mode"),
		psButton,
	))
}

func calculateFrequency(dist int, oct int) float64 {
	distance := float64(dist+(oct-4)*13) / 12
	return 440 * (math.Pow(2.0, float64(distance)))
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
				k := stdKeys[key]
				fmt.Printf("Playing note: %s at frequency %f\n", k, calculateFrequency(relations[k], octave))
				line = key
			} else if strings.HasPrefix(line, "r:") {
				key := strings.TrimPrefix(line, "r:")
				fmt.Printf("Key released: %s\n", key)
				line = key
			}

			if line == "q" {
				fmt.Println("Exiting pipe programs...")
				pipe.Close()
				break
			}
		}
	}
}

func main() {
	a := app.New()
	w := a.NewWindow("Harmonic Wave Generator")
	modeToggle := widget.NewCheck("Live Mode", func(checked bool) {})
	modeToggle.SetChecked(false)

	generatorContainer := createGeneratorContainer()
	liveContainer := createLiveContainer()
	contentCard := widget.NewCard("", "", generatorContainer)

	updateContent := func(checked bool) {
		if checked {
			contentCard.SetContent(liveContainer)
		} else {
			contentCard.SetContent(generatorContainer)
		}
	}

	modeToggle.OnChanged = updateContent

	content := container.NewVBox(
		modeToggle,
		contentCard,
	)
	w.SetContent(content)
	w.ShowAndRun()
}
