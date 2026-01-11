package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/hroc135/my-shogi-engine/domain/model"
)

const (
	USIPonder = "USI_Ponder"
	USIHash   = "USI_Hash"
)

func main() {
	position := model.Position{}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		split := strings.Fields(line)
		if len(split) == 0 {
			continue
		}

		command := split[0]
		switch model.USICommand(command) {
		case "d": // for debugging
			fmt.Print(position.BoardToString())
		case model.USICommandUSI:
			fmt.Println("id name hroc135")
			fmt.Println("id author hroc135")
			fmt.Printf("option name %s type check default true\n", USIPonder)
			fmt.Printf("option name %s type spin default 256\n", USIHash)
			fmt.Println("usiok")
			if err := os.Stdout.Sync(); err != nil {
				fmt.Println(err)
			}
		case model.USICommandIsReady:
			fmt.Println("readyok")
			if err := os.Stdout.Sync(); err != nil {
				fmt.Println(err)
			}
		case model.USICommandUSINewGame, model.USICommandSetOption, model.USICommandStop, model.USICommandPonderhit, model.USICommandGameover:
		case model.USICommandPosition:
			if len(split) < 2 {
				fmt.Println("invalid command")
				continue
			}
			// var nextIdx int
			switch model.USIPositionSubcommand(split[1]) {
			case model.USIPositionSubcommandStartpos:
				position = model.SetPosition(model.StartPosSFEN)
				// nextIdx = 2
			case model.USIPositionSubcommandSFEN:
				sfen := strings.Join(split[2:], " ")
				position = model.SetPosition(sfen)
				// nextIdx = 6
			default:
				fmt.Println("invalid command")
				continue
			}

			// todo: 指し手のパース
		case model.USICommandGo:
			fmt.Println("bestmove resign")
			if err := os.Stdout.Sync(); err != nil {
				fmt.Println(err)
			}
		case model.USICommandQuit:
			return
		case model.USICommandGenerateMove:
			ch := make(chan model.Move)
			go func() {
				position.Generate(ch)
				close(ch)
			}()
			for m := range ch {
				fmt.Printf("%s ", m.ToString())
			}
			fmt.Println()
		default:
			fmt.Printf("info string Unsupported command: %s\n", command)
			if err := os.Stdout.Sync(); err != nil {
				fmt.Println(err)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error during scan: %v\n", err)
	}
}
