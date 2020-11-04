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
	if len(arr) >=pos{
		return arr[pos-1]
	}
	return ""
}

func GetJsonStr(insrc string,keyStr string) (string,bool) {
	var bret bool = false
	var indexend int
	var retstr string
	if len(insrc)<=0 || len(keyStr)<=0 {
		return retstr,bret
	}
	var keystart string
	keystart ="\""+keyStr
	keystart = keystart+"\""

	index := strings.Index(insrc,keystart)
	if index<0{
		return retstr,bret
	}


	ss := insrc[index+len(keystart):]

	indexend1 := strings.Index(ss,",")
	if indexend1>0 {
		indexend = indexend1
	}
	indexend2 := strings.Index(ss,"}")
	if indexend2>0 {
		if indexend>indexend2{
			indexend = indexend2
		}
	}else if indexend<0 {
		return retstr, bret
	}

	retstr = ss[:indexend]

	retstr = strings.TrimLeft(retstr,":")
	retstr = strings.TrimLeft(retstr,"\"")

	retstr = strings.TrimRight(retstr,"}")
	retstr = strings.TrimRight(retstr,",")
	retstr = strings.TrimRight(retstr,"\"")

	bret = true
	return retstr,bret
}

func Uint16ToBytes(n uint16) []byte {
	return []byte{
		byte(n),
		byte(n >> 8),
	}
}

func Uint32ToBytes(n uint32) []byte {
	return []byte{
		byte(n),
		byte(n >> 8),
		byte(n >> 16),
		byte(n >> 24),
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


var g_str string= "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/="


func Base64Encode(data []byte, data_len int) ([]byte,bool){
	var prepare int
	var ret_len int
	var temp int
	var tmp int
	var changed [4]byte
	var pos int
	var base []byte = []byte(g_str)
	var ret []byte

	ret_len = data_len / 3
	temp = data_len % 3
	if (temp > 0){
		ret_len += 1
	}
	ret_len = ret_len * 4
	ret = make([]byte,ret_len)

	for ;tmp < data_len;{
		temp = 0
		prepare = 0
		for ;temp < 3; {
			if tmp >= data_len {
				break
			}
			tt := data[tmp] & 0xFF
			prepare = ((prepare << 8) | int(tt))
			tmp++
			temp++
		}
		prepare = (prepare << ((3 - temp) * 8))
		for i := 0; i < 4; i++{
			if (temp < i){
				changed[i] = 0x40
			}else{
				dsd := (prepare >> ((3 - i) * 6)) & 0x3F
				changed[i] = uint8(dsd)
			}
			ret[pos] = base[changed[i]]
			pos++
		}
	}
	//ret[pos] = 0x00
	return ret,true
}

func Base64FindPos(inchar string)  int{
	return strings.Index(g_str,inchar)
}

func Base64Decode(data []byte, data_len int) ([]byte,bool) {
	var ret_len int= (data_len / 4) * 3
	var equal_count int
	var pos int
	var tmp int
	var temp int

	var prepare int
	var i int
	if (data [data_len - 1] == '='){
		equal_count += 1
	}
	if (data [+data_len - 2] == '='){
		equal_count += 1
	}
	if (data [ data_len - 3] == '='){//seems impossible
		equal_count += 1
	}
	switch (equal_count){
	case 0:
		ret_len += 4	//3 + 1 [1 for NULL]
		break
	case 1:
		ret_len += 4	//Ceil((6*3)/8)+1
		break
	case 2:
		ret_len += 3	//Ceil((6*2)/8)+1
		break
	case 3:
		ret_len += 2	//Ceil((6*1)/8)+1
		break
	}
	var ret = make([]byte,ret_len)

	for ;tmp < (data_len - equal_count);{
		temp = 0
		prepare = 0

		for ;temp < 4;{
			if (tmp >= (data_len - equal_count)){
				break
			}
			var str string = string(data[tmp])
			prepare = (prepare << 6) | (Base64FindPos(str))
			temp++
			tmp++
		}
		prepare = prepare << ((4 - temp) * 6);
		for i = 0; i<3; i++{
			if (i == temp){
				break
			}
			ttt := (prepare >> ((2 - i) * 8))
			ttt = ttt& 0xFF
			ret[pos] = uint8(ttt)
			pos++
		}
	}
	return ret,true
}

func GetNext16(in *int16) int16 {
	*in++
	if *in == 0 {
		*in = 1
	}
	return *in
}

func GetNext32(in *int32) int32 {
	*in++
	if *in == 0 {
		*in = 1
	}
	return *in
}

