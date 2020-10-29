package udpserver

import (
	"inf"
	"net"
)

type UdpClientInfo struct {
	Addrstr 		string
	Conn			*net.UDPAddr
	Clientid		int32
	LastTime		int64

	ClientIp		string
	ClientPort		int16

	Seq				uint16
	Src				int32
	Tid				int32
	UserName		string
}

var GSimCtrlInf inf.SimCtroler

func UdpServerInit(ip string,cmdport int,dataport int,ctrol inf.SimCtroler) int {

	GSimCtrlInf = ctrol

	ret := UdpCmdServerInit(ip,cmdport,ctrol)
	if ret !=0 {
		return ret
	}
	return UdpDataServerInit(ip,dataport,ctrol)
}

