package exif

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

func GetGifMetadata(file *os.File, filename string) {
	fmt.Printf("\n\n====================================================================\nMETADATA FOR %s (GIF)\n", filename)

	// Skip signature (3 bytes "GIF") and version (3 bytes "87a" or "89a")
	file.Seek(6, 0)

	// Logical Screen Descriptor
	var width, height uint16
	binary.Read(file, binary.LittleEndian, &width)
	binary.Read(file, binary.LittleEndian, &height)
	fmt.Printf("Dimensions: %d x %d\n", width, height)

	// Skip rest of Logical Screen Descriptor (3 bytes) and Global Color Table if it exists
	var packedFields byte
	binary.Read(file, binary.LittleEndian, &packedFields)
	file.Seek(2, io.SeekCurrent) // Skip Background Color Index and Pixel Aspect Ratio

	if packedFields&0x80 != 0 {
		// Global Color Table exists
		size := 1 << ((packedFields & 0x07) + 1)
		file.Seek(int64(3*size), io.SeekCurrent)
	}

	for {
		var blockType byte
		err := binary.Read(file, binary.LittleEndian, &blockType)
		if err != nil {
			break
		}

		if blockType == 0x3B { // Trailer
			break
		}

		if blockType == 0x21 { // Extension
			var extType byte
			binary.Read(file, binary.LittleEndian, &extType)
			
			if extType == 0xFE { // Comment Extension
				fmt.Print("Comment: ")
				for {
					var size byte
					binary.Read(file, binary.LittleEndian, &size)
					if size == 0 {
						break
					}
					data := make([]byte, size)
					io.ReadFull(file, data)
					fmt.Print(string(data))
				}
				fmt.Println()
			} else {
				// Skip other extensions
				for {
					var size byte
					binary.Read(file, binary.LittleEndian, &size)
					if size == 0 {
						break
					}
					file.Seek(int64(size), io.SeekCurrent)
				}
			}
		} else if blockType == 0x2C { // Image Descriptor
			file.Seek(8, io.SeekCurrent)
			var pf byte
			binary.Read(file, binary.LittleEndian, &pf)
			if pf&0x80 != 0 {
				size := 1 << ((pf & 0x07) + 1)
				file.Seek(int64(3*size), io.SeekCurrent)
			}
			// Skip image data blocks
			for {
				var size byte
				binary.Read(file, binary.LittleEndian, &size)
				if size == 0 {
					break
				}
				file.Seek(int64(size), io.SeekCurrent)
			}
		} else {
			// Unknown block
			break
		}
	}
}
