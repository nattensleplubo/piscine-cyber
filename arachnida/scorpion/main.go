package main

import (
	"fmt"
	"os"
	"path/filepath"
	"scorpion/checker"
	"scorpion/exif"
)

/*

Scorpion:
	The program receive image files and parse them for EXIF and other metadatas,
displaying them on the screen.

	Should be compatible with the files extensions that spider handles :
		[ ] jpg
		[ ] jpeg
		[ ] png
		[ ] gif
		[ ] bmp

*/

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: ./scorpion <file>")
		return
	}

	for i := 1; i < len(os.Args); i++ {
		file, err := os.Open(os.Args[i])
		if err != nil {
			fmt.Println("error opening file ", err)
		}
		defer file.Close()

		file_extension := filepath.Ext(os.Args[i])
		switch file_extension {
		case ".jpg":
			is_valid_jpeg, marker := checker.CheckIfValidJpeg(file)
			if is_valid_jpeg {
				exif.GetJpegEXIF(file, marker, os.Args[i])
			}
		case ".png":
			fmt.Println("CHECKING PNG")
		case ".gif":
			fmt.Println("CHECKING GIF")
		case ".bmp":
			fmt.Println("CHECKING BMP")
		}
	}
}
