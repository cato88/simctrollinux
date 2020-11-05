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





//ar gCmdUdpAddrMap sync.Map
//var gCmdUdpClientMap sync.Map
var gCmdUdpClientMap = &ClientIdCLientInfoMap{mutex:new(sync.RWMutex),mp:make(map[int]*UdpClientInfo)}
var gCmdUdpAddrMap = &AddrClientIdMap{mutex:new(sync.RWMutex),mp:make(map[string]int)}

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

func SetCmdClientInfoSrc(clientid int32,src int32) bool {
	clientInfo,ok :=gCmdUdpClientMap.MapClientInfoGet(int(clientid))
	if ok == true {
		clientInfo.Src = src
		return true
	}
	return false
}

func SetCmdClientInfoTid(clientid int32,tid int32) bool {
	clientInfo,ok :=gCmdUdpClientMap.MapClientInfoGet(int(clientid))
	if ok == true {
		clientInfo.Tid = uint32(tid)
		return true
	}
	return false
}

func SetCmdClientInfoUserName(clientid int32,username string) bool {
	clientInfo,ok :=gCmdUdpClientMap.MapClientInfoGet(int(clientid))
	if ok == true {
		clientInfo.UserName = username
		return true
	}
	return false
}

func SetCmdClientInfoSeq(clientid int32,seq uint32) bool {
	clientInfo,ok :=gCmdUdpClientMap.MapClientInfoGet(int(clientid))
	if ok == true {
		clientInfo.Seq = uint16(int16(seq))
		return true
	}
	return false
}

func FindCmdClientInfo(clientid int32) (*UdpClientInfo,bool) {
	return gCmdUdpClientMap.MapClientInfoGet(int(clientid))
}

func AddCmdFifo(cmdClientId int,instr []byte)  {
	gCmdUdpSendFifo.PutEntryFifo((jsutils.NewEntryFifo(int32(cmdClientId),instr)))
}

