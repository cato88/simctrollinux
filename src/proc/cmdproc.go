package proc

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"jsutils"
	"reflect"
	"strconv"
	"unsafe"
)

func IsNeedAck(header *CmdHeadInfo) (uint8, uint16, bool) {
	var bret = false
	var flag uint8 = 0
	var status uint16 = 0

	if header.Magic != 0xa3c5 {
		return flag, status, bret
	}

	if (header.Flag&0x01) == 1 || (header.Flag&0x02) > 0 || (header.Flag&0x80) > 0 || (header.Flag&0x04) > 0 {
		status = 200
		if (header.Flag & 0x01) > 0 {
			flag = 0
			bret = true
		} else if (header.Flag & 0x02) > 0 {
			if (header.Flag & 0x62) > 0 {
				flag = 0x22
				bret = true
			} else if (header.Flag & 0x02) > 0 {
				flag = 0x22
				bret = true
			} else {
				flag = 0
				bret = false
			}
		} else if (header.Flag & 0x04) > 0 {
			flag = 1
			status = header.Status
			bret = true
		}
	}

	return flag, status, bret
}

func EncodeAck(header *CmdHeadInfo, flag uint8, status uint16) ([]byte, int, bool) {
	var bret bool = true
	var data []byte
	var datalen int

	temp := binary.BigEndian.Uint16(jsutils.Uint16ToBytes(status))

	resphear := &CmdHeadInfo{
		Ver:    header.Ver,
		Length: 0,
		Ttid:   binary.BigEndian.Uint32(jsutils.Uint32ToBytes(header.Ttid)),
		Ssrc:   binary.BigEndian.Uint32(jsutils.Uint32ToBytes(header.Ssrc)),
		Status: temp,
		Flag:   header.Flag,
		Magic:  binary.BigEndian.Uint16(jsutils.Uint16ToBytes(header.Magic)),
	}
	datalen = int(unsafe.Sizeof(*resphear))
	var x reflect.SliceHeader
	x.Len = datalen
	x.Cap = datalen
	x.Data = uintptr(unsafe.Pointer(&resphear))
	data = *(*[]byte)(unsafe.Pointer(&x))

	data = data[:datalen]
	return data, datalen, bret
}

func DecodeReg(instr string) (int32, int, string, string, bool) {
	var bret bool = true
	var seq int32 = 0
	var did int = 0
	var simstate string
	var username string

	temp, ok := jsutils.GetJsonStr(instr, ESCC_MSG_PAR_CSEQ)
	if ok == true {
		sint, _ := strconv.Atoi(temp)
		seq = int32(sint)
	}

	temp, ok = jsutils.GetJsonStr(instr, ESCC_MSG_PAR_DID)
	if ok == true {
		sint, _ := strconv.Atoi(temp)
		did = sint
	}

	simstate, ok = jsutils.GetJsonStr(instr, ESCC_MSG_PAR_SIM_STATUS)
	username, ok = jsutils.GetJsonStr(instr, ESCC_MSG_PAR_USERNAME)

	return seq, did, simstate, username, bret
}

func EncodeRegAck(header *CmdHeadInfo, seq int32, did int) ([]byte, int, bool) {
	var bret bool = false
	ack := map[string]interface{}{
		ESCC_MSG_PAR_PROTO:   SIM_PROTO,
		ESCC_MSG_PAR_MSG:     ESCC_MSG_REG_ACK,
		ESCC_MSG_PAR_CSEQ:    seq,
		ESCC_MSG_PAR_DID:     did,
		ESCC_MSG_PAR_CODE:    ESCC_MSG_CODE_OK,
		ESCC_MSG_PAR_REASON:  "OK",
		ESCC_MSG_PAR_EXPIRES: "180",
	}

	data, err := json.Marshal(ack)
	if err != nil {
		bret = false
	}
	datalen := len(data)
	data = data[:datalen]

	//temp := binary.BigEndian.Uint16(jsutils.Uint16ToBytes(200))

	resphear := &CmdHeadInfo{
		Ver:    header.Ver,
		Length: uint16(datalen),
		Ttid:   binary.BigEndian.Uint32(jsutils.Uint32ToBytes(header.Ttid)),
		Ssrc:   binary.BigEndian.Uint32(jsutils.Uint32ToBytes(header.Ssrc)),
		Status: 200,
		Flag:   0,
		Magic:  binary.BigEndian.Uint16(jsutils.Uint16ToBytes(header.Magic)),
	}

	buf := &bytes.Buffer{}
	eee := binary.Write(buf, binary.BigEndian, resphear)
	if eee != nil {
		panic(eee)
	}

	head := buf.Bytes()[:]

	bret = true

	head = append(head, data...)

	return head, len(head), bret
}

