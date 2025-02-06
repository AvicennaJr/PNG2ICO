package main

import (
	"bytes"
	"encoding/binary"
	"flag"
	"fmt"
	"image"
	"image/png"
	"io"
	"os"
	"os/signal"
	"path/filepath"
	"strings"

	"github.com/cheggaaa/pb/v3"
	"github.com/fatih/color"
)

var (
	outputDir string
	force     bool
	verbose   bool
)

type icoHeader struct {
	Reserved uint16
	Type     uint16
	Count    uint16
}

type icoDirectory struct {
	Width    byte
	Height   byte
	Colors   byte
	Reserved byte
	Planes   uint16
	BitCount uint16
	Size     uint32
	Offset   uint32
}

func init() {
	flag.StringVar(&outputDir, "o", "", "Output directory for .ico files")
	flag.StringVar(&outputDir, "output", "", "Output directory for .ico files")
	flag.BoolVar(&force, "f", false, "Overwrite existing files")
	flag.BoolVar(&force, "force", false, "Overwrite existing files")
	flag.BoolVar(&verbose, "v", false, "Verbose output")
	flag.BoolVar(&verbose, "verbose", false, "Verbose output")
}

func main() {
	handleInterrupt()
	flag.Parse()
	args := flag.Args()

	if len(args) < 1 {
		printUsage()
		return
	}

	inputPath := args[0]
	if outputDir != "" {
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			color.Red("Error creating output directory: %v", err)
			return
		}
	}

	fileInfo, err := os.Stat(inputPath)
	if err != nil {
		color.Red("Error: %v", err)
		return
	}

	if fileInfo.IsDir() {
		processDirectory(inputPath)
	} else {
		processSingleFile(inputPath)
	}
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

func processDirectory(inputPath string) {
	var pngFiles []string

	err := filepath.Walk(inputPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.EqualFold(filepath.Ext(path), ".png") {
			pngFiles = append(pngFiles, path)
		}
		return nil
	})

	if err != nil {
		color.Red("Error walking directory: %v", err)
		return
	}

	if len(pngFiles) == 0 {
		color.Yellow("No PNG files found in directory: %s", inputPath)
		return
	}

	bar := pb.StartNew(len(pngFiles))
	successCount := 0

	for _, file := range pngFiles {
		if err := convertToICO(file); err != nil {
			if verbose {
				color.Red("Error converting %s: %v", file, err)
			}
		} else {
			successCount++
			if verbose {
				color.Green("Converted: %s", file)
			}
		}
		bar.Increment()
	}

	bar.Finish()
	color.Green("Successfully converted %d/%d files", successCount, len(pngFiles))
}

func processSingleFile(inputPath string) {
	if !strings.EqualFold(filepath.Ext(inputPath), ".png") {
		color.Red("Input file must be a PNG image")
		return
	}

	if err := convertToICO(inputPath); err != nil {
		color.Red("Conversion failed: %v", err)
	} else {
		color.Green("Successfully converted: %s", inputPath)
	}
}

func convertToICO(inputPath string) error {
	inputFile, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("failed to open input file: %w", err)
	}
	defer inputFile.Close()

	img, err := png.Decode(inputFile)
	if err != nil {
		return fmt.Errorf("invalid PNG format: %w", err)
	}

	if err := validateImageSize(img); err != nil {
		return err
	}

	outputPath, err := constructOutputPath(inputPath)
	if err != nil {
		return err
	}

	if !force && fileExists(outputPath) {
		return fmt.Errorf("output file exists: %s (use -f to overwrite)", outputPath)
	}

	outputFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outputFile.Close()

	if err := writeICO(outputFile, img); err != nil {
		return fmt.Errorf("failed to write ICO file: %w", err)
	}

	return nil
}

func validateImageSize(img image.Image) error {
	bounds := img.Bounds()
	if width, height := bounds.Dx(), bounds.Dy(); width > 256 || height > 256 {
		return fmt.Errorf("image dimensions (%dx%d) exceed maximum 256x256", width, height)
	}
	return nil
}

func constructOutputPath(inputPath string) (string, error) {
	baseName := strings.TrimSuffix(filepath.Base(inputPath), filepath.Ext(inputPath))
	outputFile := baseName + ".ico"

	if outputDir != "" {
		return filepath.Join(outputDir, outputFile), nil
	}
	return filepath.Join(filepath.Dir(inputPath), outputFile), nil
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func writeICO(w io.Writer, img image.Image) error {
	var pngBuf bytes.Buffer
	if err := png.Encode(&pngBuf, img); err != nil {
		return err
	}

	header := icoHeader{
		Reserved: 0,
		Type:     1,
		Count:    1,
	}

	if err := binary.Write(w, binary.LittleEndian, &header); err != nil {
		return err
	}

	bounds := img.Bounds()
	directory := icoDirectory{
		Width:    byte(bounds.Dx()),
		Height:   byte(bounds.Dy()),
		Colors:   0,
		Reserved: 0,
		Planes:   1,
		BitCount: 32,
		Size:     uint32(pngBuf.Len()),
		Offset:   uint32(binary.Size(header) + binary.Size(icoDirectory{})),
	}

	if err := binary.Write(w, binary.LittleEndian, &directory); err != nil {
		return err
	}

	_, err := w.Write(pngBuf.Bytes())
	return err
}
