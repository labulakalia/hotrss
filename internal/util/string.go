package util

import (
	"fmt"
	"time"
	"unsafe"
)

// StringToByte string to byte
func StringToByte(data string) []byte {
	return *(*[]byte)(unsafe.Pointer(&data))
}

// JSONTime json time
type JSONTime time.Time

// MarshalJSON impl MarshalJSON method
func (t JSONTime) MarshalJSON() ([]byte, error) {
	tt := fmt.Sprintf("\"%s\"", time.Time(t).Format("2006-01-02 15:04:05"))
	fmt.Println(tt)
	return StringToByte(tt), nil
}
