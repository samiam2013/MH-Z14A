package main

import (
	"time"

	"github.com/sirupsen/logrus"
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
		logrus.Fatal(err)
	}
	defer stream.Close()

	for {
		// send the 'gas concentration' command to get the current reading
		if _, err := stream.Write([]byte{0xFF, 0x01, 0x86, 0x00, 0x00, 0x00, 0x00, 0x00, 0x79}); err != nil {
			logrus.WithError(err).Warn("Failed writing command.")
			continue
		}
		response := make([]byte, 9)
		if _, err = stream.Read(response); err != nil {
			logrus.WithError(err).Warn("Failed reading response.")
			continue
		}

		if !checksumValidate(response) {
			logrus.Warn("Checksum validation error.")
			continue
		}

		ppm := int(response[2])*256 + int(response[3])
		logrus.Info("Reported concentration:", ppm, "ppm.")

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
