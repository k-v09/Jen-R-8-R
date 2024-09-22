package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
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
