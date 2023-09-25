package egg

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/gen2brain/beeep"
	"github.com/hajimehoshi/go-mp3"
	"github.com/hajimehoshi/oto"
)

const successSound = "sounds/success.mp3"
const errorSound = "sounds/error.mp3"

// Command contains the current command run by egg
type Command struct {
	Arguments []string

	ran      bool
	cmd      *exec.Cmd
	duration time.Duration
	err      error
}

// NewCommand returns a Command if enough args are given
func NewCommand(args []string) (*Command, error) {
	if len(args) < 1 {
		return nil, errors.New("No command was given")

	}

	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return &Command{
		Arguments: args,
		cmd:       cmd,
	}, nil
}

// Run will run the command keepting track of its duration
func (c *Command) Run() error {
	if c.ran {
		return errors.New("Attempt to run command twice")
	}

	c.ran = true
	started := time.Now()
	if c.err = c.cmd.Start(); c.err != nil {
		return c.err
	}

	c.err = c.cmd.Wait()
	completed := time.Now()

	c.duration = completed.Sub(started)

	return c.err
}

// AnnounceIntent prints out what egg is about to run + workingDir
func (c *Command) AnnounceIntent() {
	fmt.Printf("📣 [%s] %s\n", c.workingDir(), strings.Join(c.Arguments, " "))
}

// NotifyStatus will send an os notification about the status of the command
func (c *Command) NotifyStatus() {
	if c.ran {
		fullCommand := strings.Join(c.Arguments, " ")
		resultIcon := "✅"
		resultText := "completed"
		if c.err != nil {
			resultIcon = "❌"
			resultText = fmt.Sprintf("exited with code %d", c.ExitCode())
		}

		beeep.Notify(
			fmt.Sprintf("%s %s", resultIcon, fullCommand),
			fmt.Sprintf("[%s] %s (%s)", c.workingDir(), resultText, c.duration),
			"na",
		)
	}
}

// PlaySound will play either the error or success sound
func (c *Command) PlaySound() error {

	if c.ran {
		sound := successSound
		if c.err != nil {
			sound = errorSound
		}
		a, err := Asset(sound)
		if err != nil {
			return err
		}
		r := bytes.NewReader(a)
		d, err := mp3.NewDecoder(r)
		if err != nil {
			return err
		}

		ctx, err := oto.NewContext(d.SampleRate(), 2, 2, 8192)
		if err != nil {
			return err
		}
		defer ctx.Close()

		p := ctx.NewPlayer()
		defer p.Close()

		if _, err := io.Copy(p, d); err != nil {
			return err
		}
	}

	return nil
}

// ExitCode will return the code the command returned
func (c *Command) ExitCode() int {
	if c.err != nil {
		if exiterr, ok := c.err.(*exec.ExitError); ok {
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				return status.ExitStatus()
			}

			return 1
		}
		return 1
	}

	return 0
}

func (c Command) workingDir() string {
	dir, err := os.Getwd()
	if err != nil {
		log.Println(err)
		return "Unknown"
	}

	return filepath.Base(dir)
}
