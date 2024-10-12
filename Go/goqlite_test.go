package main

import (
	"bytes"
	"fmt"
	// "math/rand"
	"os/exec"
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
	var debugMode bool
	var debugInput string
	runCLI(reader, writer, debugMode, debugInput)

	// Step 4: Capture and verify output
	output := writer.String()
	if output != expectedOutput {
		t.Errorf("Unexpected output: got %v, want %v", output, expectedOutput)
	}
}

func TestSimpleInsertSelect(t *testing.T) {
	input := "insert 1 braden braden@gmail.com\nselect\n"
	expectedOutput := "(1 braden braden@gmail.com)\n"

	reader := strings.NewReader(input)
	writer := &bytes.Buffer{}

	var debugMode bool
	var debugInput string
	runCLI(reader, writer, debugMode, debugInput)

	output := writer.String()
	output = strings.TrimSpace(strings.ReplaceAll(output, "\x00", ""))
	expectedOutput = strings.TrimSpace(expectedOutput)
	fmt.Println(expectedOutput)
	fmt.Println(output)
	if output != expectedOutput {
		t.Errorf("Unexpected output: got %v, want %v", output, expectedOutput)
	}
	err := exec.Command("rm", "-r", "mydb.db").Run()
	if err != nil {
		t.Fatalf("Failed to remove database file: %v", err)
	}
}

// func TestTableFull(t *testing.T) {
// 	var inputBuilder strings.Builder
// 	expectedOutput := "Error: Table Full\n"
// 	//1401
// 	for i := 0; i < 1401; i++ {
// 		username := randomString(8)                           // Generate random 8-character username
// 		email := fmt.Sprintf("%s@gmail.com", randomString(6)) // Generate random 6-character email
// 		inputBuilder.WriteString(fmt.Sprintf("insert %d %s %s\n", i, username, email))
// 	}
//
// 	reader := strings.NewReader(inputBuilder.String())
// 	writer := &bytes.Buffer{}
//
// 	var debugMode bool
// 	var debugInput string
// 	runCLI(reader, writer, debugMode, debugInput)
//
// 	output := writer.String()
// 	output = strings.TrimSpace(strings.ReplaceAll(output, "\x00", ""))
// 	expectedOutput = strings.TrimSpace(expectedOutput)
// 	if output != expectedOutput {
// 		t.Errorf("Unexpected output: got = %v; want = %v;", output, expectedOutput)
// 	}
// 	err := exec.Command("rm", "-r", "mydb.db").Run()
// 	if err != nil {
// 		t.Fatalf("Failed to remove database file: %v", err)
// 	}
//
// }
//
// func TestNegativeId(t *testing.T) {
// 	input := "insert -1 braden braden@gmail.com\n"
// 	expectedOutput := "Error: Cannot insert a negative id\n"
//
// 	reader := strings.NewReader(input)
// 	writer := &bytes.Buffer{}
//
// 	var debugMode bool
// 	var debugInput string
// 	runCLI(reader, writer, debugMode, debugInput)
//
// 	output := writer.String()
// 	output = strings.TrimSpace(strings.ReplaceAll(output, "\x00", ""))
// 	expectedOutput = strings.TrimSpace(expectedOutput)
// 	if output != expectedOutput {
// 		t.Errorf("Unexpected output: got = %v; want = %v;", output, expectedOutput)
// 	}
// }
//
// func TestPersistance(t *testing.T) {
// 	input := "insert 1 braden braden@gmail.com\n"
// 	reader := strings.NewReader(input)
// 	writer := &bytes.Buffer{}
// 	var debugMode bool
// 	var debugInput string
// 	runCLI(reader, writer, debugMode, debugInput)
//
// 	input = ".exit\n"
// 	reader = strings.NewReader(input)
// 	writer = &bytes.Buffer{}
// 	runCLI(reader, writer, debugMode, debugInput)
// 	output_test := writer.String()
// 	t.Log(output_test)
// 	t.Log("test")
//
// 	input1 := "select\n"
// 	reader1 := strings.NewReader(input1)
// 	writer1 := &bytes.Buffer{}
// 	runCLI(reader1, writer1, debugMode, debugInput)
//
// 	expectedOutput := "(1 braden braden@gmail.com)\n"
//
// 	output := writer1.String()
// 	output = strings.TrimSpace(strings.ReplaceAll(output, "\x00", ""))
// 	expectedOutput = strings.TrimSpace(expectedOutput)
// 	if output != expectedOutput {
// 		t.Errorf("Unexpected output: got %v, want %v", output, expectedOutput)
// 	}
// 	err := exec.Command("rm", "-r", "mydb.db").Run()
// 	if err != nil {
// 		t.Fatalf("Failed to remove database file: %v", err)
// 	}
// }
//
// // RandomString generates a random string of a given length.
// func randomString(n int) string {
// 	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
// 	result := make([]byte, n)
// 	for i := range result {
// 		result[i] = letters[rand.Intn(len(letters))]
// 	}
// 	return string(result)
// }
//
// func TestComplexInsertSelect(t *testing.T) {
// 	var inputBuilder strings.Builder
// 	var expectedOutputBuilder strings.Builder
//
// 	// Insert 500 rows of data with random usernames and emails
// 	for i := 1; i <= 1000; i++ {
// 		username := randomString(8)                           // Generate random 8-character username
// 		email := fmt.Sprintf("%s@gmail.com", randomString(6)) // Generate random 6-character email
// 		inputBuilder.WriteString(fmt.Sprintf("insert %d %s %s\n", i, username, email))
// 		expectedOutputBuilder.WriteString(fmt.Sprintf("(%d %s %s)\n", i, username, email))
// 	}
// 	t.Log(inputBuilder.String())
// 	input := inputBuilder.String()
// 	expectedOutput := expectedOutputBuilder.String()
// 	reader := strings.NewReader(input)
// 	writer := &bytes.Buffer{}
// 	var debugMode bool
// 	var debugInput string
// 	runCLI(reader, writer, debugMode, debugInput)
//
// 	input = ".exit\n"
// 	reader = strings.NewReader(input)
// 	writer = &bytes.Buffer{}
// 	runCLI(reader, writer, debugMode, debugInput)
//
// 	input1 := "select\n"
// 	reader1 := strings.NewReader(input1)
// 	writer1 := &bytes.Buffer{}
// 	runCLI(reader1, writer1, debugMode, debugInput)
//
// 	output := writer1.String()
// 	output = strings.TrimSpace(strings.ReplaceAll(output, "\x00", ""))
// 	expectedOutput = strings.TrimSpace(expectedOutput)
// 	if output != expectedOutput {
// 		t.Errorf("Unexpected output: got %v, want %v", "", "")
// 	}
//
// 	// Clean up by removing the database file
// 	err := exec.Command("rm", "-r", "mydb.db").Run()
// 	if err != nil {
// 		t.Fatalf("Failed to remove database file: %v", err)
// 	}
// }
