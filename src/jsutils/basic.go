package jsutils

import (
	"fmt"
	"strconv"
	"strings"
	"unsafe"
)

func DisplayHexString(source []byte,bank int) string {
	var sa = make([]string, 0)
	for _, v := range source {
		sa = append(sa, fmt.Sprintf("%02X", v))
	}
	bankstr :=""
	for i:=0 ;i<bank;i++ {
		bankstr +=" "
	}
	ss := strings.Join(sa, bankstr)
	return ss
}

func GetIpPort(src string) (string,int) {
	arr:=strings.Split(src,":")
	port,_:=strconv.Atoi(arr[1])
	return arr[0],port
}

func GetStrSlit(src string,key string,pos int) string {
	arr:=strings.Split(src,key)
	return arr[pos-1]
}

func GetJsonStr(insrc string,keyStr string) (string,bool) {
	var bret bool = false
	var retstr string
	if len(insrc)<=0 || len(keyStr)<=0 {
		return retstr,bret
	}
	var keystart string
	keystart ="\""+keyStr
	keystart = keystart+"\""

	index := strings.IndexAny(insrc,keystart)
	if index<0{
		return retstr,bret
	}


	ss := insrc[index+len(keystart):]

	indexend := strings.IndexAny(ss,",")
	if indexend<0{
		indexend = strings.IndexAny(ss,"}")
		if indexend<0 {
			return retstr, bret
		}
	}

	retstr = ss[:indexend]
	bret = true
	return retstr,bret
}

func Uint16ToBytes(n uint16) []byte {
	return []byte{
		byte(n),
		byte(n >> 8),
	}
}

func Int2Byte(data int)(ret []byte){
	var len uintptr = unsafe.Sizeof(data)
	ret = make([]byte, len)
	var tmp int = 0xff
	var index uint = 0
	for index=0; index<uint(len); index++{
		ret[index] = byte((tmp<<(index*8) & data)>>(index*8))
	}
	return ret
}

func Byte2Int(data []byte)int{
	var ret int = 0
	var len int = len(data)
	var i uint = 0
	for i=0; i<uint(len); i++{
		ret = ret | (int(data[i]) << (i*8))
	}
	return ret
}

func Hex2Byte(str string) []byte {
	slen := len(str)
	bHex := make([]byte, len(str)/2)
	ii := 0
	for i := 0; i < len(str); i = i + 2 {
		if slen != 1 {
			ss := string(str[i]) + string(str[i+1])
			bt, _ := strconv.ParseInt(ss, 16, 32)
			bHex[ii] = byte(bt)
			ii = ii + 1
			slen = slen - 2
		}
	}
	return bHex
}