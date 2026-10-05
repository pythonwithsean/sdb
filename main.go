package main

import (
	"bufio"
	"fmt"
	"os"
	queryprocessor "sdb/QueryProcessor"
	storageengine "sdb/StorageEngine"
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
		stmt, err := queryprocessor.Parse(cmd)
		if err != nil {
			fmt.Println("error:", err)
			continue
		}
		err = storageengine.ExecuteStatement(stmt)
		if err != nil {
			fmt.Println("error:", err)
			continue
		}
		fmt.Printf("ok\n")

	}
}
