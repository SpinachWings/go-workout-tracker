package initializers

import (
	"log"
	"os"
)

func InitLogs() {
	file, err := os.OpenFile("../../logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	log.SetOutput(file)
}
