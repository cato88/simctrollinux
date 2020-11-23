
package udpserver

import (
	"fmt"
	"inf"
	"jsutils"
	"net"
	"os"
	"proc"
	"sync"
	"time"
)

//var gDataUdpAddrMap sync.Map
//var gDataUdpClientMap sync.Map

var gDataUdpClientMap = &ClientIdCLientInfoMap{mutex:new(sync.RWMutex),mp:make(map[int]*UdpClientInfo)}
var gDataUdpAddrMap = &AddrClientIdMap{mutex:new(sync.RWMutex),mp:make(map[string]int)}

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
	return gDataUdpClientMap.MapClientInfoGet(int(clientid))
}

func findDataClientByRemoteIpPort(remoteIp string,remotePort int) (int,bool){
	var remoteClientAddr = fmt.Sprintf("%s:%d",remoteIp,remotePort)

	gDataUdpClientMap.mutex.RLock()
	defer gDataUdpClientMap.mutex.RUnlock()
	for _,v := range gDataUdpClientMap.mp{
		if v.Addrstr ==remoteClientAddr{
			return int(v.Clientid),true
		}
	}

	return 0,false
}

func SetDataClientInfoSeq(clientid int32,seq uint32) bool {
	clientInfo,ok := gDataUdpClientMap.MapClientInfoGet(int(clientid))
	if ok == true {
		clientInfo.Seq = uint16(seq)
		return true
	}
	return false
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
			DataUdpExec(&msg)
		}else{
			time.Sleep(100*time.Millisecond)
		}
	}
}
func DataTimeOut(){
	gDataUdpClientMap.mutex.RLock()
	defer gDataUdpClientMap.mutex.RUnlock()
	for key,v := range gDataUdpClientMap.mp{
		if time.Now().Unix() >=  (v.LastTime+gDataClientKeepActiveTime){
			gDataClientIdArr[key] = 0
			delete(gDataUdpAddrMap.mp, v.Addrstr)
			delete(gDataUdpClientMap.mp,key )
			jsutils.Fatal("---DataTimeOut delete key=",key,"value=",v)
		}
	}
}



