package main

import (
	"fmt"
	"os"
	"path/filepath"
	"scorpion/checker"
	"scorpion/exif"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: ./scorpion <file1> <file2> ...")
		return
	}

	for i := 1; i < len(os.Args); i++ {
		filename := os.Args[i]
		file, err := os.Open(filename)
		if err != nil {
			fmt.Printf("Error opening file %s: %v\n", filename, err)
			continue
		}

		ext := filepath.Ext(filename)
		switch ext {
		case ".jpg", ".jpeg":
			isValid, marker := checker.CheckIfValidJpeg(file)
			if isValid {
				exif.GetJpegEXIF(file, marker, filename)
			} else {
				fmt.Printf("%s is not a valid JPEG\n", filename)
			}
		case ".png":
			if checker.CheckIfValidPng(file) {
				exif.GetPngMetadata(file, filename)
			} else {
				fmt.Printf("%s is not a valid PNG\n", filename)
			}
		case ".gif":
			if checker.CheckIfValidGif(file) {
				exif.GetGifMetadata(file, filename)
			} else {
				fmt.Printf("%s is not a valid GIF\n", filename)
			}
		case ".bmp":
			if checker.CheckIfValidBmp(file) {
				exif.GetBmpMetadata(file, filename)
			} else {
				fmt.Printf("%s is not a valid BMP\n", filename)
			}
		default:
			fmt.Printf("Unsupported file extension: %s\n", ext)
		}
		file.Close()
	}
}
