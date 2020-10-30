package proc

import (
	"inf"
	"strconv"
	"time"
)

type RandomInfo struct {
	Random			string
	Auth			string
	Sres			string
	Ck				string
	Ik				string
	Result			uint32
	StartTime		uint64
	EndTime			uint64
}

type SimInfo struct {
	From		string
	To			string
	SimSlot		string
	GoipSlot	string
	Iccid		string
	Imsi		string
	Imei		string
	UpdateTime	uint64
	State		uint32
	Seq			uint16
	SimType		uint32
	CmdClient	uint32
	DataClient	uint32
	DataState	uint32
	DataStepState	uint32
	AuthErrorCount	uint32
	RandList	[MAX_SAVE_RANDOM_COUNT]RandomInfo
	RandListPos	int
}

var GSimInfoList [CMD_CLIENT_MAX_COUNT][SLOT_MAX_COUNT]SimInfo


func SetFrom(from string,slot uint32,clientid uint32)  bool {
	if clientid>= CMD_CLIENT_MAX_COUNT {
		return false
	}
	if slot>= SLOT_MAX_COUNT {
		return false
	}
	GSimInfoList[clientid][slot].From = from
	return true
}

func GetFrom(slot uint32,clientid uint32)  (string,bool) {
	var temp string
	if clientid>= CMD_CLIENT_MAX_COUNT {
		return temp,false
	}
	if slot>= SLOT_MAX_COUNT {
		return temp,false
	}
	return GSimInfoList[clientid][slot].From,true
}


func SetTo(to string,slot uint32,clientid uint32)  bool {
	if clientid>= CMD_CLIENT_MAX_COUNT {
		return false
	}
	if slot>= SLOT_MAX_COUNT {
		return false
	}
	GSimInfoList[clientid][slot].To = to
	return true
}

func GetTo(slot uint32,clientid uint32)  (string,bool) {
	var temp string
	if clientid>= CMD_CLIENT_MAX_COUNT {
		return temp,false
	}
	if slot>= SLOT_MAX_COUNT {
		return temp,false
	}
	return GSimInfoList[clientid][slot].To,true
}


func SetSimSlot(simslot string,slot uint32,clientid uint32)  bool {
	if clientid>= CMD_CLIENT_MAX_COUNT {
		return false
	}
	if slot>= SLOT_MAX_COUNT {
		return false
	}
	GSimInfoList[clientid][slot].SimSlot = simslot
	return true
}

func GetSimSlot(slot uint32,clientid uint32)  (string,bool) {
	var temp string
	if clientid>= CMD_CLIENT_MAX_COUNT {
		return temp,false
	}
	if slot>= SLOT_MAX_COUNT {
		return temp,false
	}
	return GSimInfoList[clientid][slot].SimSlot,true
}


func SetGoipSlot(goiplot string,slot uint32,clientid uint32)  bool {
	if clientid>= CMD_CLIENT_MAX_COUNT {
		return false
	}
	if slot>= SLOT_MAX_COUNT {
		return false
	}
	GSimInfoList[clientid][slot].GoipSlot = goiplot
	return true
}

func GetGoipSlot(slot uint32,clientid uint32)  (string,bool) {
	var temp string
	if clientid>= CMD_CLIENT_MAX_COUNT {
		return temp,false
	}
	if slot>= SLOT_MAX_COUNT {
		return temp,false
	}
	return GSimInfoList[clientid][slot].GoipSlot,true
}


func SetIccid(iccid string,slot uint32,clientid uint32)  bool {
	if clientid>= CMD_CLIENT_MAX_COUNT {
		return false
	}
	if slot>= SLOT_MAX_COUNT {
		return false
	}
	GSimInfoList[clientid][slot].Iccid = iccid
	return true
}

func GetIccid(slot uint32,clientid uint32)  (string,bool) {
	var temp string
	if clientid>= CMD_CLIENT_MAX_COUNT {
		return temp,false
	}
	if slot>= SLOT_MAX_COUNT {
		return temp,false
	}
	return GSimInfoList[clientid][slot].Iccid,true
}

