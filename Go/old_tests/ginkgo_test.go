package main_test

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"regexp"
	"strings"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestMain(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Main Suite")
}

var _ = Describe("CLI", func() {
	It("Test Meta Commands", func() {
		cmd := exec.Command("make", "run")

		// Create pipes for stdin and stdout
		stdin, err := cmd.StdinPipe()
		Expect(err).To(BeNil())
		stdout, err := cmd.StdoutPipe()
		Expect(err).To(BeNil())

		// Start the command
		err = cmd.Start()
		Expect(err).To(BeNil())

		// Write .help to stdin
		_, err = io.WriteString(stdin, ".help\n")
		Expect(err).To(BeNil())

		// Close stdin to signal we're done writing
		stdin.Close()

		// Read the output
		var output strings.Builder
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			line := scanner.Text()
			//fmt.Println("Output:", line) // Print each line for debugging
			output.WriteString(line + "\n")
			if strings.Contains(line, "Helping...") {
				break // Exit after we've seen the expected output
			}
		}

		// Wait for the command to finish
		err = cmd.Wait()
		Expect(err).To(BeNil())

		// Check the output
		Expect(output.String()).To(ContainSubstring("Helping..."))
	})

	It("Test Simple Insert and Select", func() {
		cmd := exec.Command("make", "run")

		// Create pipes for stdin and stdout
		stdin, err := cmd.StdinPipe()
		Expect(err).To(BeNil())
		stdout, err := cmd.StdoutPipe()
		Expect(err).To(BeNil())

		// Start the command
		err = cmd.Start()
		Expect(err).To(BeNil())

		// Write first command to stdin
		_, err = io.WriteString(stdin, "insert 1 braden braden@gmail.com\n")
		Expect(err).To(BeNil())

		// Write second command to stdin
		_, err = io.WriteString(stdin, "select\n")
		Expect(err).To(BeNil())

		// Close stdin to signal we're done writing
		stdin.Close()

		// Read the output
		var output strings.Builder
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			line := scanner.Text()
			fmt.Println(strings.TrimSpace(line))
			fmt.Println(strings.TrimSpace("(1 braden braden@gmail.com)\n"))
			matched, _ := regexp.MatchString("Goqlite> Goqlite> (1 braden braden@gmail.com)\n", line)
			fmt.Println(matched)
			// if regexp.MatchString("(1 braden braden@gmail.com)\n", line) {
			// 	output.WriteString(line + "\n")
			// 	break // Exit after we've seen the expected output
			// }
		}

		// Wait for the command to finish
		err = cmd.Wait()
		Expect(err).To(BeNil())

		// Check that the output contains the expected result
		outputStr := strings.TrimSpace(output.String())
		Expect(outputStr).To(ContainSubstring("(1 braden braden@gmail.com)"))

		// Optionally, you can also print the output for debugging
	})

})
