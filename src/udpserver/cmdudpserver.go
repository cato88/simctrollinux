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
	"unsafe"
)



////////////////////////////////////////////////////////////////////////////////////////////////////////////////

var gCmdUdpExit int = 0

var gCmdUdpAddrMap sync.Map
var gCmdUdpClientMap sync.Map

var gCmdUdpClientPos int32 =  1

var gCmdClientKeepActiveTime int64 = 300

var gCmdUdpRecvFifo *jsutils.LimitedEntryFifo
var gCmdUdpSendFifo *jsutils.LimitedEntryFifo

var gCmdUdpServerIp string
var gCmdUdpServerPort uint16

var gCmdClientIdArr  [proc.CMD_CLIENT_MAX_COUNT+1] int32

var gCmdSimCtrlInf inf.SimCtroler

///////////////////////////////////////////////////////////

func getCmdUdpClientPos()  int32 {

	var ret int32 = 0
	for n:=gCmdUdpClientPos;n<proc.CMD_CLIENT_MAX_COUNT;n++ {
		if gCmdClientIdArr[n] == 0{
			gCmdClientIdArr[n] = 1
			ret = gCmdUdpClientPos
			gCmdUdpClientPos++
			if gCmdUdpClientPos> proc.CMD_CLIENT_MAX_COUNT {
				gCmdUdpClientPos = 1
			}
			return ret
		}

	}
	for n:=1;n<int(gCmdUdpClientPos);n++ {
		if gCmdClientIdArr[n] == 0{
			gCmdClientIdArr[n] = 1
			ret = gCmdUdpClientPos
			gCmdUdpClientPos++
			if gCmdUdpClientPos> proc.CMD_CLIENT_MAX_COUNT {
				gCmdUdpClientPos = 1
			}
			return ret
		}
	}

	return ret
}


func FindCmdClientInfo(clientid int32) (*UdpClientInfo,bool) {
	clientInfo,ok :=gCmdUdpClientMap.Load(clientid)
	if ok == true {
		cc := clientInfo.(UdpClientInfo)
		return &cc,true
	}
	return nil,false
}

func AddCmdFifo(cmdClientId int,instr []byte)  {
	gCmdUdpSendFifo.PutEntryFifo((jsutils.NewEntryFifo(int32(cmdClientId),instr)))
}

func GetCmdClientIdByClientIp(cmdClientIp string) (int,bool){
	var cmdClientId int = 0
	gCmdUdpClientMap.Range(func(k,v interface{}) bool{
		if v.(UdpClientInfo).ClientIp ==cmdClientIp {

			cmdClientId = int(v.(UdpClientInfo).Clientid)
			return true
		}
		return true
	})

	return cmdClientId,true
}


func CmdUdpExecProcess()  {
	t1 := time.Now()
	for{
		if gCmdUdpExit != 0 {
			return
		}

		t2 := time.Now()
 		if t2.Sub(t1) >1 {
			CmdTimeout()
			CheckUpdate()
			t1 = t2
		}

		msg,ok := gCmdUdpRecvFifo.GetEntryFifo()
		if ok == true{
			fmt.Printf("CmdUdpExecProcess fifo msg.key=%d msg.valuses=%s\n",msg.Key,msg.Values)
			CmdUdpExec(&msg)
		}else{
			time.Sleep(1)
		}
	}
}


