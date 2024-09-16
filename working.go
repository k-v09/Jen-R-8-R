package main

import (
	"fmt"
	"math"

	"github.com/eiannone/keyboard"
	"github.com/hajimehoshi/oto"
)

const (
	sq         = 44100
	bufferSize = 4096
	frequency  = 440 // A4 note
)

func main() {
	context, err := oto.NewContext(sq, 1, 2, bufferSize)
	if err != nil {
		panic(err)
	}
	defer context.Close()

	player := context.NewPlayer()
	defer player.Close()

	// Open keyboard
	if err := keyboard.Open(); err != nil {
		panic(err)
	}
	defer keyboard.Close()

	fmt.Println("Press and hold 'k' to play sound. Press 'q' to quit.")

	// Generate one second of sample data
	samples := make([]byte, 4*sq)
	for i := 0; i < len(samples)/4; i++ {
		t := float64(i) / float64(sq)
		amplitude := int16(math.Sin(2*math.Pi*frequency*t) * 32767)
		samples[4*i] = byte(amplitude)
		samples[4*i+1] = byte(amplitude >> 8)
		samples[4*i+2] = byte(amplitude)
		samples[4*i+3] = byte(amplitude >> 8)
	}

	playing := false
	playChannel := make(chan bool)

	go func() {
		for play := range playChannel {
			if play {
				player.Write(samples)
			}
		}
	}()

	for {
		char, key, err := keyboard.GetKey()
		if err != nil {
			panic(err)
		}

		if key == keyboard.KeyEsc || char == 'q' {
			fmt.Println("Quitting...")
			close(playChannel)
			return
		}

		if char == 'k' {
			if !playing {
				playing = true
				go func() {
					for playing {
						playChannel <- true
					}
				}()
			}
		} else {
			playing = false
		}
	}
}
