package main

import (
	"bufio"
	"fmt"
	"log"

	"github.com/tarm/serial"
)

func main() {
	config := &serial.Config{
		Name:        "/dev/ttyAMA0",
		Baud:        9600,
		ReadTimeout: 1,
		Size:        8,
	}

	stream, err := serial.OpenPort(config)
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(stream)
	for scanner.Scan() {
		fmt.Println(scanner.Text()) // Println will add back the final '\n'
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

}