func EncodeOpen(header *CmdHeadInfo, serverip string, serverport uint16, seq uint32, tid uint32, slot int, username string) ([]byte, int, bool) {
	var bret bool = false

	var simslot, goipname string
	simslot = fmt.Sprintf("%s.%03d", username, slot)
	goipname = fmt.Sprintf("%s.%03d", username, slot)
	ack := map[string]interface{}{
		ESCC_MSG_PAR_PROTO:     SIM_PROTO,
		ESCC_MSG_PAR_MSG:       ESCC_MSG_OPEN,
		ESCC_MSG_PAR_CSEQ:      seq,
		ESCC_MSG_PAR_FROM:      LOCAL_FROM,
		ESCC_MSG_PAR_TO:        "0",
		ESCC_MSG_PAR_GOIP_SLOT: goipname,
		ESCC_MSG_PAR_SIM_SLOT:  simslot,

		ESCC_MSG_PAR_EXPIRES:          120,
		ESCC_MSG_PAR_CONN_PROTO:       "UDP",
		ESCC_MSG_PAR_CONN_IP:          serverip,
		ESCC_MSG_PAR_CONN_PORT:        serverport,
		ESCC_VAL_CONN_MPORT:           slot,
		ESCC_MSG_PAR_CONN_ACKTIME:     0,
		ESCC_MSG_PAR_CONN_RETRANSTIME: 0,

		ESCC_MSG_PAR_CONN_RETRANSINTVL: 0,
		ESCC_MSG_PAR_CONN_RETRANSCOUNT: 0,
		ESCC_MSG_PAR_MOD_TYPE:          "GSM",
		ESCC_VAL_MOD_VALUE:             "",
		ESCC_MSG_PAR_MOD_IMEI:          "",
		ESCC_MSG_PAR_GET:               "SIM-DATA",
	}

	data, err := json.Marshal(ack)
	if err != nil {
		bret = false
	}
	datalen := len(data)
	data = data[:datalen]

	resphear := &CmdHeadInfo{
		Ver:    header.Ver,
		Length: uint16(datalen),
		Ttid:   tid,
		Ssrc:   binary.BigEndian.Uint32(jsutils.Uint32ToBytes(header.Ssrc)),
		Status: 0,
		Flag:   0x40,
		Magic:  binary.BigEndian.Uint16(jsutils.Uint16ToBytes(header.Magic)),
	}

	buf := &bytes.Buffer{}
	eee := binary.Write(buf, binary.BigEndian, resphear)
	if eee != nil {
		panic(eee)
	}
	bret = true
	head := buf.Bytes()[:]

	head = append(head, data...)

	return head, len(head), bret
}

func DecodeOpenAck(instr string) (string, string, string, string, string, int, int, int, string, int, int) {
	var imsi, imei, iccid, from, to, remoteip, simslot, temp, contentlength string
	var code, slot, clen, unclen, remoteport, simtype int
	var ok bool

	from, ok = jsutils.GetJsonStr(instr, ESCC_MSG_PAR_FROM)
	to, ok = jsutils.GetJsonStr(instr, ESCC_MSG_PAR_TO)
	iccid, ok = jsutils.GetJsonStr(instr, ESCC_MSG_PAR_SIM_ICCID)
	imsi, ok = jsutils.GetJsonStr(instr, ESCC_MSG_PAR_SIM_IMSI)
	imei, ok = jsutils.GetJsonStr(instr, ESCC_MSG_PAR_SIM_IMEI)
	remoteip, ok = jsutils.GetJsonStr(instr, ESCC_MSG_PAR_CONN_IP)

	simslot, ok = jsutils.GetJsonStr(instr, ESCC_MSG_PAR_SIM_SLOT)
	if ok == true {
		retStr := jsutils.GetStrSlit(simslot, ".", 2)
		slot, _ = strconv.Atoi(retStr)
	}

	contentlength, ok = jsutils.GetJsonStr(instr, ESCC_MSG_PAR_CTT_LEN)
	if ok == true {
		jsutils.Fatal("----------->contentlength", contentlength, "<----------------")
		temp = jsutils.GetStrSlit(contentlength, "/", 1)
		clen, _ = strconv.Atoi(temp)

		temp = jsutils.GetStrSlit(contentlength, "/", 2)
		unclen, _ = strconv.Atoi(temp)

		jsutils.Fatal("clen", clen, "unclen", unclen)
	} else {
		jsutils.Fatal("----------->get contentlength  error<----------------")
	}

	temp, ok = jsutils.GetJsonStr(instr, ESCC_MSG_PAR_CODE)
	if ok == true {
		code, _ = strconv.Atoi(temp)
	}

	temp, ok = jsutils.GetJsonStr(instr, ESCC_MSG_PAR_CONN_PORT)
	if ok == true {
		remoteport, _ = strconv.Atoi(temp)
	}
	if clen > 0 && unclen > 0 {
		unpackdata := instr[len(instr)-clen:]
		var srcdata []byte = []byte(unpackdata)
		dst := jsutils.MyZlibUnCompress(srcdata)
		if len(dst) >= 16 {
			tt := dst[12:16]
			simtype = int(binary.BigEndian.Uint32(tt))
		} else {
			jsutils.Fatal("DecodeOpenAck jsutils.MyZlibUnCompress error,len", len(dst))
		}
	}

	return imsi, imei, iccid, from, to, code, slot, unclen, remoteip, remoteport, simtype
}

