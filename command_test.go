package egg

import (
	"bytes"
	"io"
	"log"
	"os"
	"sync"
	"testing"
)

func TestNewCommand(t *testing.T) {
	var command *Command
	var err error

	// NewCommand with no args
	command, err = NewCommand([]string{})

	if err == nil {
		t.Errorf("NewCommand with no args should of returned an error")
	}
	if command != nil {
		t.Errorf("NewCommand with no args return %v, it should have been nil", command)
	}

	// NewCommand with one arg
	command, err = NewCommand([]string{"ls"})

	if err != nil {
		t.Errorf("NewCommand with one arg shouldn't of returned an error: %v", err)
	}
	if command == nil {
		t.Errorf("NewCommand with one args should have returns a command but it did not")
	}

	// NewCommand with multiple args
	command, err = NewCommand([]string{"ls", "-al"})

	if err != nil {
		t.Errorf("NewCommand with one arg shouldn't of returned an error: %v", err)
	}
	if command == nil {
		t.Errorf("NewCommand with one args should have returns a command but it did not")
	}
}

func TestRun(t *testing.T) {
	var command *Command
	var err error
	var output string
	var expectedOutput string

	// Run with echo
	output = captureOutput(func() {
		command, err = NewCommand([]string{"echo", "whatup"})
		if err != nil {
			t.Errorf("Error creating command: %v", err)
		}

		err = command.Run()
		if err != nil {
			t.Errorf("Error running echo %v", err)
		}
	})

	expectedOutput = "whatup\n"
	if output != expectedOutput {
		t.Errorf("Command output should have been %#v, but was %#v", expectedOutput, output)
	}

	// Run with command that doesn't exist
	output = captureOutput(func() {
		command, err = NewCommand([]string{"idontexist"})
		if err != nil {
			t.Errorf("Error creating command: %v", err)
		}

		err = command.Run()
		if err == nil {
			t.Errorf("Non existing command should have raise error")
		}
	})

	expectedOutput = ""
	if output != expectedOutput {
		t.Errorf("Command output should have been %#v, but was %#v", expectedOutput, output)
	}

	// Run called twice
	output = captureOutput(func() {
		command, err = NewCommand([]string{"echo"})
		if err != nil {
			t.Errorf("Error creating command: %v", err)
		}

		err = command.Run()
		if err != nil {
			t.Errorf("Error running command: %v", err)
		}

		err = command.Run()
		if err == nil {
			t.Errorf("No error was raised when Run was called twice")
		}
	})

	expectedOutput = "\n"
	if output != expectedOutput {
		t.Errorf("Command output should have been %#v, but was %#v", expectedOutput, output)
	}
}

func captureOutput(f func()) string {
	reader, writer, err := os.Pipe()
	if err != nil {
		panic(err)
	}
	stdout := os.Stdout
	stderr := os.Stderr
	defer func() {
		os.Stdout = stdout
		os.Stderr = stderr
		log.SetOutput(os.Stderr)
	}()
	os.Stdout = writer
	os.Stderr = writer
	log.SetOutput(writer)
	out := make(chan string)
	wg := new(sync.WaitGroup)
	wg.Add(1)
	go func() {
		var buf bytes.Buffer
		wg.Done()
		io.Copy(&buf, reader)
		out <- buf.String()
	}()
	wg.Wait()
	f()
	writer.Close()
	return <-out
}