func GetCmdClientIdByClientIp(cmdClientIp string) (int,bool){
	var cmdClientId int = 0
	gCmdUdpClientMap.mutex.RLock()
	defer gCmdUdpClientMap.mutex.RUnlock()
	for _,v := range gCmdUdpClientMap.mp{
		if v.ClientIp ==cmdClientIp{
			return int(v.Clientid), true
		}
	}

	return cmdClientId,false
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
			fmt.Printf("CmdUdpExecProcess fifo msg.key=%d msg.valuses.len=%d\n",msg.Key, len(msg.Values))
			CmdUdpExec(&msg)
		}else{
			time.Sleep(100*time.Millisecond)
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
		fmt.Printf("CmdUdpServerRecv conn.ReadFromUDP len=%d addr=%v \n",len,clientAddr)
		//ss := jsutils.DisplayHexString(buf,3)
		//fmt.Println(ss)
		clientid,ok :=gCmdUdpAddrMap.MapClientIdGet(clientAddr.String())
		if ok != true {
			clientid = int(getCmdUdpClientPos())
			gCmdUdpAddrMap.MapClientIdSet(clientAddr.String(),clientid)
			ip,port:=jsutils.GetIpPort(clientAddr.String())

			clientinfo :=UdpClientInfo{Addrstr:clientAddr.String(),Conn:clientAddr,Clientid:int32(clientid),LastTime:time.Now().Unix(),ClientIp: ip,ClientPort: int16(port),Seq: 1,Tid: 1}
			gCmdUdpClientMap.MapClientInfoSet(clientid,&clientinfo)

			fmt.Println("!!! CmdUdpServerRecv add ",clientinfo)
			gCmdUdpRecvFifo.PutEntryFifo(jsutils.NewEntryFifo(int32(clientid),buf))

			gCmdSimCtrlInf.SimClientConnectNotify(clientid,ip,port)
		}else {
			clientinfo,ok1 :=FindCmdClientInfo(int32(clientid))
			if ok1 == true{
				clientinfo.LastTime = time.Now().Unix();
				fmt.Printf("CmdUdpServerRecv mutil client_id = %v  addr=%s\n",clientid,clientAddr.String())
				gCmdUdpRecvFifo.PutEntryFifo(jsutils.NewEntryFifo(int32(clientid),buf))
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
			time.Sleep(100*time.Millisecond)
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

	gCmdUdpClientMap.mutex.RLock()
	defer gCmdUdpClientMap.mutex.RUnlock()
	for key,v := range gCmdUdpClientMap.mp{
		if time.Now().Unix() >=  (v.LastTime+gCmdClientKeepActiveTime){
			ip,port:=jsutils.GetIpPort(v.Addrstr)
			gCmdSimCtrlInf.SimClientDisConnectNotify(int(v.Clientid),ip,port)
			gCmdClientIdArr[key] = 0
			delete(gCmdUdpAddrMap.mp, v.Addrstr)
			delete(gCmdUdpClientMap.mp,key )
			fmt.Printf("CmdTimeout delete key=%v value=%v\n",key,v)
		}
	}
}
/*
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
*/
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
		fmt.Println("switch retstr",retstr)
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


	ssrc := header.Ssrc


	clientinfo,ok1 := FindCmdClientInfo(clientid)
	if ok1 == true{

		SetCmdClientInfoSrc(clientid,int32(ssrc))
		SetCmdClientInfoUserName(clientid,username)

		regack,regacklen,ok2 :=proc.EncodeRegAck(header,seq,did)
		if ok2 == true{
			regack = regack[:regacklen]
			gCmdUdpSendFifo.PutEntryFifo(jsutils.NewEntryFifo(clientid,regack))
		}

	}else {
		return false
	}

	clientinfo,_ = FindCmdClientInfo(clientid)



	fmt.Println("DealReg1 cseq",clientinfo.Seq,"ttid",clientinfo.Tid)

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

				jsutils.GetNext32(&clientinfo.Tid)

				jsutils.GetNext16(&clientinfo.Seq)


				fmt.Println("DealReg2 cseq",clientinfo.Seq,"ttid",clientinfo.Tid)
				openstr,openlen,ok3 :=proc.EncodeOpen(header,gDataUdpServerIp,gDataUdpServerPort, uint32(clientinfo.Seq),uint32(clientinfo.Tid),int(slot),username)
				if ok3 {
					openstr = openstr[:openlen]
					gCmdUdpSendFifo.PutEntryFifo(jsutils.NewEntryFifo(clientid,openstr))
					gCmdSimCtrlInf.SimState(int(slot),imsi,inf.SIM_STATE_DISABLE,int(simtype),clientinfo.ClientIp)
				}
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


		}else{
			return false
		}
	}else{
		return false
	}
	openackack,openackacklen,ok:= proc.EncodeOpenAckAck(header)
	if ok == true{
		openackack = openackack[:openackacklen]
		gCmdUdpSendFifo.PutEntryFifo(jsutils.NewEntryFifo(clientid,openackack))
	}
	dataclientid,_:=findDataClientByRemoteIpPort(remoteip, remoteport)
	proc.SetDataClient(uint32(dataclientid),uint32(slot), uint32(clientid))

	cmdClientInfo,ok := FindCmdClientInfo(int32(clientid))
	if ok == false{
		return false
	}
	gCmdSimCtrlInf.SimState(int(slot),imsi,inf.SIM_STATE_USABLE,int(simtype),cmdClientInfo.ClientIp)
	fmt.Println("DealOpenAck cseq",cmdClientInfo.Seq,"ttid",cmdClientInfo.Tid)
	jsutils.GetNext32(&cmdClientInfo.Tid)
	jsutils.GetNext16(&cmdClientInfo.Seq)
	fmt.Println("DealOpenAck EncodeInfoReset cseq",cmdClientInfo.Seq,"ttid",cmdClientInfo.Tid)
	retstr,_,ok3 :=proc.EncodeInfoReset(from,to, int32(cmdClientInfo.Seq),slot, int32(cmdClientInfo.Tid),cmdClientInfo.Src)
	if ok3 == true{
		AddCmdFifo(int(clientid),retstr)
		return true
	}

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
		if ok == false{
			return false
		}
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

		cmdClientInfo,err := FindCmdClientInfo(int32(clientid))
		if err == false{
			return false
		}
		cmdClientIp := cmdClientInfo.ClientIp
		proc.SetState(inf.SIM_STATE_DISABLE,slot, uint32(clientid))

		gCmdSimCtrlInf.SimState(int(slot),imsi,inf.SIM_STATE_DISABLE,int(simtype),cmdClientIp)

		fmt.Println("DealCloseAck cseq",cmdClientInfo.Seq,"ttid",cmdClientInfo.Tid)

		jsutils.GetNext32(&cmdClientInfo.Tid)
		jsutils.GetNext16(&cmdClientInfo.Seq)

		fmt.Println("DealCloseAck EncodeOpen cseq",cmdClientInfo.Seq,"ttid",cmdClientInfo.Tid)
		//reopen
		openstr,openlen,ok3 :=proc.EncodeOpen(header,gDataUdpServerIp,gDataUdpServerPort, uint32(cmdClientInfo.Seq),uint32(cmdClientInfo.Tid),int(slot),cmdClientInfo.UserName)
		if ok3 {
			openstr = openstr[:openlen]
			gCmdUdpSendFifo.PutEntryFifo(jsutils.NewEntryFifo(clientid,openstr))
		}

	}
	return true
}