func EncodeOpenAckAck(header *CmdHeadInfo) ([]byte, int, bool) {
	var bret bool = false

	var status uint16 = 200
	//status = binary.BigEndian.Uint16(jsutils.Uint16ToBytes(uint16(status)))

	resphear := &CmdHeadInfo{
		Ver:    header.Ver,
		Length: 0,
		Ttid:   binary.BigEndian.Uint32(jsutils.Uint32ToBytes(header.Ttid)),
		Ssrc:   binary.BigEndian.Uint32(jsutils.Uint32ToBytes(header.Ssrc)),
		Status: status,
		Flag:   0x01,
		Magic:  binary.BigEndian.Uint16(jsutils.Uint16ToBytes(header.Magic)),
	}

	buf := &bytes.Buffer{}
	eee := binary.Write(buf, binary.BigEndian, resphear)
	if eee != nil {
		panic(eee)
	}
	bret = true
	head := buf.Bytes()[:]

	return head, len(head), bret
}

func DecodeClose(instr string) (string, string, int, int) {
	var from, to, simslot, temp string
	var slot, cseq int
	var ok bool

	from, ok = jsutils.GetJsonStr(instr, ESCC_MSG_PAR_FROM)
	to, ok = jsutils.GetJsonStr(instr, ESCC_MSG_PAR_TO)
	temp, ok = jsutils.GetJsonStr(instr, ESCC_MSG_PAR_CSEQ)
	if ok == true {
		cseq, _ = strconv.Atoi(temp)
	}

	simslot, ok = jsutils.GetJsonStr(instr, ESCC_MSG_PAR_SIM_SLOT)
	if ok == true {
		retStr := jsutils.GetStrSlit(simslot, ".", 2)
		slot, _ = strconv.Atoi(retStr)
	}

	return from, to, cseq, slot
}

func EncodeCloseAck(header *CmdHeadInfo, cSeq int, from string, to string) ([]byte, int, bool) {
	var bret bool = false

	var status uint16 = 200
	//status = binary.BigEndian.Uint16(jsutils.Uint16ToBytes(uint16(status)))

	ack := map[string]interface{}{
		ESCC_MSG_PAR_PROTO: SIM_PROTO,
		ESCC_MSG_PAR_MSG:   ESCC_MSG_CLOSE_ACK,
		ESCC_MSG_PAR_CSEQ:  cSeq,
		ESCC_MSG_PAR_FROM:  from,
		ESCC_MSG_PAR_TO:    to,

		ESCC_MSG_PAR_CODE:   ESCC_MSG_CODE_OK,
		ESCC_MSG_PAR_REASON: "OK",
	}

	data, err := json.Marshal(ack)
	if err != nil {
		bret = false
	}
	datalen := len(data)
	data = data[:datalen]

	resphear := &CmdHeadInfo{
		Ver:    header.Ver,
		Length: uint16(datalen),
		Ttid:   binary.BigEndian.Uint32(jsutils.Uint32ToBytes(header.Ttid)),
		Ssrc:   binary.BigEndian.Uint32(jsutils.Uint32ToBytes(header.Ssrc)),
		Status: status,
		Flag:   0x00,
		Magic:  binary.BigEndian.Uint16(jsutils.Uint16ToBytes(header.Magic)),
	}

	buf := &bytes.Buffer{}
	eee := binary.Write(buf, binary.BigEndian, resphear)
	if eee != nil {
		panic(eee)
	}

	head := buf.Bytes()[:]

	bret = true
	head = append(head, data...)
	return head, len(head), bret
}

