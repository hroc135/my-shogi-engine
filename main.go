package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const (
	USIPonder = "USI_Ponder"
	USIHash   = "USI_Hash"
)

func run() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		split := strings.Fields(line)
		if len(split) == 0 {
			continue
		}

		command := split[0]
		switch command {
		case "usi":
			fmt.Println("id name hroc135")
			fmt.Println("id author hroc135")
			fmt.Printf("option name %s type check default true\n", USIPonder)
			fmt.Printf("option name %s type spin default 256\n", USIHash)
			fmt.Println("usiok")
			os.Stdout.Sync()
		case "isready":
			fmt.Println("readyok")
			os.Stdout.Sync()
		case "usinewgame", "setoption", "position", "stop", "ponderhit", "gameover":
		case "go":
			fmt.Println("bestmove resign")
			os.Stdout.Sync()
		case "quit":
			return
		default:
			fmt.Printf("info string Unsupported command: %s\n", command)
			os.Stdout.Sync()
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error during scan: %v\n", err)
	}
}

func main() {
	run()
}
