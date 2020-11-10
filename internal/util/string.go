package util

import (
	"fmt"
	"reflect"
	"strings"
	"time"
	"unsafe"
)

// String2ByteSlice string to []byte nocopy
func String2ByteSlice(str string) (bs []byte) {
	strHdr := (*reflect.StringHeader)(unsafe.Pointer(&str))
	sliceHdr := (*reflect.SliceHeader)(unsafe.Pointer(&bs))
	sliceHdr.Data = strHdr.Data
	sliceHdr.Cap = strHdr.Len
	sliceHdr.Len = strHdr.Len
	return
}

// ByteSlice2String []byte to string nocopy
func ByteSlice2String(bs []byte) string {
	return *(*string)(unsafe.Pointer(&bs))
}

// JSONTime json time
type JSONTime time.Time

// MarshalJSON impl MarshalJSON method
func (t JSONTime) MarshalJSON() ([]byte, error) {
	tt := fmt.Sprintf("\"%s\"", time.Time(t).Format("2006-01-02 15:04:05"))
	fmt.Println(tt)
	return String2ByteSlice(tt), nil
}

func RemoveSlash(url string) string {
	urla := []string{}
	for _, v := range strings.Split(url, "\n") {
		urla = append(urla, strings.TrimRight(v, `\`))
	}
	return strings.Join(urla, "")
}
