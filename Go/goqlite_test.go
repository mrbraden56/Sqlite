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
//5. Insert statement persists

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
	input := "insert 1 braden braden@gmail.com\nselect"
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

	for i := 0; i < 1401; i++ {
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

// func TestPersistance(t *testing.T) {
// 	input := "insert 1 braden braden@gmail.com\n"
// 	reader := strings.NewReader(input)
// 	writer := &bytes.Buffer{}
// 	runCLI(reader, writer)
//
// 	input = ".exit\n"
// 	reader = strings.NewReader(input)
// 	writer = &bytes.Buffer{}
// 	runCLI(reader, writer)
// 	output_test := writer.String()
// 	t.Log(output_test)
// 	t.Log("test")
//
// 	input1 := "select 1 braden braden@gmail.com\n"
// 	reader1 := strings.NewReader(input1)
// 	writer1 := &bytes.Buffer{}
// 	runCLI(reader1, writer1)
//
// 	expectedOutput := "(1 braden braden@gmail.com)\n"
//
// 	output := writer1.String()
// 	if output != expectedOutput {
// 		t.Errorf("Persistance test failed: %v", output)
// 	}
// }