func SetImsi(imsi string,slot uint32,clientid uint32)  bool {
	if clientid>= CMD_CLIENT_MAX_COUNT {
		return false
	}
	if slot>= SLOT_MAX_COUNT {
		return false
	}
	GSimInfoList[clientid][slot].Imsi = imsi
	return true
}

func GetImsi(slot uint32,clientid uint32)  (string,bool) {
	var temp string
	if clientid>= CMD_CLIENT_MAX_COUNT {
		return temp,false
	}
	if slot>= SLOT_MAX_COUNT {
		return temp,false
	}
	return GSimInfoList[clientid][slot].Imsi,true
}


func SetImei(imei string,slot uint32,clientid uint32)  bool {
	if clientid>= CMD_CLIENT_MAX_COUNT {
		return false
	}
	if slot>= SLOT_MAX_COUNT {
		return false
	}
	GSimInfoList[clientid][slot].Imei = imei
	return true
}

func GetImei(slot uint32,clientid uint32)  (string,bool) {
	var temp string
	if clientid>= CMD_CLIENT_MAX_COUNT {
		return temp,false
	}
	if slot>= SLOT_MAX_COUNT {
		return temp,false
	}
	return GSimInfoList[clientid][slot].Imei,true
}


func SetState(state uint32,slot uint32,clientid uint32)  bool {
	if clientid>= CMD_CLIENT_MAX_COUNT {
		return false
	}
	if slot>= SLOT_MAX_COUNT {
		return false
	}
	GSimInfoList[clientid][slot].State = state
	return true
}

func GetState(slot uint32,clientid uint32)  (uint32,bool) {
	var temp uint32
	if clientid>= CMD_CLIENT_MAX_COUNT {
		return temp,false
	}
	if slot>= SLOT_MAX_COUNT {
		return temp,false
	}
	return GSimInfoList[clientid][slot].State,true
}


func SetSeq(seq uint16,slot uint32,clientid uint32)  bool {
	if clientid>= CMD_CLIENT_MAX_COUNT {
		return false
	}
	if slot>= SLOT_MAX_COUNT {
		return false
	}
	GSimInfoList[clientid][slot].Seq = seq
	return true
}

func GetSeq(slot uint32,clientid uint32)  (uint16,bool) {
	var temp uint16
	if clientid>= CMD_CLIENT_MAX_COUNT {
		return temp,false
	}
	if slot>= SLOT_MAX_COUNT {
		return temp,false
	}
	return GSimInfoList[clientid][slot].Seq,true
}


func SetSUpdateTime(updatetime uint64,slot uint32,clientid uint32)  bool {
	if clientid>= CMD_CLIENT_MAX_COUNT {
		return false
	}
	if slot>= SLOT_MAX_COUNT {
		return false
	}
	GSimInfoList[clientid][slot].UpdateTime = updatetime
	return true
}

func GetUpdateTime(slot uint32,clientid uint32)  (uint64,bool) {
	var temp uint64
	if clientid>= CMD_CLIENT_MAX_COUNT {
		return temp,false
	}
	if slot>= SLOT_MAX_COUNT {
		return temp,false
	}
	return GSimInfoList[clientid][slot].UpdateTime,true
}


func SetSimType(simtype uint32,slot uint32,clientid uint32)  bool {
	if clientid>= CMD_CLIENT_MAX_COUNT {
		return false
	}
	if slot>= SLOT_MAX_COUNT {
		return false
	}
	GSimInfoList[clientid][slot].SimType = simtype
	return true
}

func GetSimType(slot uint32,clientid uint32)  (uint32,bool) {
	var temp uint32
	if clientid>= CMD_CLIENT_MAX_COUNT {
		return temp,false
	}
	if slot>= SLOT_MAX_COUNT {
		return temp,false
	}
	return GSimInfoList[clientid][slot].SimType,true
}



func SetDataClient(dataclient uint32,slot uint32,clientid uint32)  bool {
	if clientid>= CMD_CLIENT_MAX_COUNT {
		return false
	}
	if slot>= SLOT_MAX_COUNT {
		return false
	}
	GSimInfoList[clientid][slot].DataClient = dataclient
	return true
}

