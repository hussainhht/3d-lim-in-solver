package lemin

import (
	"log"
	"os"
	"bufio"
)

func WriteFile(
	data string,
) {
	file, err := os.Create("result.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	_, err = file.WriteString(data)
	if err != nil {
		log.Println("Error writing to file:", err)
		return
	}
}

func TakeData() string {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Println("Error reading from file:", err)
		return ""
	}
	defer file.Close()

	var data string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		data += scanner.Text() + "\n"
	}

	return string(data+"\n\n\n\n\n\n\n\n\n\n\n\n")
}