func CmdUdpServerRecv(conn *net.UDPConn)  {

	data := make([]byte,9999)

	for{
		if gCmdUdpExit != 0 {
			return
		}

		len,clientAddr,err :=conn.ReadFromUDP(data)

		if err != nil {
			continue
		}
		buf := [] byte(data[:len])
		fmt.Printf("CmdUdpServerRecv conn.ReadFromUDP len=%d addr=%v  data=%s\n",len,clientAddr,buf)
		//ss := jsutils.DisplayHexString(buf,3)
		//fmt.Println(ss)

		clientid,ok :=gCmdUdpAddrMap.Load(clientAddr.String())
		if ok != true {
			clientid = getCmdUdpClientPos();
			gCmdUdpAddrMap.Store(clientAddr.String(),clientid)

			ip,port:=jsutils.GetIpPort(clientAddr.String())
			nClientId,_ :=clientid.(int32)
			clientinfo :=UdpClientInfo{Addrstr:clientAddr.String(),Conn:clientAddr,Clientid:nClientId,LastTime:time.Now().Unix(),ClientIp: ip,ClientPort: int16(port)}
			gCmdUdpClientMap.Store(clientid,clientinfo)

			fmt.Printf("!!! CmdUdpServerRecv new client_id = %v  addr=%s\n",clientid,clientAddr.String())
			gCmdUdpRecvFifo.PutEntryFifo(jsutils.NewEntryFifo(nClientId ,buf))

			gCmdSimCtrlInf.SimClientConnectNotify(int(nClientId),ip,port)
		}else {
			nClientId,_ :=clientid.(int32)
			clientinfo,ok1 := FindCmdClientInfo(nClientId)
			if ok1 == true{
				clientinfo.LastTime = time.Now().Unix();
				fmt.Printf("CmdUdpServerRecv mutil client_id = %v  addr=%s\n",clientid,clientAddr.String())
				gCmdUdpRecvFifo.PutEntryFifo(jsutils.NewEntryFifo(nClientId,buf))
			}
		}
	}

	defer conn.Close()
}

func CmdUdpServerSend(conn *net.UDPConn)  {

	for{
		if gCmdUdpExit != 0 {
			return
		}

		msg,ok := gCmdUdpSendFifo.GetEntryFifo()
		if ok == true{
			fmt.Printf("CmdUdpServerSend fifo msg.key=%d msg.valuses=%s\n",msg.Key,msg.Values)
			clientinfo,ok := FindCmdClientInfo(msg.Key)
			if ok  == true {
				_,err := conn.WriteToUDP([]byte(msg.Values),clientinfo.Conn)
				if err != nil{
					fmt.Println("CmdUdpServerSend WriteToUDP err="+err.Error() )
				}
			}

		}else{
			time.Sleep(1)
		}
	}
}

func UdpCmdServerInit(ip string,cmdport int,ctrol inf.SimCtroler) int {

	var addrStr string
	addrStr = fmt.Sprintf("%s:%d",ip,cmdport)
	gCmdUdpServerIp = ip
	gCmdUdpServerPort = uint16(cmdport)
	fmt.Println("UdpServerInit addrInfo = "+addrStr)
	gCmdSimCtrlInf = ctrol

	udpAddr,err := net.ResolveUDPAddr("udp",addrStr)
	if err != nil{
		fmt.Println("UdpServerInit ResolveUDPAddr err="+err.Error())
		os.Exit(1)
	}

	conn,err:=net.ListenUDP("udp",udpAddr)
	if err != nil{
		fmt.Println("UdpServerInit ListenUDP err="+err.Error())
		os.Exit(2)
	}


	gCmdUdpRecvFifo =jsutils.NewLimitedEntryFifo(1000)
	gCmdUdpSendFifo =jsutils.NewLimitedEntryFifo(1000)

	go CmdUdpServerRecv(conn)
	go CmdUdpServerSend(conn)
	go CmdUdpExecProcess()


	return 0
}

func UdpCmdServerRelease()  {
	gCmdUdpExit = 1
}

func CmdTimeout()  {
	gCmdUdpClientMap.Range(func(k,v interface{}) bool{
		if time.Now().Unix() >= (v.(UdpClientInfo).LastTime +gCmdClientKeepActiveTime){

			ip,port:=jsutils.GetIpPort(v.(UdpClientInfo).Addrstr)

			gCmdSimCtrlInf.SimClientDisConnectNotify(int(v.(UdpClientInfo).Clientid),ip,port)
			nClientId,_ :=k.(int32)
			gDataClientIdArr[nClientId] = 0
			gCmdUdpClientMap.Delete(k)
			gCmdUdpAddrMap.Delete(v.(UdpClientInfo).Addrstr)
			fmt.Printf("CmdTimeout delete key=%v value=%v\n",k,v)
			return true
		}
		return true
	})
}

