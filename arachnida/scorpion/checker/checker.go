package checker

import (
	"bytes"
	"encoding/binary"
	"os"
)

func CheckIfValidJpeg(file *os.File) (bool, uint16) {
	file.Seek(0, 0)
	var marker uint16
	binary.Read(file, binary.BigEndian, &marker)
	if marker != 0xFFD8 {
		return false, marker
	}
	return true, marker
}

func CheckIfValidPng(file *os.File) bool {
	file.Seek(0, 0)
	header := make([]byte, 8)
	file.Read(header)
	pngSignature := []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}
	return bytes.Equal(header, pngSignature)
}

func CheckIfValidGif(file *os.File) bool {
	file.Seek(0, 0)
	header := make([]byte, 3)
	file.Read(header)
	return string(header) == "GIF"
}

func CheckIfValidBmp(file *os.File) bool {
	file.Seek(0, 0)
	header := make([]byte, 2)
	file.Read(header)
	return string(header) == "BM"
}
