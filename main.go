package main

import (
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
		// send the 'gas concentration' command to get the current reading
		_, err := stream.Write([]byte{0xFF, 0x01, 0x86, 0x00, 0x00, 0x00, 0x00, 0x00, 0x79})
		if err != nil {
			log.Fatal(err)
		}
		b := make([]byte, 9)
		_, err = stream.Read(b)
		if err != nil {
			log.Fatal(err)
		}

		if !checksumValidate(b) {
			log.Print("checksum error")
			continue
		}

		ppm := int(b[2])*256 + int(b[3])
		log.Println("reported concentration: ", ppm, "ppm")

		time.Sleep(1 * time.Second)
	}

}

func checksumValidate(b []byte) bool {
	var sum byte
	if len(b) != 9 {
		return false
	}
	for i := 1; i < len(b)-1; i++ {
		sum += b[i]
	}
	sum = 0xFF - sum
	return (sum + 1) == b[8]
}