func CmdUdpExec(msg *jsutils.EntryFifo)  {
	clientid := msg.Key
	data := msg.Values
	var header *proc.CmdHeadInfo = *(**proc.CmdHeadInfo)(unsafe.Pointer(&data))
	flag,status,bret := proc.IsNeedAck(header)
	if bret == true {
		buf,len,ok := proc.EncodeAck(header,flag,status)
		if ok == true{
			buf1 :=buf[:len]
			gCmdUdpRecvFifo.PutEntryFifo(jsutils.NewEntryFifo(clientid,buf1))
		}
	}

	datalen := len(data)
	if datalen<=16 {
		return
	}

	msgdata := string(data[16:])

	retstr,ok := jsutils.GetJsonStr(msgdata,proc.ESCC_MSG_PAR_MSG)
	if ok != true{
		return
	}

	switch retstr {
	case proc.ESCC_MSG_REG:
		DealReg(clientid,header,msgdata)
	case proc.ESCC_MSG_OPEN_ACK:
		DealOpenAck(clientid,header,msgdata)
	case proc.ESCC_MSG_CLOSE:
		DealClose(clientid,header,msgdata)
	case proc.ESCC_MSG_PUBLISH:
		DealPublish(clientid,header,msgdata)
	case proc.ESCC_MSG_CLOSE_ACK:
		DealCloseAck(clientid,header,msgdata)
	case proc.ESCC_MSG_UPDATE_ACK:
		DealUpdateAck(clientid,header,msgdata)
	case proc.ESCC_MSG_INFO_ACK:
	default:

	}


	fmt.Println(clientid ,header)

}

func DealReg(clientid int32,header *proc.CmdHeadInfo,instr string) bool {
	seq,did,simregstate,username,ok := proc.DecodeReg(instr)
	if ok != true{
		return false
	}

	clientinfo,ok1 := FindCmdClientInfo(clientid)
	if ok1 == true{
		clientinfo.Src = int32(header.Ssrc)
		clientinfo.Tid = int32(header.Ttid)
		clientinfo.UserName = username

		regack,regacklen,ok2 :=proc.EncodeRegAck(header,seq,did)
		if ok2 == true{
			regack = regack[:regacklen]
			gCmdUdpSendFifo.PutEntryFifo(jsutils.NewEntryFifo(clientid,regack))
		}

	}else {
		return false
	}

	for i:=0;i<len(simregstate);i++{

		var slot uint32
		slot = uint32(i+1)

		simtype ,_ := proc.GetSimType(slot,uint32(clientid))
		imsi ,_ := proc.GetImsi(slot,uint32(clientid))
		simstate,_ := proc.GetState(slot,uint32(clientid))
		if simregstate[i] == 0x30{
			if simstate != inf.SIM_STATE_IDLE{
				proc.SetState(inf.SIM_STATE_IDLE,slot,uint32(clientid))
				gCmdSimCtrlInf.SimState(int(slot),imsi,inf.SIM_STATE_IDLE,int(simtype),clientinfo.ClientIp)
			}
		}else if simregstate[i] == 0x31 {
			proc.SetCmdClient(uint32(clientid),slot,uint32(clientid))
			if simstate == inf.SIM_STATE_IDLE || simstate == inf.SIM_STATE_DISABLE {
				proc.SetState(inf.SIM_STATE_DISABLE,slot,uint32(clientid))
				proc.SetErroCount(0,slot,uint32(clientid))
				clientinfo.Tid++
				clientinfo.Seq++
			}
			openstr,openlen,ok3 :=proc.EncodeOpen(header,gDataUdpServerIp,gDataUdpServerPort,clientinfo.Seq,uint32(clientinfo.Tid),int(slot),username)
			if ok3 {
				openstr = openstr[:openlen]
				gCmdUdpSendFifo.PutEntryFifo(jsutils.NewEntryFifo(clientid,openstr))
				gCmdSimCtrlInf.SimState(int(slot),imsi,inf.SIM_STATE_DISABLE,int(simtype),clientinfo.ClientIp)
			}
		}
	}

	return true
}

