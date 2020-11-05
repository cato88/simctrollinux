package udpserver

import (
	"inf"
	"net"
	"sync"
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
	Tid				uint32
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




type ClientIdCLientInfoMap struct{
	mutex *sync.RWMutex
	mp map[int]*UdpClientInfo
}

func (this *ClientIdCLientInfoMap) MapClientInfoGet(key int)  (*UdpClientInfo,bool){
	this.mutex.RLock()
	ret,ok:=this.mp[key]
	this.mutex.RUnlock()
	if ok == false{
		return nil,false
	}
	return ret,true
}

func (this *ClientIdCLientInfoMap) MapClientInfoSet(key int,info *UdpClientInfo)  bool{
	this.mutex.RLock()
	defer this.mutex.RUnlock()
	this.mp[key] = info
	return true
}



type AddrClientIdMap struct {
	mutex *sync.RWMutex
	mp map[string]int
}


func (this *AddrClientIdMap) MapClientIdGet(key string)  (int,bool){
	this.mutex.RLock()
	ret,ok:=this.mp[key]
	this.mutex.RUnlock()
	if ok == false{
		return 0,false
	}
	return ret,true
}

func (this *AddrClientIdMap) MapClientIdSet(key string,info int)  bool{
	this.mutex.RLock()
	defer this.mutex.RUnlock()
	this.mp[key] = info
	return true
}
