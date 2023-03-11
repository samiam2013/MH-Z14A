package main

import (
	"fmt"
	"log"
	"time"

	"github.com/tarm/serial"
)

func main() {
	config := &serial.Config{
		Name:        "/dev/serial0",
		Baud:        9600,
		ReadTimeout: 1,
		Size:        8,
	}

	stream, err := serial.OpenPort(config)
	if err != nil {
		log.Fatal(err)
	}
	defer stream.Close()

	for {
		_, err := stream.Write([]byte{0xFF, 0x01, 0x86, 0x00, 0x00, 0x00, 0x00, 0x00, 0x79})
		if err != nil {
			log.Fatal(err)
		}
		b := make([]byte, 9)
		_, err = stream.Read(b)
		if err != nil {
			log.Fatal(err)
		}

		upperB := b[2]
		lowerB := b[3]
		ppm := int(upperB)*256 + int(lowerB)
		fmt.Println("reported concentration: ", ppm, "ppm")

		time.Sleep(1 * time.Second)
	}

}