func DealOpenAck(clientid int32,header *proc.CmdHeadInfo,instr string) bool {
	imsi,imei,iccid,from,to,code,slot,_,remoteip,remoteport,simtype := proc.DecodeOpenAck(instr)
	if imsi != ""{
		if code == 200{
			proc.SetState(inf.SIM_STATE_USABLE, uint32(slot), uint32(clientid))
			proc.SetImsi(imsi,uint32(slot), uint32(clientid))
			proc.SetImei(imei,uint32(slot), uint32(clientid))
			proc.SetIccid(iccid,uint32(slot), uint32(clientid))
			proc.SetFrom(from,uint32(slot), uint32(clientid))
			proc.SetTo(to,uint32(slot), uint32(clientid))
			proc.SetSimType(uint32(simtype),uint32(slot), uint32(clientid))
			proc.SetSUpdateTime(uint64(time.Now().Unix()),uint32(slot), uint32(clientid))
		}
	}
	openackack,openackacklen,ok:= proc.EncodeOpenAckAck(header)
	if ok == true{
		openackack = openackack[:openackacklen]
		gCmdUdpSendFifo.PutEntryFifo(jsutils.NewEntryFifo(clientid,openackack))
	}

	proc.SetDataClient(uint32(findDataClientByRemoteIpPort(remoteip, remoteport)),uint32(slot), uint32(clientid))

	return true
}

func DealClose(clientid int32,header *proc.CmdHeadInfo,instr string) bool {
	from,to,cseq,slot := proc.DecodeClose(instr)
	closeack,closeacklen,ok:= proc.EncodeCloseAck(header,cseq,from,to)
	if ok == true{
		closeack = closeack[:closeacklen]
		gCmdUdpSendFifo.PutEntryFifo(jsutils.NewEntryFifo(clientid,closeack))
	}
	lastState,_ := proc.GetState(uint32(slot), uint32(clientid))
	proc.SetState(inf.SIM_STATE_DISABLE,uint32(slot), uint32(clientid))
	proc.SetErroCount(0,uint32(slot), uint32(clientid))
	if lastState == inf.SIM_STATE_USABLE {
		dataClientId, ok := proc.GetDataClient(uint32(slot), uint32(clientid))
		if ok == true{
			CloseDataClient(int(dataClientId))
		}
	}
	return true
}


func DealPublish(clientid int32,header *proc.CmdHeadInfo,instr string) bool {
	username,from,to,slot,cseq,simstate := proc.DecodePublish(instr)
	publishack,publishacklen,ok:= proc.EncodePublishAck(header,cseq,from,to,username,simstate,slot)
	if ok == true{
		publishack = publishack[:publishacklen]
		gCmdUdpSendFifo.PutEntryFifo(jsutils.NewEntryFifo(clientid,publishack))
	}

	imsi,_ := proc.GetImsi(uint32(slot), uint32(clientid))
	simtype,_ := proc.GetSimType(uint32(slot), uint32(clientid))
	proc.SetErroCount(0,uint32(slot), uint32(clientid))
	if simstate == 0{		//卡拔出
		dataClientId,ok := proc.GetDataClient(uint32(slot), uint32(clientid))
		if ok== true{
			CloseDataClient(int(dataClientId))
		}

		cmdClientInfo,ok := FindCmdClientInfo(int32(clientid))
		cmdClientIp := cmdClientInfo.ClientIp
		proc.SetState(inf.SIM_STATE_IDLE,uint32(slot), uint32(clientid))
		gCmdSimCtrlInf.SimState(int(slot),imsi,inf.SIM_STATE_IDLE,int(simtype),cmdClientIp)
	}
	return true
}

func DealCloseAck(clientid int32,header *proc.CmdHeadInfo,instr string) bool {
	from,to,_ := proc.DecodeCloseAck(instr)
	slot,ok := proc.FindSimslotByFromTo(uint32(clientid),from,to)
	if ok == true {
		proc.SetErroCount(0,slot, uint32(clientid))
		imsi,_:= proc.GetImsi(slot, uint32(clientid))
		simtype,_:= proc.GetSimType(slot, uint32(clientid))
		simstate,_:= proc.GetSimType(slot, uint32(clientid))

		if simstate == inf.SIM_STATE_USABLE{
			dataclient,_ := proc.GetDataClient(slot, uint32(clientid))
			CloseDataClient(int(dataclient))
		}

		cmdClientInfo,_ := FindCmdClientInfo(int32(clientid))
		cmdClientIp := cmdClientInfo.ClientIp
		proc.SetState(inf.SIM_STATE_DISABLE,slot, uint32(clientid))

		gCmdSimCtrlInf.SimState(int(slot),imsi,inf.SIM_STATE_DISABLE,int(simtype),cmdClientIp)

		cmdClientInfo.Seq++;
		cmdClientInfo.Tid++;
		if cmdClientInfo.Seq == 0 {
			cmdClientInfo.Seq = 1
		}

		//reopen
		openstr,openlen,ok3 :=proc.EncodeOpen(header,gDataUdpServerIp,gDataUdpServerPort,cmdClientInfo.Seq,uint32(cmdClientInfo.Tid),int(slot),cmdClientInfo.UserName)
		if ok3 {
			openstr = openstr[:openlen]
			gCmdUdpSendFifo.PutEntryFifo(jsutils.NewEntryFifo(clientid,openstr))
		}

	}
	return true
}

