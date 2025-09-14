package lemin

import (
	"bufio"
	"fmt"
	"log"
	"os"
	// "os/exec"
)

// * Write a result to a file
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

// * Read result from file and print to stdout
func TypeResult() {
	file, err := os.Open("result.txt")
	if err != nil {
		log.Println("Error reading from file:", err)
		os.Exit(1)
	}
	defer file.Close()
	result := ""
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		result += scanner.Text() + "\n"
	}
	if err := scanner.Err(); err != nil {
		log.Println("Error reading from file:", err)
		os.Exit(1)
	}
	fmt.Print(result)
}

// * Take data from a input file
func TakeData(f_name string) string {
	// str, err := exec.Command("cat", f_name).Output()
	// if err != nil {
	// 	log.Println("Error getting current directory:", err)
	// 	return
	// }
	// fmt.Println(string(str))
	file, err := os.Open(f_name)
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
	fmt.Println(data)
	return string(data + "\n")
}
