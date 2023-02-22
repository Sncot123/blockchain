package BLC

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
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

// 标准的JSON格式转换为切片
// 格式 cli.exe send -from "[\"tronytan\"]" -to "[\"Alice\"" -amount "[\"100\"]"
func JSONToSlice(jsonString string) []string {
	var strSlice []string
	//json
	if err := json.Unmarshal([]byte(jsonString), &strSlice); err != nil {
		log.Panicf("json to []string failed err:%v\n", err)
	}
	return strSlice
}
