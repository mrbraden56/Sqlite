package main

import (
	"bytes"
	"strings"
	"testing"
)

//NOTE: Tests
//1. Testing meta command outputs
//2. Testing single insert and select
//3. Testing that our db throws error when table is full
//4. Testing that db throws error when a user enter a negative id

func TestMetaCommands(t *testing.T) {
	// Step 1: Set up input and expected output
	input := ".help\n.dexit\n.exit\n" // Input simulates two commands: `.dexit` and `.exit`
	expectedOutput := "Helping...\nUnrecognized Command\nExiting...\n"

	// Step 2: Create a reader and writer
	reader := strings.NewReader(input)
	writer := &bytes.Buffer{}

	// Step 3: Run the CLI with the mocked input/output
	runCLI(reader, writer)

	// Step 4: Capture and verify output
	output := writer.String()
	if output != expectedOutput {
		t.Errorf("Unexpected output: got %v, want %v", output, expectedOutput)
	}
}

func TestSimpleInsertSelect(t *testing.T) {
	input := "insert 1 braden braden@gmail.com\nselect 1 braden braden@gmail.com"
	expectedOutput := "(1 braden braden@gmail.com)\n"

	reader := strings.NewReader(input)
	writer := &bytes.Buffer{}

	runCLI(reader, writer)

	output := writer.String()
	if output != expectedOutput {
		t.Errorf("Unexpected output: got %v, want %v", output, expectedOutput)
	}
}

func TestTableFull(t *testing.T) {
	input := ""
	expectedOutput := "Error: Table Full\n"

	for i := 0; i < 1400; i++ {
		input = input + "insert 1 braden braden@gmail.com\n"
	}

	reader := strings.NewReader(input)
	writer := &bytes.Buffer{}

	runCLI(reader, writer)

	output := writer.String()
	if output != expectedOutput {
		t.Errorf("Unexpected output: got = %v; want = %v;", output, expectedOutput)
	}

}

func TestNegativeId(t *testing.T) {
	input := "insert -1 braden braden@gmail.com\n"
	expectedOutput := "Error: Cannot insert a negative id\n"

	reader := strings.NewReader(input)
	writer := &bytes.Buffer{}

	runCLI(reader, writer)

	output := writer.String()
	if output != expectedOutput {
		t.Errorf("Unexpected output: got = %v; want = %v;", output, expectedOutput)
	}

}
