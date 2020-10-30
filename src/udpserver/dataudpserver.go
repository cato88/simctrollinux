
package udpserver

import (
	"encoding/binary"
	"fmt"
	"inf"
	"jsutils"
	"net"
	"os"
	"proc"
	"sync"
	"time"
	"unsafe"
)

var gDataUdpAddrMap sync.Map
var gDataUdpClientMap sync.Map

var gDataUdpServerIp string
var gDataUdpServerPort uint16

var gDataUdpClientPos int32 =  1
var gDataUdpExit int = 0

var gDataClientKeepActiveTime int64 = 300

var gDataUdpRecvFifo *jsutils.LimitedEntryFifo
var gDataUdpSendFifo *jsutils.LimitedEntryFifo

var gDataClientIdArr  [proc.DATA_CLIENT_MAX_COUNT+1] int32

var gDataSimCtrlInf inf.SimCtroler

///////////////////////////////////////////////////////////

func AddDataFifo(dataClientId int,instr []byte)  {
	gDataUdpSendFifo.PutEntryFifo((jsutils.NewEntryFifo(int32(dataClientId),instr)))
}

func getDataUdpClientPos()  int32 {
	var ret int32 = 0
	for n:=gDataUdpClientPos;n<proc.DATA_CLIENT_MAX_COUNT;n++ {
		if gDataClientIdArr[n] == 0{
			gDataClientIdArr[n] = 1
			ret = gDataUdpClientPos
			gDataUdpClientPos++
			if gDataUdpClientPos> proc.DATA_CLIENT_MAX_COUNT {
				gDataUdpClientPos = 1
			}
			return ret
		}

	}
	for n:=1;n<int(gDataUdpClientPos);n++ {
		if gDataClientIdArr[n] == 0{
			gDataClientIdArr[n] = 1
			ret = gDataUdpClientPos
			gDataUdpClientPos++
			if gDataUdpClientPos> proc.DATA_CLIENT_MAX_COUNT {
				gDataUdpClientPos = 1
			}
			return ret
		}
	}

	return ret
}


func FindDataClientInfo(clientid int32) (*UdpClientInfo,bool) {
	clientInfo,ok :=gDataUdpClientMap.Load(clientid)
	if ok == true {
		cc := clientInfo.(UdpClientInfo)
		return &cc,true
	}
	return nil,false
}

func findDataClientByRemoteIpPort(remoteIp string,remotePort int) int{
	var remoteClientAddr = fmt.Sprintf("%s:%d",remoteIp,remotePort)
	var retId int
	gDataUdpClientMap.Range(func(k,v interface{}) bool{
		if (v.(UdpClientInfo).Addrstr == remoteClientAddr){
			fmt.Printf("findDataClientByRemoteIpPort  key=%v value=%v\n",k,v)
			retId = int(v.(UdpClientInfo).Clientid)
			return true
		}
		return false
	})
	return retId
}


func DataUdpExecProcess()  {

	t1 := time.Now()
	for{
		if gCmdUdpExit != 0 {
			return
		}

		t2 := time.Now()
		if t2.Sub(t1) >1 {
			DataTimeOut()
			t1 = t2
		}

		msg,ok := gDataUdpRecvFifo.GetEntryFifo()

		if ok == true{
			fmt.Printf("DataUdpExecProcess fifo msg.key=%d msg.valuses=%s\n",msg.Key,msg.Values)
			DataUdpExec(&msg)
		}else{
			time.Sleep(1)
		}
	}
}
func DataTimeOut(){
	gDataUdpClientMap.Range(func(k,v interface{}) bool{
		if time.Now().Unix() >= (v.(UdpClientInfo).LastTime +gDataClientKeepActiveTime){
			ip,port:=jsutils.GetIpPort(v.(UdpClientInfo).Addrstr)
			gDataSimCtrlInf.SimClientDisConnectNotify(int(v.(UdpClientInfo).Clientid),ip,port)
			gDataUdpClientMap.Delete(k)
			gDataUdpAddrMap.Delete(v.(UdpClientInfo).Addrstr)
			nClientId,_ :=k.(int32)
			gDataClientIdArr[nClientId] = 0
			fmt.Printf("DataTimeOut delete key=%v value=%v\n",k,v)
			return true
		}
		return true
	})
}



