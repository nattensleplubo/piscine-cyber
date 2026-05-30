package checker

import (
	"encoding/binary"
	"fmt"
	"os"
)

func CheckIfValidJpeg(file *os.File) (bool, uint16) {
	var marker uint16
	binary.Read(file, binary.BigEndian, &marker)
	if marker != 0xFFD8 {
		fmt.Println("Not a valid JPEG")
		return false, marker
	}
	return true, marker
}
