package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"time"

	"github.com/hajimehoshi/oto"
)

const (
	sql = 44100
)

func generateSineWave(frequency float64, duration time.Duration) []byte {
	length := int(float64(sql) * duration.Seconds())
	data := make([]byte, length*2) // 16-bit audio (2 bytes per sample)

	fmt.Printf("Generating sine wave for frequency: %.2f Hz\n", frequency)

	for i := 0; i < length; i++ {
		t := float64(i) / float64(sql)
		sample := int16(math.Sin(2.0*math.Pi*frequency*t) * 32767)
		data[2*i] = byte(sample)
		data[2*i+1] = byte(sample >> 8)
	}

	fmt.Println("Sine wave generation complete.")
	return data
}

func q() {
	pipePath := "/tmp/pipe_frequency"

	fmt.Printf("Opening pipe: %s\n", pipePath)
	pipe, err := os.Open(pipePath)
	if err != nil {
		fmt.Printf("Error opening pipe: %v\n", err)
		return
	}
	defer pipe.Close()
	fmt.Println("Pipe opened successfully. Waiting for frequency data...")

	// Initialize audio player
	playerCtx, err := oto.NewContext(sql, 1, 2, 4096)
	if err != nil {
		fmt.Printf("Error initializing audio player: %v\n", err)
		return
	}
	player := playerCtx.NewPlayer()

	scanner := bufio.NewScanner(pipe)
	for {
		if scanner.Scan() {
			freqStr := scanner.Text()
			fmt.Printf("Received data from pipe: %s\n", freqStr)

			// Parse frequency
			frequency, err := strconv.ParseFloat(freqStr, 64)
			if err != nil {
				fmt.Printf("Invalid frequency data: %s\n", freqStr)
				continue
			}

			fmt.Printf("Playing frequency: %.2f Hz\n", frequency)

			// Generate wave for 1 second
			data := generateSineWave(frequency, 1*time.Second)
			fmt.Println("Playing sound...")
			player.Write(data)
			fmt.Println("Sound played.")
		} else {
			fmt.Println("No data received, waiting for next input...")
		}
		if err := scanner.Err(); err != nil {
			fmt.Printf("Error reading from pipe: %v\n", err)
			break
		}
	}
}
