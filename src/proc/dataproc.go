package proc

import (
	"encoding/binary"
	"jsutils"
	"unsafe"
)

func EncodeAuth2GStep1(random string,seq uint16) ([]byte,bool) {
	var rrr [] byte
	var temp = [...]byte {
		0x00, 0x04, 0x05, 0x10, 0x00, 0x02, 0x00, 0x02, 0x3f, 0x00, 0x7f, 0x20, 0x00, 0x00, 0x00, 0x00, 0xa0, 0x88, 0x00, 0x00, 0x10,0x74, 0x11, 0x8e, 0x74, 0x35, 0xec, 0x63, 0x98, 0xfb, 0x57, 0xe9, 0xb4, 0x81, 0xc1, 0x0e, 0x3c}
	if len(random) <32{
		return rrr,false
	}
	ran := jsutils.Hex2Byte(random)


	newseq := binary.BigEndian.Uint16(jsutils.Uint16ToBytes(uint16(seq)))

	temp[0] =  byte(newseq%256)
	temp[1] = byte(newseq/256)

	ret := temp[:21]

	ret = append(ret,ran...)
	return ret,true

}

func EncodeAuth2GStep2(seq uint16) ([]byte,bool) {

	var temp = [...]byte {
		0x00, 0x05, 0x05, 0x00, 0x00, 0x0f, 0x00, 0x00 , 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xa0, 0xc0, 0x00, 0x00, 0x0c}

	newseq := binary.BigEndian.Uint16(jsutils.Uint16ToBytes(uint16(seq)))

	temp[0] =  byte(newseq%256)
	temp[1] = byte(newseq/256)

	ret := temp[:]
	return ret,true

}


func EncodeAuth4GStep1(random string,autn string,seq uint16) ([]byte,bool) {
	var rrr [] byte
	var temp = [...]byte {
		0x00, 0x0d, 0x05, 0x22, 0x00, 0x02, 0x00, 0x02,	0x3f, 0x00, 0x7f, 0xff, 0x00, 0x00, 0x00, 0x00, 0x00, 0x88, 0x00, 0x81, 0x22, 0x10, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,	0x10, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00	}
	if len(random) <32{
		return rrr,false
	}
	ran := jsutils.Hex2Byte(random)
	aut := jsutils.Hex2Byte(autn)

	newseq := binary.BigEndian.Uint16(jsutils.Uint16ToBytes(uint16(seq)))

	temp[0] =  byte(newseq%256)
	temp[1] = byte(newseq/256)

	ret := temp[:22]

	ret = append(ret,ran...)
	ret = append(ret,aut...)
	return ret,true

}

func EncodeAuth4GStep2(seq uint16) ([]byte,bool) {

	var temp = [...]byte {
		0x00, 0x0d, 0x05, 0x00, 0x00, 0x38, 0x00, 0x00,	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xc0, 0x00, 0x00, 0x35}

	newseq := binary.BigEndian.Uint16(jsutils.Uint16ToBytes(uint16(seq)))

	temp[0] =  byte(newseq%256)
	temp[1] = byte(newseq/256)

	ret := temp[:]
	return ret,true

}
//result,sres,kc,ret
func DecodeAuth2GStep1(instr []byte) (int,string,string,int) {
	var sres,kc string


	var head DataRespHead
	headLen := unsafe.Sizeof(head)
	srcLen := len(instr)


	if srcLen <  int(headLen)+4{
		if srcLen == 12 && instr[10] == 0x9f && instr[11] == 0x0c{
			return 1, "", "", 99
		}
		return 1, "", "", 0
	}
	temp := instr[headLen:]
	if temp[1] == 0x0c || temp[2] == 0xc0{
		tpSres := temp[3:7]
		sres = jsutils.DisplayHexString(tpSres,0)
		tpKc := temp[7:15]
		kc = jsutils.DisplayHexString(tpKc,0)
		return 0,sres,kc,0
	}
	return 1,sres,kc,2
}


//result,sres,kc,ret
func DecodeAuth2GStep2(instr []byte) (int,string,string,int) {
	var sres,kc string


	var head DataRespHead
	headLen := unsafe.Sizeof(head)
	srcLen := len(instr)
	if srcLen <  int(headLen)+4{
		return 1, "", "", 0
	}
	temp := instr[headLen:]

	if temp[0] == 0xc0{
		tpSres := temp[1:5]
		sres = jsutils.DisplayHexString(tpSres,0)
		tpKc := temp[5:13]
		kc = jsutils.DisplayHexString(tpKc,0)
		return 0,sres,kc,0
	}
	return 1,sres,kc,2
}

