package main

import (
	"fmt"
	"log"
	"os"

	"github.com/jpoz/egg"
)

func main() {
	arguments := os.Args[1:]
	cmd, err := egg.NewCommand(arguments)
	if err != nil {
		log.Fatalf("‼️  %s", err)
	}

	cmd.AnnounceIntent()
	err = cmd.Run()
	if err != nil {
		fmt.Println(err)
	}

	cmd.NotifyStatus()

	// TODO make this configurable
	err = cmd.PlaySound()
	if err != nil {
		log.Fatalf("‼️  %s", err)
	}

	os.Exit(cmd.ExitCode())
}
