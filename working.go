package main

import (
	"fmt"
	"math"
	"time"

	"github.com/hajimehoshi/oto"
	"github.com/nsf/termbox-go"
)

const (
	sq         = 44100
	bufferSize = 4096
	frequency  = 440 // A4 note
)

func makeSamps(freq int) []byte {
	samps := make([]byte, 4*sq)
	for i := 0; i < len(samps)/4; i++ {
		t := float64(i) / float64(sq)
		amplitude := int16(math.Sin(2*math.Pi*float64(freq)*t) * 32767)
		samps[4*i] = byte(amplitude)
		samps[4*i+1] = byte(amplitude >> 8)
		samps[4*i+2] = byte(amplitude)
		samps[4*i+3] = byte(amplitude >> 8)
	}
	return samps
}

func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	context, err := oto.NewContext(sq, 1, 2, bufferSize)
	if err != nil {
		panic(err)
	}
	defer context.Close()

	player := context.NewPlayer()
	defer player.Close()

	fmt.Println("Press and hold 'k' to play sound. Press 'q' to quit.")

	samples := makeSamps(fundamental)

	playing := false
	playChannel := make(chan bool)

	go func() {
		for play := range playChannel {
			if play {
				player.Write(samples)
			}
		}
	}()

	eventQueue := make(chan termbox.Event)
	go func() {
		for {
			eventQueue <- termbox.PollEvent()
		}
	}()

	for {
		select {
		case ev := <-eventQueue:
			if ev.Type == termbox.EventKey {
				switch ev.Key {
				case termbox.KeyEsc:
					fmt.Println("Quitting...")
					close(playChannel)
					return
				default:
					switch ev.Ch {
					case 'q':
						fmt.Println("Quitting...")
						close(playChannel)
						return
					case 'k':
						if ev.Type == termbox.EventKey {
							if !playing {
								playing = true
								fmt.Println("Key 'k' pressed, starting sound...")
								go func() {
									for playing {
										playChannel <- true
									}
								}()
							}
						}
					case 'p':
						if playing {
							playing = false
							fmt.Println("Key 'p' pressed, stopping sound...")
						}
					}
				}
			}
		default:
			time.Sleep(10 * time.Millisecond)
		}
	}
}

func adaptedLive() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	context, err := oto.NewContext(sq, 1, 2, bufferSize)
	if err != nil {
		panic(err)
	}
	defer context.Close()

	player := context.NewPlayer()
	defer player.Close()

	fmt.Println("Press and hold 'k' to play sound. Press 'q' to quit.")

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

	eventQueue := make(chan termbox.Event)
	go func() {
		for {
			eventQueue <- termbox.PollEvent()
		}
	}()

	for {
		select {
		case ev := <-eventQueue:
			if ev.Type == termbox.EventKey {
				switch ev.Key {
				case termbox.KeyEsc:
					fmt.Println("Quitting...")
					close(playChannel)
					return
				default:
					switch ev.Ch {
					case 'q':
						fmt.Println("Quitting...")
						close(playChannel)
						return
					case 'k':
						if ev.Type == termbox.EventKey {
							if !playing {
								playing = true
								fmt.Println("Key 'k' pressed, starting sound...")
								go func() {
									for playing {
										playChannel <- true
									}
								}()
							}
						}
					case 'p':
						if playing {
							playing = false
							fmt.Println("Key 'p' pressed, stopping sound...")
						}
					}
				}
			}
		default:
			time.Sleep(10 * time.Millisecond)
		}
	}
}