func DataUdpExec(msg *jsutils.EntryFifo)  {
	dataClientId := msg.Key
	data := msg.Values
	var header *proc.DataRespHead = *(**proc.DataRespHead)(unsafe.Pointer(&data))


	datalen := len(data)
	if datalen<=12 {
		return
	}

	dataClientInfo,ok:= FindCmdClientInfo(int32(dataClientId))
	if ok == false{
		return
	}

	cmdClientId,ok := proc.GetCmdClientByDataclient(uint32(dataClientId))
	if ok == false{
		return
	}

	cmdClientInfo,ok:= FindCmdClientInfo(int32(cmdClientId))
	if ok == false{
		return
	}
	cmdClientip := cmdClientInfo.ClientIp

	slot,ok := proc.GetSlotByCmdClientDataclient(cmdClientId, uint32(cmdClientId))
	if ok == false{
		return
	}

	dataStepState,ok:= proc.GetDataStepState(slot,cmdClientId)
	if ok == false{
		return
	}


	simType,ok:= proc.GetSimType(slot,cmdClientId)
	if ok == false{
		return
	}

	imsi,ok:= proc.GetImsi(slot,cmdClientId)
	if ok == false{
		return
	}

	switch dataStepState {
	case proc.CMD_STATE_IMSI_1:
		if dataClientInfo.Seq == binary.BigEndian.Uint16(jsutils.Uint16ToBytes(header.Index)) {
			dataClientInfo.Seq++
			if dataClientInfo.Seq == 0 {
				dataClientInfo.Seq = 1
			}
			if simType == inf.SIM_TYPE_2G {
				temp, ok := proc.EncodeImsi2GStep2(dataClientInfo.Seq)
				if ok == true {
					gDataUdpSendFifo.PutEntryFifo(jsutils.NewEntryFifo(dataClientId, temp))
					proc.SetDataStepState(proc.CMD_STATE_IMSI_2, slot, cmdClientId)
				}
			} else if simType == inf.SIM_TYPE_4G {
				temp, ok := proc.EncodeImsi4GStep2(dataClientInfo.Seq)
				if ok == true {
					gDataUdpSendFifo.PutEntryFifo(jsutils.NewEntryFifo(dataClientId, temp))
					proc.SetDataStepState(proc.CMD_STATE_IMSI_2, slot, cmdClientId)
				}
			}
		}

	case proc.CMD_STATE_IMSI_2:
		if dataClientInfo.Seq == binary.BigEndian.Uint16(jsutils.Uint16ToBytes(header.Index)) {
			proc.SetDataStepState(proc.CMD_STATE_NULL, slot, cmdClientId)
		}
	case proc.CMD_STATE_AUTH_1:
		if simType == inf.SIM_TYPE_2G {
			if dataClientInfo.Seq == binary.BigEndian.Uint16(jsutils.Uint16ToBytes(header.Index)) {
				result, sres, kc, ret := proc.DecodeAuth2GStep1(data)
				if ret == 0 {
					proc.SetDataStepState(proc.CMD_STATE_NULL, slot, cmdClientId)
					if result == 0 {
						gDataSimCtrlInf.AuthResult(int(slot), imsi, result, sres, kc, "", cmdClientip)
						proc.SetErroCount(0, slot, cmdClientId)
					}else{
						proc.SetErroCount(proc.MAX_AUTH_ERROR_FOR_CLOSE_COUNT, slot, cmdClientId)
					}
				} else if ret == 99 {
					dataClientInfo.Seq++
					if dataClientInfo.Seq == 0 {
						dataClientInfo.Seq = 1
					}
					temp, ok := proc.EncodeAuth2GStep2(dataClientInfo.Seq)
					if ok == true {
						gDataUdpSendFifo.PutEntryFifo(jsutils.NewEntryFifo(dataClientId, temp))
						proc.SetDataStepState(proc.CMD_STATE_AUTH_2, slot, cmdClientId)
						proc.SetDataState(proc.DATA_STATE_WAITING, slot, cmdClientId)
					}
				} else {
					proc.SetDataStepState(proc.CMD_STATE_NULL, slot, cmdClientId)
					proc.SetErroCount(proc.MAX_AUTH_ERROR_FOR_CLOSE_COUNT, slot, cmdClientId)
				}
			}
		} else if simType == inf.SIM_TYPE_4G {
			if dataClientInfo.Seq == binary.BigEndian.Uint16(jsutils.Uint16ToBytes(header.Index)) {
				result, sres, ik,kc, ret := proc.DecodeAuth4GStep1(data)
				if ret == 0 {
					proc.SetDataStepState(proc.CMD_STATE_NULL, slot, cmdClientId)
					if result == 0 {
						gDataSimCtrlInf.AuthResult(int(slot), imsi, result, sres, kc, ik, cmdClientip)
						proc.SetErroCount(0, slot, cmdClientId)
					}else{
						proc.SetErroCount(proc.MAX_AUTH_ERROR_FOR_CLOSE_COUNT, slot, cmdClientId)
					}
				} else if ret == 99 {
					dataClientInfo.Seq++
					if dataClientInfo.Seq == 0 {
						dataClientInfo.Seq = 1
					}
					temp, ok := proc.EncodeAuth4GStep2(dataClientInfo.Seq)
					if ok == true {
						gDataUdpSendFifo.PutEntryFifo(jsutils.NewEntryFifo(dataClientId, temp))
						proc.SetDataStepState(proc.CMD_STATE_AUTH_2, slot, cmdClientId)
						proc.SetDataState(proc.DATA_STATE_WAITING, slot, cmdClientId)
					}
				} else {
					proc.SetDataStepState(proc.CMD_STATE_NULL, slot, cmdClientId)
					proc.SetErroCount(proc.MAX_AUTH_ERROR_FOR_CLOSE_COUNT, slot, cmdClientId)
				}
			}
		}

	case proc.CMD_STATE_AUTH_2:
		if simType == inf.SIM_TYPE_2G {
			if dataClientInfo.Seq == binary.BigEndian.Uint16(jsutils.Uint16ToBytes(header.Index)) {
				result, sres, kc, ret := proc.DecodeAuth2GStep2(data)
				if ret == 0 {
					proc.SetDataStepState(proc.CMD_STATE_NULL, slot, cmdClientId)
					if result == 0 {
						gDataSimCtrlInf.AuthResult(int(slot), imsi, result, sres, kc, "", cmdClientip)
						proc.SetErroCount(0, slot, cmdClientId)
					}else{
						proc.SetErroCount(proc.MAX_AUTH_ERROR_FOR_CLOSE_COUNT, slot, cmdClientId)
					}
				} else {
					proc.SetDataStepState(proc.CMD_STATE_NULL, slot, cmdClientId)
					proc.SetErroCount(proc.MAX_AUTH_ERROR_FOR_CLOSE_COUNT, slot, cmdClientId)
				}
			}
		} else if simType == inf.SIM_TYPE_4G {
			if dataClientInfo.Seq == binary.BigEndian.Uint16(jsutils.Uint16ToBytes(header.Index)) {
				result, sres,ik, kc, ret := proc.DecodeAuth4GStep2(data)
				if ret == 0 {
					proc.SetDataStepState(proc.CMD_STATE_NULL, slot, cmdClientId)
					if result == 0 {
						gDataSimCtrlInf.AuthResult(int(slot), imsi, result, sres, kc, ik, cmdClientip)
						proc.SetErroCount(0, slot, cmdClientId)
					}else{
						proc.SetErroCount(proc.MAX_AUTH_ERROR_FOR_CLOSE_COUNT, slot, cmdClientId)
					}
				} else {
					proc.SetDataStepState(proc.CMD_STATE_NULL, slot, cmdClientId)
					proc.SetErroCount(proc.MAX_AUTH_ERROR_FOR_CLOSE_COUNT, slot, cmdClientId)
				}
			}
		}
	default:


	}
}