func DecodePublish(instr string) (string, string, string, int, int, int) {
	var username, from, to, simslot, temp string
	var simstate, slot, cseq int
	var ok bool

	from, ok = jsutils.GetJsonStr(instr, ESCC_MSG_PAR_FROM)
	to, ok = jsutils.GetJsonStr(instr, ESCC_MSG_PAR_TO)
	temp, ok = jsutils.GetJsonStr(instr, ESCC_MSG_PAR_CSEQ)
	if ok == true {
		cseq, _ = strconv.Atoi(temp)
	}

	simslot, ok = jsutils.GetJsonStr(instr, ESCC_MSG_PAR_SIM_SLOT)
	if ok == true {
		retStr := jsutils.GetStrSlit(simslot, ".", 2)
		slot, _ = strconv.Atoi(retStr)
	}

	username, ok = jsutils.GetJsonStr(instr, ESCC_MSG_PAR_USERNAME)

	temp, ok = jsutils.GetJsonStr(instr, ESCC_VAL_HAS_SIM)
	if ok == true {
		simstate, _ = strconv.Atoi(temp)
	}

	return username, from, to, slot, cseq, simstate
}

func EncodePublishAck(header *CmdHeadInfo, cSeq int, from string, to string, username string, simstate int, slot int) ([]byte, int, bool) {
	var bret bool = false

	var status uint16 = 200
	var simslot string
	simslot = fmt.Sprintf("%s.%03d", username, slot)
	//status = binary.BigEndian.Uint16(jsutils.Uint16ToBytes(uint16(status)))

	ack := map[string]interface{}{
		ESCC_MSG_PAR_PROTO:      SIM_PROTO,
		ESCC_MSG_PAR_MSG:        ESCC_MSG_PUBLISH_ACK,
		ESCC_MSG_PAR_CSEQ:       cSeq,
		ESCC_MSG_PAR_FROM:       from,
		ESCC_MSG_PAR_TO:         to,
		ESCC_MSG_PAR_SIM_SLOT:   simslot,
		ESCC_MSG_PAR_USERNAME:   username,
		ESCC_MSG_PAR_CODE:       ESCC_MSG_CODE_OK,
		ESCC_MSG_PAR_REASON:     "OK",
		ESCC_MSG_PAR_SIM_STATUS: simstate,
	}

	data, err := json.Marshal(ack)
	if err != nil {
		bret = false
	}
	datalen := len(data)
	data = data[:datalen]

	resphear := &CmdHeadInfo{
		Ver:    header.Ver,
		Length: uint16(datalen),
		Ttid:   binary.BigEndian.Uint32(jsutils.Uint32ToBytes(header.Ttid)),
		Ssrc:   binary.BigEndian.Uint32(jsutils.Uint32ToBytes(header.Ssrc)),
		Status: status,
		Flag:   0x00,
		Magic:  binary.BigEndian.Uint16(jsutils.Uint16ToBytes(header.Magic)),
	}

	buf := &bytes.Buffer{}
	eee := binary.Write(buf, binary.BigEndian, resphear)
	if eee != nil {
		panic(eee)
	}

	head := buf.Bytes()[:]

	bret = true
	head = append(head, data...)
	return head, len(head), bret
}

func DecodeCloseAck(instr string) (string, string, int) {
	var from, to, temp string
	var code int
	var ok bool

	from, ok = jsutils.GetJsonStr(instr, ESCC_MSG_PAR_FROM)
	to, ok = jsutils.GetJsonStr(instr, ESCC_MSG_PAR_TO)
	temp, ok = jsutils.GetJsonStr(instr, ESCC_MSG_PAR_CODE)
	if ok == true {
		code, _ = strconv.Atoi(temp)
	}

	return from, to, code
}

func EncodeUpdate(cSeq int, from string, to string, slot int, ttid int32, ssrc int32) ([]byte, int, bool) {
	var bret bool = false

	var status uint16 = 0
	status = binary.BigEndian.Uint16(jsutils.Uint16ToBytes(uint16(status)))

	ack := map[string]interface{}{
		ESCC_MSG_PAR_PROTO:      SIM_PROTO,
		ESCC_MSG_PAR_MSG:        ESCC_MSG_UPDATE,
		ESCC_MSG_PAR_CSEQ:       cSeq,
		ESCC_MSG_PAR_FROM:       from,
		ESCC_MSG_PAR_TO:         to,
		ESCC_MSG_PAR_TYPE:       "KEEPALIVE",
		ESCC_MSG_PAR_EXPIRES:    120,
		ESCC_MSG_PAR_SIM_NUM:    "",
		ESCC_MSG_PAR_SIM_SIGNAL: 26,
		ESCC_MSG_PAR_SIM_BAL:    -1,
	}

	data, err := json.Marshal(ack)
	if err != nil {
		bret = false
	}
	datalen := len(data)
	data = data[:datalen]

	resphear := &CmdHeadInfo{
		Ver:    0x10,
		Length: uint16(datalen),
		Ttid:   uint32(ttid),
		Ssrc:   binary.BigEndian.Uint32(jsutils.Uint32ToBytes(uint32(ssrc))),
		Status: status,
		Flag:   0x40,
		Magic:  0xc5a3,
	}

	buf := &bytes.Buffer{}
	eee := binary.Write(buf, binary.BigEndian, resphear)
	if eee != nil {
		panic(eee)
	}

	head := buf.Bytes()[:]

	bret = true
	head = append(head, data...)
	return head, len(head), bret
}

