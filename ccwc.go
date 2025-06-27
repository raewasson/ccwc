package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/urfave/cli/v3"
)

const (
	maxFileSize   = 100 * 1024 * 1024 // 100MB limit
	maxLineLength = 1024 * 1024       // 1MB line limit
	maxBufferSize = 64 * 1024         // 64KB scanner buffer
)

func main() {
	cmd := &cli.Command{
		Name:  "ccwc",
		Usage: "custom word, line, character, and byte count",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "c",
				Value: false,
				Usage: "The number of bytes in each input file is written to the standard output.",
			},
			&cli.BoolFlag{
				Name:  "l",
				Value: false,
				Usage: "The number of lines in each input file is written to the standard output.",
			},
			&cli.BoolFlag{
				Name:  "w",
				Value: false,
				Usage: "The number of words in each input file is written to the standard output.",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			var reader io.Reader
			var filename string

			if cmd.Args().Len() == 0 {
				reader = os.Stdin
				filename = ""
			} else {
				filename = cmd.Args().Get(0)
				fileHandle, err := os.Open(filename)
				if err != nil {
					return fmt.Errorf("failed to open file: %w", err)
				}
				defer fileHandle.Close()
				reader = fileHandle
			}

			linesCount := 0
			wordsCount := 0
			bytesCount := 0

			scanner := bufio.NewScanner(reader)
			buf := make([]byte, maxBufferSize)
			scanner.Buffer(buf, maxLineLength)

			for scanner.Scan() {
				line := scanner.Text()
				linesCount++
				wordsCount += len(strings.Fields(line))
				bytesCount += len(line) + 1 // +1 for the newline character
			}
			if err := scanner.Err(); err != nil {
				return fmt.Errorf("failed to read input: %w", err)
			}

			showLines := cmd.Bool("l")
			showWords := cmd.Bool("w")
			showBytes := cmd.Bool("c")
			if !showLines && !showWords && !showBytes {
				showLines, showWords, showBytes = true, true, true
			}

			if showLines {
				fmt.Printf("%d\t", linesCount)
			}
			if showWords {
				fmt.Printf("%d\t", wordsCount)
			}
			if showBytes {
				fmt.Printf("%d\t", bytesCount)
			}
			if filename != "" {
				fmt.Printf("%s", filename)
			}
			fmt.Println()
			return nil
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
