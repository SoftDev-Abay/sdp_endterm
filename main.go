package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	// Create a new scanner to read from standard input
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("Enter some text (write 'end' to end):")

	// Read input line by line
	for scanner.Scan() {
		text := scanner.Text() // Get the current line of text
		if text == "end" {
			break // Exit loop if an empty line is entered
		}
		fmt.Println("You entered:", text)
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error:", err)
	}
}
