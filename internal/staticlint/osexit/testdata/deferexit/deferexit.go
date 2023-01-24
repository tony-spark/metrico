package main

import (
	"fmt"
	"os"
)

func main() {
	defer func() {
		fmt.Println("exiting...")
		os.Exit(-1) // want "os.Exit in main"
	}()
}
