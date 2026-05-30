package exif

import (
	"encoding/binary"
	"fmt"
	"os"
)

func GetBmpMetadata(file *os.File, filename string) {
	fmt.Printf("\n\n====================================================================\nMETADATA FOR %s (BMP)\n", filename)

	// Skip File Header (14 bytes)
	file.Seek(14, 0)

	// DIB Header
	var headerSize uint32
	binary.Read(file, binary.LittleEndian, &headerSize)
	
	if headerSize >= 12 {
		var width, height int32
		if headerSize == 12 {
			// BITMAPCOREHEADER
			var w, h uint16
			binary.Read(file, binary.LittleEndian, &w)
			binary.Read(file, binary.LittleEndian, &h)
			width = int32(w)
			height = int32(h)
		} else {
			// BITMAPINFOHEADER and newer
			binary.Read(file, binary.LittleEndian, &width)
			binary.Read(file, binary.LittleEndian, &height)
		}
		
		var planes uint16
		var bitCount uint16
		binary.Read(file, binary.LittleEndian, &planes)
		binary.Read(file, binary.LittleEndian, &bitCount)
		
		fmt.Printf("Dimensions: %d x %d\n", width, height)
		fmt.Printf("Planes: %d\n", planes)
		fmt.Printf("Bits per Pixel: %d\n", bitCount)
	} else {
		fmt.Println("Unknown DIB header size")
	}
}
