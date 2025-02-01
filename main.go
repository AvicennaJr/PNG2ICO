package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"image"
	"image/png"
	"io"
	"os"
	"path/filepath"
	"strings"
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

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: program <input_path>")
		fmt.Println("Input path can be a single PNG file or a directory containing PNG files")
		return
	}

	inputPath := os.Args[1]
	fileInfo, err := os.Stat(inputPath)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if fileInfo.IsDir() {
		// Process directory
		err := filepath.Walk(inputPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() && strings.ToLower(filepath.Ext(path)) == ".png" {
				err := convertToICO(path)
				if err != nil {
					fmt.Printf("Error converting %s: %v\n", path, err)
				} else {
					fmt.Printf("Successfully converted %s\n", path)
				}
			}
			return nil
		})
		if err != nil {
			fmt.Printf("Error walking directory: %v\n", err)
		}
	} else {
		// Process single file
		if strings.ToLower(filepath.Ext(inputPath)) != ".png" {
			fmt.Println("Error: Input file must be a PNG image")
			return
		}
		err := convertToICO(inputPath)
		if err != nil {
			fmt.Printf("Error converting file: %v\n", err)
		} else {
			fmt.Printf("Successfully converted %s\n", inputPath)
		}
	}
}

func convertToICO(inputPath string) error {
	// Read PNG file
	inputFile, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer inputFile.Close()

	// Decode PNG
	img, err := png.Decode(inputFile)
	if err != nil {
		return err
	}

	// Create output file
	outputPath := strings.TrimSuffix(inputPath, filepath.Ext(inputPath)) + ".ico"
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	// Write ICO
	err = writeICO(outputFile, img)
	if err != nil {
		return err
	}

	return nil
}

func writeICO(w io.Writer, img image.Image) error {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	// Create a buffer for the PNG data
	var pngBuf bytes.Buffer
	err := png.Encode(&pngBuf, img)
	if err != nil {
		return err
	}

	// Write ICO header
	header := icoHeader{
		Reserved: 0,
		Type:     1,
		Count:    1,
	}
	err = binary.Write(w, binary.LittleEndian, &header)
	if err != nil {
		return err
	}

	// Write ICO directory
	directory := icoDirectory{
		Width:    byte(width),
		Height:   byte(height),
		Colors:   0,
		Reserved: 0,
		Planes:   1,
		BitCount: 32,
		Size:     uint32(pngBuf.Len()),
		Offset:   uint32(6 + 16), // size of header + size of directory
	}
	err = binary.Write(w, binary.LittleEndian, &directory)
	if err != nil {
		return err
	}

	// Write PNG data
	_, err = w.Write(pngBuf.Bytes())
	if err != nil {
		return err
	}

	return nil
}