func CloseDataClient(clientId int) bool {
	udpInfo,ok := FindDataClientInfo(int32(clientId))
	if ok == true{
		gDataUdpAddrMap.Delete(udpInfo.Addrstr)
	}
	gDataUdpClientMap.Delete(clientId)
	return true
}

func DataUdpServerRecv(conn *net.UDPConn)  {

	data := make([]byte,9999)

	for{

		if gDataUdpExit != 0 {
			return
		}

		len,clientAddr,err :=conn.ReadFromUDP(data)

		if err != nil {
			continue
		}
		buf := [] byte(data[:len])
		fmt.Printf("DataUdpServerRecv conn.ReadFromUDP len=%d addr=%v  data=%s\n",len,clientAddr,buf)
		//ss := jsutils.DisplayHexString(buf,3)
		//fmt.Println(ss)

		clientid,ok :=gDataUdpAddrMap.Load(clientAddr.String())
		if ok != true {
			clientid = getDataUdpClientPos();
			gDataUdpAddrMap.Store(clientAddr.String(),clientid)
			nClientId,_ :=clientid.(int32)
			clientinfo :=UdpClientInfo{Addrstr:clientAddr.String(),Conn:clientAddr,Clientid:nClientId,LastTime:time.Now().Unix()}
			gDataUdpClientMap.Store(clientid,clientinfo)

			fmt.Printf("!!! DataUdpServerRecv new client_id = %v  addr=%s\n",clientid,clientAddr.String())
			gDataUdpRecvFifo.PutEntryFifo(jsutils.NewEntryFifo(nClientId ,buf))
		}else {
			nClientId,_ :=clientid.(int32)
			clientinfo,ok1 := FindDataClientInfo(nClientId)
			if ok1 == true{
				clientinfo.LastTime = time.Now().Unix();
				fmt.Printf("DataUdpServerRecv mutil client_id = %v  addr=%s\n",clientid,clientAddr.String())
				gDataUdpRecvFifo.PutEntryFifo(jsutils.NewEntryFifo(nClientId,buf))
			}
		}
	}

	defer conn.Close()
}