//result,sres,ik,kc,ret
func DecodeAuth4GStep1(instr []byte) (int,string,string,string,int) {
	var sres,ik,kc string


	var head DataRespHead
	headLen := unsafe.Sizeof(head)
	srcLen := len(instr)


	if srcLen <  int(headLen)+4{
		if srcLen == 12 && instr[10] == 0x61 && instr[11] == 0x35{
			return 1, "", "", "",99
		}
		return 1, "", "","", 0
	}

	temp := instr[headLen:]
	if temp[2] == 0xc0{
		if temp[3] == 0xdb {		//鉴权成功
			sres_len := temp[4]
			if sres_len == 8{
				tpSres := temp[5:13]
				sres = jsutils.DisplayHexString(tpSres,0)
				kc_len := temp[13]

				if kc_len == 0x10{
					tpKc := temp[14:30]
					kc = jsutils.DisplayHexString(tpKc,0)
					ik_len := temp[30]
					if ik_len == 0x10{
						tpIk := temp[31:47]
						ik = jsutils.DisplayHexString(tpIk,0)
						return 0,sres,ik,kc,0
					}
				}
			}
		}else if temp[3] == 0xdc{		//鉴权失败
			return 1,sres,ik,kc,1
		}

	}
	return 1,sres,ik,kc,2
}

//result,sres,ik,kc,ret
func DecodeAuth4GStep2(instr []byte) (int,string,string,string,int) {
	var sres,ik,kc string


	var head DataRespHead
	headLen := unsafe.Sizeof(head)
	srcLen := len(instr)


	if srcLen <  int(headLen)+4{
		if srcLen == 12 && instr[10] == 0x61 && instr[11] == 0x35{
			return 1, "", "", "",99
		}
		return 1, "", "","", 1
	}

	temp := instr[12:]
	if temp[0] == 0xc0{
		if temp[1] == 0xdb {		//鉴权成功
			sres_len := temp[2]
			if sres_len == 8{
				tpSres := temp[3:11]
				sres = jsutils.DisplayHexString(tpSres,0)
				kc_len := temp[11]

				if kc_len == 0x10{
					tpKc := temp[12:28]
					kc = jsutils.DisplayHexString(tpKc,0)
					ik_len := temp[28]
					if ik_len == 0x10{
						tpIk := temp[29:45]
						ik = jsutils.DisplayHexString(tpIk,0)
						return 0,sres,ik,kc,0
					}
				}
			}
		}

	}
	return 1,sres,ik,kc,2
}

func EncodeImsi2GStep1(index uint16) ([]byte,bool) {

	var temp = [...]byte {
		0x00, 0x05, 0x05, 0x00, 0x00, 0x12, 0x00, 0x03, 0x3f, 0x00, 0x7f, 0x20, 0x6f, 0x07, 0x00, 0x00, 0xa0, 0xc0, 0x00, 0x00, 0x0f}

	newseq := binary.BigEndian.Uint16(jsutils.Uint16ToBytes(uint16(index)))

	temp[0] =  byte(newseq%256)
	temp[1] = byte(newseq/256)

	ret := temp[:]

	return ret,true
}


func EncodeImsi2GStep2(index uint16) ([]byte,bool) {

	var temp = [...]byte {
		0x00, 0x05, 0x05, 0x00, 0x00, 0x0c, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xa0, 0xb0, 0x00, 0x00, 0x09}

	newseq := binary.BigEndian.Uint16(jsutils.Uint16ToBytes(uint16(index)))

	temp[0] =  byte(newseq%256)
	temp[1] = byte(newseq/256)

	ret := temp[:]

	return ret,true
}

func EncodeImsi4GStep1(index uint16) ([]byte,bool) {

	var temp = [...]byte {
		0x00, 0x0d, 0x05, 0x28, 0x00, 0x02, 0x00, 0x03,0x3f, 0x00, 0x7f, 0xff, 0x6f, 0x07, 0x00, 0x00, 0x00, 0xc0, 0x00, 0x00, 0x25}

	newseq := binary.BigEndian.Uint16(jsutils.Uint16ToBytes(uint16(index)))

	temp[0] =  byte(newseq%256)
	temp[1] = byte(newseq/256)

	ret := temp[:]

	return ret,true
}


func EncodeImsi4GStep2(index uint16) ([]byte,bool) {

	var temp = [...]byte {
		0x00, 0x05, 0x05, 0x00, 0x00, 0x0c, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xb0, 0x00, 0x00, 0x09}

	newseq := binary.BigEndian.Uint16(jsutils.Uint16ToBytes(uint16(index)))

	temp[0] =  byte(newseq%256)
	temp[1] = byte(newseq/256)

	ret := temp[:]

	return ret,true
}