func EncodeClose(cSeq int, from string, to string, username string, slot int, ttid int32, ssrc int32) ([]byte, int, bool) {
	var bret bool = false

	var status uint16 = 0
	status = binary.BigEndian.Uint16(jsutils.Uint16ToBytes(uint16(status)))
	simslot := fmt.Sprintf("%s.%03d", username, slot)
	ack := map[string]interface{}{
		ESCC_MSG_PAR_PROTO:    SIM_PROTO,
		ESCC_MSG_PAR_MSG:      ESCC_MSG_CLOSE,
		ESCC_MSG_PAR_CSEQ:     cSeq,
		ESCC_MSG_PAR_FROM:     from,
		ESCC_MSG_PAR_TO:       to,
		ESCC_MSG_PAR_SIM_SLOT: simslot,
	}

	data, err := json.Marshal(ack)
	if err != nil {
		bret = false
	}
	datalen := len(data)
	data = data[:datalen]

	resphear := &CmdHeadInfo{
		Ver:    0x10,
		Length: uint16(datalen),
		Ttid:   uint32(ttid),
		Ssrc:   binary.BigEndian.Uint32(jsutils.Uint32ToBytes(uint32(ssrc))),
		Status: status,
		Flag:   0x40,
		Magic:  0xc5a3,
	}

	buf := &bytes.Buffer{}
	eee := binary.Write(buf, binary.BigEndian, resphear)
	if eee != nil {
		panic(eee)
	}

	head := buf.Bytes()[:]
	bret = true
	head = append(head, data...)
	return head, len(head), bret
}

func DecodeUpdateAck(instr string) (string, string, int) {
	var from, to, temp string
	var code int
	var ok bool

	from, ok = jsutils.GetJsonStr(instr, ESCC_MSG_PAR_FROM)
	to, ok = jsutils.GetJsonStr(instr, ESCC_MSG_PAR_TO)
	temp, ok = jsutils.GetJsonStr(instr, ESCC_MSG_PAR_CODE)
	if ok == true {
		code, _ = strconv.Atoi(temp)
	}

	return from, to, code
}

func EncodeInfoReset(from string, to string, seq int32, slot int, ttid int32, ssrc int32) ([]byte, int, bool) {
	var bret bool = false

	var tt = [...]byte{0x11, 0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	tt[0] = byte(slot)

	tt1 := tt[:]
	dst, _ := jsutils.Base64Encode(tt1, len(tt1))

	//var dst = "BwIAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA"
	ack := map[string]interface{}{
		ESCC_MSG_PAR_PROTO:       SIM_PROTO,
		ESCC_MSG_PAR_MSG:         ESCC_MSG_INFO,
		ESCC_MSG_PAR_CSEQ:        seq,
		ESCC_MSG_PAR_FROM:        from,
		ESCC_MSG_PAR_TO:          to,
		ESCC_MSG_PAR_TYPE:        ESCC_MSG_VAL_CMD,
		ESCC_MSG_PAR_SMS_CHARSET: "B64",
		ESCC_MSG_PAR_DATA:        string(dst[:]),
	}

	data, err := json.Marshal(ack)
	if err != nil {
		bret = false
	}
	datalen := len(data)
	data = data[:datalen]

	resphear := &CmdHeadInfo{
		Ver:    0x10,
		Length: uint16(datalen),
		Ttid:   uint32(ttid),
		Ssrc:   binary.BigEndian.Uint32(jsutils.Uint32ToBytes(uint32(ssrc))),
		Status: 0,
		Flag:   0x40,
		Magic:  0xc5a3,
	}
	buf := &bytes.Buffer{}
	eee := binary.Write(buf, binary.BigEndian, resphear)
	if eee != nil {
		panic(eee)
	}

	head := buf.Bytes()[:]
	bret = true
	head = append(head, data...)

	return head, len(head), bret
}
