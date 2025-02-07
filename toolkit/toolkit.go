package toolkit

import (
	"flag"
	"fmt"
	"os"

	"github.com/fatih/color"
)
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func handleInterrupt() {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	go func() {
		<-signalChan
		color.Yellow("\nOperation cancelled by user")
		os.Exit(1)
	}()
}

func printUsage() {
	fmt.Println("Usage: png2ico [OPTIONS] <input-path>")
	fmt.Println("Options:")
	flag.PrintDefaults()
}