func DataUdpServerSend(conn *net.UDPConn)  {

	for{

		if gDataUdpExit != 0 {
			return
		}

		msg,ok := gDataUdpSendFifo.GetEntryFifo()
		if ok == true{
			fmt.Printf("DataUdpServerSend fifo msg.key=%d msg.valuses=%s\n",msg.Key,msg.Values)
			clientinfo,ok := FindDataClientInfo(msg.Key)
			if ok  == true {
				_,err := conn.WriteToUDP([]byte(msg.Values),clientinfo.Conn)
				if err != nil{
					fmt.Println("DataUdpServerSend WriteToUDP err="+err.Error() )
				}
			}

		}else{
			time.Sleep(1)
		}
	}
}

func UdpDataServerInit(ip string,cmdport int,ctrol inf.SimCtroler) int {

	var addrStr string
	addrStr = fmt.Sprintf("%s:%d",ip,cmdport)
	gDataUdpServerIp = ip
	gDataUdpServerPort = uint16(cmdport)
	fmt.Println("UdpDataServerInit addrInfo = "+addrStr)
	gDataSimCtrlInf = ctrol
	udpAddr,err := net.ResolveUDPAddr("udp",addrStr)
	if err != nil{
		fmt.Println("UdpDataServerInit ResolveUDPAddr err="+err.Error())
		os.Exit(1)
	}

	conn,err:=net.ListenUDP("udp",udpAddr)
	if err != nil{
		fmt.Println("UdpDataServerInit ListenUDP err="+err.Error())
		os.Exit(2)
	}



	gDataUdpRecvFifo =jsutils.NewLimitedEntryFifo(1000)
	gDataUdpSendFifo =jsutils.NewLimitedEntryFifo(1000)

	go DataUdpServerRecv(conn)
	go DataUdpServerSend(conn)
	go DataUdpExecProcess()

	return 0
}


func UdpDataServerRelease()  {
	gDataUdpExit = 1
}