func DataUdpExec(msg *jsutils.EntryFifo)  {
	dataClientId := msg.Key
	data := msg.Values

	datalen := len(data)
	if datalen<=12 {
		return
	}

	var rspindex uint16
	rspindex = uint16(data[0]) * 256 + uint16(data[1])


	dataClientInfo,ok:= FindDataClientInfo(int32(dataClientId))
	if ok == false{
		jsutils.Fatal("DataUdpExec FindDataClientInfo error ",dataClientId)
		return
	}


	cmdClientId,ok := proc.GetCmdClientByDataclient(uint32(dataClientId))
	if ok == false{
		jsutils.Fatal("DataUdpExec GetCmdClientByDataclient error ",dataClientId)
		return
	}

	cmdClientInfo,ok:= FindCmdClientInfo(int32(cmdClientId))
	if ok == false{
		jsutils.Fatal("DataUdpExec FindCmdClientInfo error ",dataClientId)
		return
	}

	cmdClientip := cmdClientInfo.ClientIp

	slot,ok := proc.GetSlotByCmdClientDataclient(cmdClientId, uint32(dataClientId))
	if ok == false{
		jsutils.Fatal("DataUdpExec GetSlotByCmdClientDataclient error ",cmdClientId)
		return
	}

	dataStepState,ok:= proc.GetDataStepState(slot,cmdClientId)
	if ok == false{
		jsutils.Fatal("DataUdpExec GetDataStepState error ",cmdClientId)
		return
	}


	simType,ok:= proc.GetSimType(slot,cmdClientId)
	if ok == false{
		jsutils.Fatal("DataUdpExec GetSimType error ",cmdClientId)
		return
	}

	imsi,ok:= proc.GetImsi(slot,cmdClientId)
	if ok == false{
		jsutils.Fatal("DataUdpExec GetImsi error ",cmdClientId)
		return
	}

	switch dataStepState {
	case proc.CMD_STATE_IMSI_1:
		if dataClientInfo.Seq == rspindex {
			jsutils.GetNext16(&dataClientInfo.Seq)

			if simType == inf.SIM_TYPE_2G {
				temp, ok := proc.EncodeImsi2GStep2(uint16(dataClientInfo.Seq))
				if ok == true {
					gDataUdpSendFifo.PutEntryFifo(jsutils.NewEntryFifo(dataClientId, temp))
					proc.SetDataStepState(proc.CMD_STATE_IMSI_2, slot, cmdClientId)
				}
			} else if simType == inf.SIM_TYPE_4G {
				temp, ok := proc.EncodeImsi4GStep2(uint16(dataClientInfo.Seq))
				if ok == true {
					gDataUdpSendFifo.PutEntryFifo(jsutils.NewEntryFifo(dataClientId, temp))
					proc.SetDataStepState(proc.CMD_STATE_IMSI_2, slot, cmdClientId)
				}
			}
		}

	case proc.CMD_STATE_IMSI_2:
		if dataClientInfo.Seq == rspindex {
			proc.SetDataStepState(proc.CMD_STATE_NULL, slot, cmdClientId)
		}
	case proc.CMD_STATE_AUTH_1:
		if simType == inf.SIM_TYPE_2G {
			if dataClientInfo.Seq == rspindex {
				result, sres, kc, ret := proc.DecodeAuth2GStep1(data)
				if ret == 0 {
					proc.SetDataStepState(proc.CMD_STATE_NULL, slot, cmdClientId)
					if result == 0 {
						jsutils.Fatal("2G slot=",slot,"imsi=",imsi,"result=",result,"sres=",sres,"kc=",kc,"ik= ","cmdClientip=",cmdClientip)
						gDataSimCtrlInf.AuthResult(int(slot), imsi, result, sres, kc, "", cmdClientip)
						ok = proc.SaveAuthResp(slot,cmdClientId,sres,"",kc,result)
						if ok == false{
							jsutils.Fatal("DataUdpExec CMD_STATE_AUTH_1 SaveAuthResp error ",slot,cmdClientId,sres,kc,result)
						}
						proc.SetErroCount(0, slot, cmdClientId)
					}else{
						proc.SetErroCount(proc.MAX_AUTH_ERROR_FOR_CLOSE_COUNT, slot, cmdClientId)
					}
				} else if ret == 99 {
					jsutils.GetNext16(&dataClientInfo.Seq)

					temp, ok := proc.EncodeAuth2GStep2(uint16(dataClientInfo.Seq))
					if ok == true {
						gDataUdpSendFifo.PutEntryFifo(jsutils.NewEntryFifo(dataClientId, temp))
						proc.SetDataStepState(proc.CMD_STATE_AUTH_2, slot, cmdClientId)
						proc.SetDataState(proc.DATA_STATE_WAITING, slot, cmdClientId)
					}
				} else {
					fmt.Println("DataUdpExec CMD_STATE_AUTH_1 DecodeAuth2GStep1 error ")
					proc.SetDataStepState(proc.CMD_STATE_NULL, slot, cmdClientId)
					proc.SetErroCount(proc.MAX_AUTH_ERROR_FOR_CLOSE_COUNT, slot, cmdClientId)
				}
			}else {
			}
		} else if simType == inf.SIM_TYPE_4G {
			if dataClientInfo.Seq == rspindex {
				result, sres, ik,kc, ret := proc.DecodeAuth4GStep1(data)
				if ret == 0 {
					proc.SetDataStepState(proc.CMD_STATE_NULL, slot, cmdClientId)
					if result == 0 {
						gDataSimCtrlInf.AuthResult(int(slot), imsi, result, sres, kc, ik, cmdClientip)
						proc.SaveAuthResp(slot,cmdClientId,sres,ik,kc,result)
						jsutils.Fatal("4G slot=",slot,"imsi=",imsi,"result=",result,"sres=",sres,"kc=",kc,"ik=",ik,"cmdClientip=",cmdClientip)
						proc.SetErroCount(0, slot, cmdClientId)
					}else{
						proc.SetErroCount(proc.MAX_AUTH_ERROR_FOR_CLOSE_COUNT, slot, cmdClientId)
					}
				} else if ret == 99 {

					jsutils.GetNext16(&dataClientInfo.Seq)

					temp, ok := proc.EncodeAuth4GStep2(uint16(dataClientInfo.Seq))
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
			if dataClientInfo.Seq == rspindex {
				result, sres, kc, ret := proc.DecodeAuth2GStep2(data)
				if ret == 0 {
					proc.SetDataStepState(proc.CMD_STATE_NULL, slot, cmdClientId)
					if result == 0 {
						jsutils.Fatal("2G slot=",slot,"imsi=",imsi,"result=",result,"sres=",sres,"kc=",kc,"ik=","","cmdClientip=",cmdClientip)
						gDataSimCtrlInf.AuthResult(int(slot), imsi, result, sres, kc, "", cmdClientip)
						proc.SaveAuthResp(slot,cmdClientId,sres,"",kc,result)
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
			if dataClientInfo.Seq == rspindex {
				result, sres,ik, kc, ret := proc.DecodeAuth4GStep2(data)
				if ret == 0 {
					proc.SetDataStepState(proc.CMD_STATE_NULL, slot, cmdClientId)
					if result == 0 {
						jsutils.Fatal("4G slot=",slot,"imsi=",imsi,"result=",result,"sres=",sres,"kc=",kc,"ik=",ik,"cmdClientip=",cmdClientip)
						gDataSimCtrlInf.AuthResult(int(slot), imsi, result, sres, kc, ik, cmdClientip)
						proc.SaveAuthResp(slot,cmdClientId,sres,ik,kc,result)
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
	gDataUdpClientMap.mutex.RLock()
	defer gDataUdpClientMap.mutex.RUnlock()
	if ok == true{
		delete(gDataUdpAddrMap.mp, udpInfo.Addrstr)
	}
	delete(gDataUdpClientMap.mp, clientId)
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
		//fmt.Printf("DataUdpServerRecv conn.ReadFromUDP len=%d addr=%v\n",len,clientAddr)
		//ss := jsutils.DisplayHexString(buf,3)
		//fmt.Println(ss)


		clientid,ok :=gDataUdpAddrMap.MapClientIdGet(clientAddr.String())
		if ok != true {
			clientid = int(getDataUdpClientPos())
			gDataUdpAddrMap.MapClientIdSet(clientAddr.String(),clientid)
			ip,port:=jsutils.GetIpPort(clientAddr.String())

			clientinfo :=UdpClientInfo{Addrstr:clientAddr.String(),Conn:clientAddr,Clientid:int32(clientid),LastTime:time.Now().Unix(),ClientIp: ip,ClientPort: int16(port),Seq: 1,Tid: 1}
			gDataUdpClientMap.MapClientInfoSet(clientid,&clientinfo)
			jsutils.Fatal("!!! DataUdpServerRecv add ",clientinfo)
			gDataUdpRecvFifo.PutEntryFifo(jsutils.NewEntryFifo(int32(clientid),buf))
		}else {

			clientinfo,ok1 := FindDataClientInfo(int32(clientid))
			if ok1 == true{
				clientinfo.LastTime = time.Now().Unix();
				//fmt.Printf("DataUdpServerRecv mutil client_id = %v  addr=%s\n",clientid,clientAddr.String())
				gDataUdpRecvFifo.PutEntryFifo(jsutils.NewEntryFifo(int32(clientid),buf))
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
			//fmt.Printf("DataUdpServerSend fifo msg.key=%d msg.valuses.len=%d\n",msg.Key, len(msg.Values))
			clientinfo,ok := FindDataClientInfo(msg.Key)
			if ok  == true {
				_,err := conn.WriteToUDP([]byte(msg.Values),clientinfo.Conn)
				if err != nil{
					fmt.Println("DataUdpServerSend WriteToUDP err="+err.Error() )
				}
			}

		}else{
			time.Sleep(100*time.Millisecond)
		}
	}
}

func UdpDataServerInit(ip string,dataport int,ctrol inf.SimCtroler) int {

	var addrStr string
	addrStr = fmt.Sprintf("%s:%d",ip,dataport)
	gDataUdpServerIp = ip
	gDataUdpServerPort = uint16(dataport)
	jsutils.Fatal("UdpDataServerInit addrInfo = "+addrStr)
	gDataSimCtrlInf = ctrol
	udpAddr,err := net.ResolveUDPAddr("udp",addrStr)
	if err != nil{
		jsutils.Fatal("UdpDataServerInit ResolveUDPAddr err="+err.Error())
		os.Exit(1)
	}

	conn,err:=net.ListenUDP("udp",udpAddr)
	if err != nil{
		jsutils.Fatal("UdpDataServerInit ListenUDP err="+err.Error())
		os.Exit(2)
	}



	gDataUdpRecvFifo =jsutils.NewLimitedEntryFifo(1000)
	gDataUdpSendFifo =jsutils.NewLimitedEntryFifo(1000)

	go DataUdpServerRecv(conn)
	go DataUdpServerSend(conn)
	go DataUdpExecProcess()
	jsutils.Fatal("UdpDataServerInit OK")
	return 0
}


func UdpDataServerRelease()  {
	gDataUdpExit = 1
}