func GetDataClient(slot uint32,clientid uint32)  (uint32,bool) {
	var temp uint32
	if clientid>= CMD_CLIENT_MAX_COUNT {
		return temp,false
	}
	if slot>= SLOT_MAX_COUNT {
		return temp,false
	}
	return GSimInfoList[clientid][slot].DataClient,true
}


func SetCmdClient(cmdclient uint32,slot uint32,clientid uint32)  bool {
	if clientid>= CMD_CLIENT_MAX_COUNT {
		return false
	}
	if slot>= SLOT_MAX_COUNT {
		return false
	}
	GSimInfoList[clientid][slot].CmdClient = cmdclient
	return true
}

func GetCmdClient(slot uint32,clientid uint32)  (uint32,bool) {
	var temp uint32
	if clientid>= CMD_CLIENT_MAX_COUNT {
		return temp,false
	}
	if slot>= SLOT_MAX_COUNT {
		return temp,false
	}
	return GSimInfoList[clientid][slot].CmdClient,true
}


func SetDataState(datastate uint32,slot uint32,clientid uint32)  bool {
	if clientid>= CMD_CLIENT_MAX_COUNT {
		return false
	}
	if slot>= SLOT_MAX_COUNT {
		return false
	}
	GSimInfoList[clientid][slot].DataState = datastate
	return true
}

func GetDataState(slot uint32,clientid uint32)  (uint32,bool) {
	var temp uint32
	if clientid>= CMD_CLIENT_MAX_COUNT {
		return temp,false
	}
	if slot>= SLOT_MAX_COUNT {
		return temp,false
	}
	return GSimInfoList[clientid][slot].DataState,true
}


func SetDataStepState(datassteptate uint32,slot uint32,clientid uint32)  bool {
	if clientid>= CMD_CLIENT_MAX_COUNT {
		return false
	}
	if slot>= SLOT_MAX_COUNT {
		return false
	}
	GSimInfoList[clientid][slot].DataStepState = datassteptate
	return true
}

func GetDataStepState(slot uint32,clientid uint32)  (uint32,bool) {
	var temp uint32
	if clientid>= CMD_CLIENT_MAX_COUNT {
		return temp,false
	}
	if slot>= SLOT_MAX_COUNT {
		return temp,false
	}
	return GSimInfoList[clientid][slot].DataStepState,true
}


func SetErroCount(errorcount uint32,slot uint32,clientid uint32)  bool {
	if clientid>= CMD_CLIENT_MAX_COUNT {
		return false
	}
	if slot>= SLOT_MAX_COUNT {
		return false
	}
	GSimInfoList[clientid][slot].AuthErrorCount = errorcount
	return true
}

func GetErrorCount(slot uint32,clientid uint32)  (uint32,bool) {
	var temp uint32
	if clientid>= CMD_CLIENT_MAX_COUNT {
		return temp,false
	}
	if slot>= SLOT_MAX_COUNT {
		return temp,false
	}
	return GSimInfoList[clientid][slot].AuthErrorCount,true
}

func GetCmdClientByDataclient(dataclientid uint32)  (uint32,bool) {

	for n:=0;n< CMD_CLIENT_MAX_COUNT;n++{
		for m :=0 ;m< SLOT_MAX_COUNT;m++{

			if GSimInfoList[n][m].DataClient == dataclientid{
				return GSimInfoList[n][m].CmdClient,true
			}
		}
	}
	return 0,false
}

func GetSlotByCmdClientDataclient(cmdclientid uint32,dataclientid uint32)  (uint32,bool) {

	for m :=0 ;m< SLOT_MAX_COUNT;m++{
		if GSimInfoList[cmdclientid][m].DataClient == dataclientid {
			return uint32(m),true
		}
	}

	return 0,false
}


