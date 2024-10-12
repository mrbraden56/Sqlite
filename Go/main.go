package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

var PREPARE_UNRECOGNIZED_STATEMENT = errors.New("Prepare Unrecognized Statement")
var PREPARE_PARSING_ERROR = errors.New("Prepare Parsing Error")
var PREPARE_NEGATIVE_ID = errors.New("Error: Cannot insert a negative id")
var EXECUTE_TABLE_FULL = errors.New("Error: Table Full")
var PAGER_OPENING_ERROR = errors.New("Error: Opening Pager Failed")

func runCLI(reader io.Reader, writer io.Writer, debug bool, debugInput string) {
	var filename string
	if debug {
		filename = os.Args[3]
	} else {
		filename = os.Args[1]
	}
	var err error
	table, err = db_open(filename)
	if err != nil {
		switch {
		case errors.Is(err, PAGER_OPENING_ERROR):
			fmt.Fprintln(writer, PAGER_OPENING_ERROR)

		}
	}

	scanner := bufio.NewScanner(reader)

	for {
		if !debug {
			fmt.Print("Goqlite> ")
			if !scanner.Scan() {
				if err := scanner.Err(); err != nil {
					fmt.Fprintln(writer, "Error reading input:", err)
				}
				return
			}
		}

		var input string
		if debug {
			input = debugInput
		} else {

			input = strings.TrimSpace(scanner.Text())
		}

		// Meta Commands
		if len(input) > 0 && input[0] == '.' {
			switch input {
			case ".exit":
				fmt.Fprintln(writer, "Exiting...")
				db_close(table)
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
				fmt.Fprintln(writer, err)
			}
		}
		if debug {
			break
		}
	}
}

var debugMode bool
var debugInput string

func main() {

	flag.BoolVar(&debugMode, "d", false, "Run in debug mode")
	flag.StringVar(&debugInput, "input", "", "Input for debug mode")
	flag.Parse()
	if debugMode {
		fmt.Println("Running in debug mode")
		fmt.Println(debugInput)
		// Add your debug mode logic here
	} else {
		fmt.Println("Running in normal mode")
		// Add your normal mode logic here
	}
	runCLI(os.Stdin, os.Stdout, debugMode, debugInput)
}
