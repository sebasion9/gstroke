package main

import (
	"fmt"
	"gstroke/decoder"
	"os"
)

func main() {
	args := os.Args
	if len(os.Args) != 2 {
		fmt.Println("usage: gstroke <file.jpeg>")
		os.Exit(1)
	}

	path := args[1]

	src, err := os.ReadFile(path)
	if err != nil {
		fmt.Printf("[ERR] failed to read file at %s\n", path)
		fmt.Printf("[ERR] %s\n", err.Error())
		os.Exit(1)
	}

	decoder := decoder.NewDecoder(src)
	if err := decoder.Decode(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
