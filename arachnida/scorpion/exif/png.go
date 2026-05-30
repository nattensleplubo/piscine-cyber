package exif

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

func GetPngMetadata(file *os.File, filename string) {
	fmt.Printf("\n\n====================================================================\nMETADATA FOR %s (PNG)\n", filename)
	
	// Skip PNG signature
	file.Seek(8, 0)

	for {
		var length uint32
		err := binary.Read(file, binary.BigEndian, &length)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("Error reading chunk length:", err)
			return
		}

		chunkType := make([]byte, 4)
		io.ReadFull(file, chunkType)
		typeName := string(chunkType)

		switch typeName {
		case "IHDR":
			var width, height uint32
			binary.Read(file, binary.BigEndian, &width)
			binary.Read(file, binary.BigEndian, &height)
			fmt.Printf("Dimensions: %d x %d\n", width, height)
			file.Seek(int64(length-8), io.SeekCurrent) // Skip rest of IHDR
		case "tEXt":
			data := make([]byte, length)
			io.ReadFull(file, data)
			fmt.Printf("Text Chunk: %s\n", string(data))
		case "eXIf":
			fmt.Printf("Found eXIf chunk (Length: %d)\n", length)
			// Handle EXIF data similarly to JPEG if needed, but for now just acknowledge it
			file.Seek(int64(length), io.SeekCurrent)
		default:
			// Skip unknown chunks
			file.Seek(int64(length), io.SeekCurrent)
		}

		// Skip CRC
		file.Seek(4, io.SeekCurrent)

		if typeName == "IEND" {
			break
		}
	}
}