func DealUpdateAck(clientid int32,header *proc.CmdHeadInfo,instr string) bool {
	from,to,code := proc.DecodeUpdateAck(instr)
	if code == 200{
		return true
	}

	slot,ok := proc.FindSimslotByFromTo(uint32(clientid),from,to)
	if ok == true {

		siminfo := proc.GetSimInfo(slot, uint32(clientid))
		if siminfo == nil{
			return false
		}

		imsi := siminfo.Imsi
		simtype := siminfo.SimType
		lststate := siminfo.State
		if lststate == inf.SIM_STATE_USABLE {
			dataClientId := siminfo.DataClient
			CloseDataClient(int(dataClientId))
			siminfo.State = inf.SIM_STATE_DISABLE
			cmdClientInfo,ok := FindCmdClientInfo(clientid)
			if ok == true {
				gCmdSimCtrlInf.SimState(int(slot),imsi,inf.SIM_STATE_DISABLE,int(simtype), cmdClientInfo.ClientIp)
				jsutils.GetNext32(&cmdClientInfo.Tid)
				jsutils.GetNext16(&cmdClientInfo.Seq)

				//reopen
				openstr,openlen,ok3 :=proc.EncodeOpen(header,gDataUdpServerIp,gDataUdpServerPort, uint32(cmdClientInfo.Seq),uint32(cmdClientInfo.Tid),int(slot),cmdClientInfo.UserName)
				if ok3 {
					openstr = openstr[:openlen]
					gCmdUdpSendFifo.PutEntryFifo(jsutils.NewEntryFifo(clientid,openstr))
				}
			}

		}


	}
	return true
}

func CheckUpdate()  {
	curTime := uint64(time.Now().Unix())
	for m:=0;m<proc.CMD_CLIENT_MAX_COUNT;m++{
		cmdClientId := m
		cmdClientInfo,err := FindCmdClientInfo(int32(cmdClientId))
		if err == false {
			continue
		}

		for n:=0;n<proc.SLOT_MAX_COUNT;n++{
			if proc.GSimInfoList[m][n].State == 0 {
				continue
			}

			if err == false{
				continue
			}




			if proc.GSimInfoList[m][n].UpdateTime >0 && proc.GSimInfoList[m][n].UpdateTime+proc.DEFAULT_UPDATE_TIME_INTERVAL<=curTime {
				proc.GSimInfoList[m][n].UpdateTime = curTime
				fmt.Println("CheckUpdate cseq",cmdClientInfo.Seq,"ttid",cmdClientInfo.Tid)
				jsutils.GetNext32(&cmdClientInfo.Tid)
				jsutils.GetNext16(&cmdClientInfo.Seq)
				fmt.Println("CheckUpdate 2 cseq",cmdClientInfo.Seq,"ttid",cmdClientInfo.Tid)
				updatestr,_,_:= proc.EncodeUpdate(int(cmdClientInfo.Seq),proc.GSimInfoList[m][n].From,proc.GSimInfoList[m][n].To,n, int32(cmdClientInfo.Tid),cmdClientInfo.Src)
				gCmdUdpSendFifo.PutEntryFifo(jsutils.NewEntryFifo(int32(cmdClientId),updatestr))
			}

			if proc.GSimInfoList[m][n].AuthErrorCount>=proc.MAX_AUTH_ERROR_FOR_CLOSE_COUNT{
				proc.GSimInfoList[m][n].AuthErrorCount = 0

				fmt.Println("CheckUpdate 3 cseq",cmdClientInfo.Seq,"ttid",cmdClientInfo.Tid)
				jsutils.GetNext32(&cmdClientInfo.Tid)
				jsutils.GetNext16(&cmdClientInfo.Seq)
				fmt.Println("CheckUpdate 4 cseq",cmdClientInfo.Seq,"ttid",cmdClientInfo.Tid)
				clsoestr,_,_:= proc.EncodeClose(int(cmdClientInfo.Seq),proc.GSimInfoList[m][n].From,proc.GSimInfoList[m][n].To,n, int32(cmdClientInfo.Tid),cmdClientInfo.Src)
				gCmdUdpSendFifo.PutEntryFifo(jsutils.NewEntryFifo(int32(cmdClientId),clsoestr))
			}
		}
	}
}


