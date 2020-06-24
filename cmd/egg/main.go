package main

import (
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
	cmd.Run()
	cmd.NotifyStatus()

	os.Exit(cmd.ExitCode())
}
