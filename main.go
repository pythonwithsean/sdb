package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("sdb> ")
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println()
			return
		}
		cmd := strings.TrimSpace(line)
		if cmd == "exit" || cmd == "quit" {
			return
		}
		if cmd == "" {
			continue
		}
		fmt.Println("unknown command:", cmd)
	}
}