func DealUpdateAck(clientid int32,header *proc.CmdHeadInfo,instr string) bool {
	from,to,_ := proc.DecodeUpdateAck(instr)
	slot,ok := proc.FindSimslotByFromTo(uint32(clientid),from,to)
	if ok == true {
		proc.SetErroCount(0,slot, uint32(clientid))
		imsi,_:= proc.GetImsi(slot, uint32(clientid))
		simtype,_:= proc.GetSimType(slot, uint32(clientid))
		simstate,_:= proc.GetSimType(slot, uint32(clientid))

		if simstate == inf.SIM_STATE_USABLE{
			dataclient,_ := proc.GetDataClient(slot, uint32(clientid))
			CloseDataClient(int(dataclient))
		}

		cmdClientInfo,_ := FindCmdClientInfo(int32(clientid))
		cmdClientIp := cmdClientInfo.ClientIp
		proc.SetState(inf.SIM_STATE_DISABLE,slot, uint32(clientid))

		gCmdSimCtrlInf.SimState(int(slot),imsi,inf.SIM_STATE_DISABLE,int(simtype),cmdClientIp)

		cmdClientInfo.Seq++;
		cmdClientInfo.Tid++;
		if cmdClientInfo.Seq == 0 {
			cmdClientInfo.Seq = 1
		}

		//reopen
		openstr,openlen,ok3 :=proc.EncodeOpen(header,gDataUdpServerIp,gDataUdpServerPort,cmdClientInfo.Seq,uint32(cmdClientInfo.Tid),int(slot),cmdClientInfo.UserName)
		if ok3 {
			openstr = openstr[:openlen]
			gCmdUdpSendFifo.PutEntryFifo(jsutils.NewEntryFifo(clientid,openstr))
		}

	}
	return true
}

func CheckUpdate()  {
	curTime := uint64(time.Now().Unix())
	for m:=0;m<proc.CMD_CLIENT_MAX_COUNT;m++{
		for n:=0;n<proc.SLOT_MAX_COUNT;n++{
			if proc.GSimInfoList[m][n].State == 0 {
				continue
			}
			cmdClientId := m
			cmdClientInfo,_ := FindCmdClientInfo(int32(cmdClientId))
			cmdClientInfo.Seq++
			if cmdClientInfo.Seq == 0 {
				cmdClientInfo.Seq = 1
			}

			cmdClientInfo.Tid++


			if proc.GSimInfoList[m][n].UpdateTime >0 && proc.GSimInfoList[m][n].UpdateTime+proc.DEFAULT_UPDATE_TIME_INTERVAL<=curTime {
				proc.GSimInfoList[m][n].UpdateTime = curTime
				updatestr,_,_:= proc.EncodeUpdate(int(cmdClientInfo.Seq),proc.GSimInfoList[m][n].From,proc.GSimInfoList[m][n].To,n,cmdClientInfo.Tid,cmdClientInfo.Src)
				gCmdUdpSendFifo.PutEntryFifo(jsutils.NewEntryFifo(int32(cmdClientId),updatestr))
			}

			if proc.GSimInfoList[m][n].AuthErrorCount>=proc.MAX_AUTH_ERROR_FOR_CLOSE_COUNT{
				proc.GSimInfoList[m][n].AuthErrorCount = 0

				cmdClientInfo.Seq++
				if cmdClientInfo.Seq == 0 {
					cmdClientInfo.Seq = 1
				}

				cmdClientInfo.Tid++

				clsoestr,_,_:= proc.EncodeClose(int(cmdClientInfo.Seq),proc.GSimInfoList[m][n].From,proc.GSimInfoList[m][n].To,n,cmdClientInfo.Tid,cmdClientInfo.Src)
				gCmdUdpSendFifo.PutEntryFifo(jsutils.NewEntryFifo(int32(cmdClientId),clsoestr))
			}
		}
	}
}


