package exif

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

func GetJpegEXIF(file *os.File, marker uint16, filename string) {
	fmt.Printf("\n\n====================================================================\nEXIF DATA FOR %s\n", filename)
	for {
		// Read next marker
		err := binary.Read(file, binary.BigEndian, &marker)
		if err != nil {
			break
		}

		// Read segment length (length includes these 2 bytes)
		var length uint16
		binary.Read(file, binary.BigEndian, &length)

		// 0xFFE1 is the APP1 Marker (EXIF)
		if marker == 0xFFE1 {
			fmt.Printf("Found APP1 (EXIF) Segment. Length: %d\n", length)

			// Verify Exif header: "Exif\0\0"
			buf := make([]byte, 6)
			io.ReadFull(file, buf)
			if string(buf) != "Exif\x00\x00" {
				fmt.Println("Invalid Exif header")
				return
			}

			// Read TIFF header (shoud start now)
			// 1. Endianness (2 bytes): "II" (little) or "MM" (big)
			endian_buf := make([]byte, 2)
			io.ReadFull(file, endian_buf)

			var byte_order binary.ByteOrder
			if string(endian_buf) == "II" {
				byte_order = binary.LittleEndian
				// fmt.Println("Endianness: Little endian -> Intel")
			} else {
				byte_order = binary.BigEndian
				// fmt.Println("Endianness: Big endian -> Motorola")
			}

			// 2. Skip 2 bytes (there's a fixed number 0x002A)
			var fixed uint16
			binary.Read(file, byte_order, &fixed)

			// 3. Offset to first IFD (usually 8 bytes from start to TIFF header)
			var first_IFD_offset uint32
			binary.Read(file, byte_order, &first_IFD_offset)
			// fmt.Printf("First IFD offset : %d\n", first_IFD_offset)

			// Record the start of the TIFF header for relative offsets
			// It started at the II/MM bytes, which is 8 bytes before our current position
			// (2 bytes endian + 2 bytes fixed + 4 bytes offset)
			tiffHeaderOffset, _ := file.Seek(0, io.SeekCurrent)
			tiffHeaderOffset -= 8

			// Seek to the first IFD
			file.Seek(tiffHeaderOffset+int64(first_IFD_offset), io.SeekStart)

			// Read Number of Tags (2 bytes)
			var numTags uint16
			binary.Read(file, byte_order, &numTags)
			fmt.Printf("Number of tags: %d\n", numTags)

			// Loop through tags (each is 12 bytes)
			for i := 0; i < int(numTags); i++ {
				var tagID uint16
				var dataType uint16
				var count uint32
				var valueOffset uint32

				binary.Read(file, byte_order, &tagID)
				binary.Read(file, byte_order, &dataType)
				binary.Read(file, byte_order, &count)
				binary.Read(file, byte_order, &valueOffset)

				value := readTagValue(file, byte_order, tiffHeaderOffset, dataType, count, valueOffset)
				fmt.Printf("Tag: 0x%04X | Value: %s\n", tagID, value)
			}

			return
		}

		file.Seek(int64(length-2), io.SeekCurrent)
	}
}

func readTagValue(file *os.File, order binary.ByteOrder, tiffOffset int64, dataType uint16, count uint32, valueOffset uint32) string {
	// Type 2 is ASCII
	if dataType == 2 {
		currentPos, _ := file.Seek(0, io.SeekCurrent)
		defer file.Seek(currentPos, io.SeekStart)

		var data []byte
		if count <= 4 {
			data = make([]byte, 4)
			order.PutUint32(data, valueOffset)
			data = data[:count]
		} else {
			file.Seek(tiffOffset+int64(valueOffset), io.SeekStart)
			data = make([]byte, count)
			io.ReadFull(file, data)
		}
		// Clean up null bytes and whitespace
		return string(data)
	}

	// For other types, just return the number as a string for now
	return fmt.Sprintf("%d", valueOffset)
}
