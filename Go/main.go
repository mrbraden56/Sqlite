package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

var PREPARE_UNRECOGNIZED_STATEMENT = errors.New("Prepare Unrecognized Statement")
var PREPARE_PARSING_ERROR = errors.New("Prepare Parsing Error")
var PREPARE_NEGATIVE_ID = errors.New("Error: Cannot insert a negative id")
var EXECUTE_TABLE_FULL = errors.New("Error: Table Full")

func runCLI(reader io.Reader, writer io.Writer) {
	scanner := bufio.NewScanner(reader)

	for {
		fmt.Print("Goqlite> ")
		if !scanner.Scan() {
			if err := scanner.Err(); err != nil {
				fmt.Fprintln(writer, "Error reading input:", err)
			}
			return
		}

		input := strings.TrimSpace(scanner.Text())

		// Meta Commands
		if len(input) > 0 && input[0] == '.' {
			switch input {
			case ".exit":
				fmt.Fprintln(writer, "Exiting...")
				return
			case ".help":
				fmt.Fprintln(writer, "Helping...")
				continue
			default:
				fmt.Fprintln(writer, "Unrecognized Command")
				continue
			}
		}

		var statement Statement
		err := prepare_statement(input, &statement)
		if err != nil {
			switch {
			case errors.Is(err, PREPARE_UNRECOGNIZED_STATEMENT):
				fmt.Fprintln(writer, PREPARE_UNRECOGNIZED_STATEMENT)
			case errors.Is(err, PREPARE_PARSING_ERROR):
				fmt.Fprintln(writer, PREPARE_PARSING_ERROR)
			case errors.Is(err, PREPARE_NEGATIVE_ID):
				fmt.Fprintln(writer, PREPARE_NEGATIVE_ID)

			}
			continue
		}

		err = execute_statement(&statement, writer)
		if err != nil {
			switch {
			case errors.Is(err, EXECUTE_TABLE_FULL):
				fmt.Fprintln(writer, EXECUTE_TABLE_FULL)
			default:
				fmt.Fprintln(writer, "Unrecognized Error")
			}
		}
	}
}

func main() {
	runCLI(os.Stdin, os.Stdout)
}
