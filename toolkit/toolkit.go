package toolkit

import (
	"flag"
	"fmt"
	"os"
	"os/signal"

	"github.com/fatih/color"
)
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func HandleInterrupt() {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	go func() {
		<-signalChan
		color.Yellow("\nOperation cancelled by user")
		os.Exit(1)
	}()
}

func PrintUsage() {
	fmt.Println("Usage: png2ico [OPTIONS] <input-path>")
	fmt.Println("Options:")
	flag.PrintDefaults()
}