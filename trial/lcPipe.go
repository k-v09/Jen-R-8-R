package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
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
	keys = map[string]string{
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

func main() {
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
				fmt.Printf("Key pressed: %s\n", key)
				fmt.Printf("Playing note: %s\n", keys[key])
				fmt.Printf("Distance from A: %s\n", string(relations[keys[key]]))
			} else if strings.HasPrefix(line, "r:") {
				key := strings.TrimPrefix(line, "r:")
				fmt.Printf("Key released: %s\n", key)
			}

			if line == "p:q" || line == "r:q" {
				fmt.Println("Exiting both programs...")
				os.Exit(0)
			}
		}
	}
}
