package BLC

import (
	"bytes"
	"encoding/binary"
	"log"
)

func IntToHex(data int64) []byte {
	buffer := new(bytes.Buffer)
	err := binary.Write(buffer, binary.BigEndian, data)
	if err != nil {
		log.Panicf("parse int64 to []byte failed,err:%v\n ", err)
	}
	return buffer.Bytes()
}