//最后一个返回值 <0--失败  =0--新的鉴权  1--重复鉴权有结果  2--重复鉴权无结果
func SetRandInfo(slot uint32,clientid uint32,random string,auth string)  (string,string,string,int32,int){
	//var temp uint32
	if clientid>= CMD_CLIENT_MAX_COUNT {
		return "","","",0,-1
	}
	if slot>= SLOT_MAX_COUNT {
		return "","","",0,-2
	}

	if GSimInfoList[clientid][slot].SimType == inf.SIM_TYPE_4G && len(auth)<=0 {
		return "","","",0,-3
	}

	for n:=0;n< MAX_SAVE_RANDOM_COUNT;n++{
		if GSimInfoList[clientid][slot].RandList[n].Random == random {
			if GSimInfoList[clientid][slot].RandList[n].EndTime > 0 {
				return GSimInfoList[clientid][slot].RandList[n].Sres, GSimInfoList[clientid][slot].RandList[n].Ik, GSimInfoList[clientid][slot].RandList[n].Ck,int32(GSimInfoList[clientid][slot].RandList[n].Result),1
			}else{
				return GSimInfoList[clientid][slot].RandList[n].Sres, GSimInfoList[clientid][slot].RandList[n].Ik, GSimInfoList[clientid][slot].RandList[n].Ck,int32(GSimInfoList[clientid][slot].RandList[n].Result),2
			}
		}
	}

	GSimInfoList[clientid][slot].RandList[GSimInfoList[clientid][slot].RandListPos].Random = random
	GSimInfoList[clientid][slot].RandList[GSimInfoList[clientid][slot].RandListPos].Auth = auth
	GSimInfoList[clientid][slot].RandList[GSimInfoList[clientid][slot].RandListPos].Sres = ""
	GSimInfoList[clientid][slot].RandList[GSimInfoList[clientid][slot].RandListPos].Result = 0
	GSimInfoList[clientid][slot].RandList[GSimInfoList[clientid][slot].RandListPos].Ik = ""
	GSimInfoList[clientid][slot].RandList[GSimInfoList[clientid][slot].RandListPos].Ck = ""
	GSimInfoList[clientid][slot].RandList[GSimInfoList[clientid][slot].RandListPos].StartTime = uint64(time.Now().Unix())
	GSimInfoList[clientid][slot].RandList[GSimInfoList[clientid][slot].RandListPos].EndTime = 0

	GSimInfoList[clientid][slot].RandListPos++
	if GSimInfoList[clientid][slot].RandListPos>= MAX_SAVE_RANDOM_COUNT {
		GSimInfoList[clientid][slot].RandListPos = 0
	}

	return "","","",0,0
}

func SaveAuthResp(slot uint32,clientid uint32,sres string,ik string,ck string,result int)  bool {
	if clientid>= CMD_CLIENT_MAX_COUNT {
		return false
	}
	if slot>= SLOT_MAX_COUNT {
		return false
	}

	for n:=0;n< MAX_SAVE_RANDOM_COUNT;n++{
		if GSimInfoList[clientid][slot].RandList[n].StartTime > 0 && GSimInfoList[clientid][slot].RandList[n].EndTime == 0 {
			GSimInfoList[clientid][slot].RandList[GSimInfoList[clientid][slot].RandListPos].Sres = sres
			GSimInfoList[clientid][slot].RandList[GSimInfoList[clientid][slot].RandListPos].Ik = ik
			GSimInfoList[clientid][slot].RandList[GSimInfoList[clientid][slot].RandListPos].Ck = ck
			GSimInfoList[clientid][slot].RandList[GSimInfoList[clientid][slot].RandListPos].Result = uint32(result)
			return true
		}
	}
	return false
}

func FindSimslotByFromTo(clientid uint32,from string,to string) (uint32,bool) {
	for n:=0;n< MAX_SAVE_RANDOM_COUNT;n++{
		if GSimInfoList[clientid][n].From == from && GSimInfoList[clientid][n].To == to {
			slot,_ := strconv.Atoi(GSimInfoList[clientid][n].SimSlot)
			return uint32(slot),true
		}
	}
	return 0,false
}

func GetSimInfo(slot uint32,clientid uint32) *SimInfo {
	if clientid>= CMD_CLIENT_MAX_COUNT {
		return nil
	}
	if slot>= SLOT_MAX_COUNT {
		return nil
	}

	return &GSimInfoList[clientid][slot]
}

