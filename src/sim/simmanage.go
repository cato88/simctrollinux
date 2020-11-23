package sim

import (
	"inf"
	"jsutils"
	"proc"
	"udpserver"
)




/******************************************************************
	* 功能说明：初始化卡池控制函数
	* 参数说明：ip[IN]	- 本地UDP服务端IP地址
	*			cmdport[IN]	- 本地UDP服务端命令控制端口
	*			dataport[IN]	- 本地UDP服务端数据控制端口
	*			ctrol[IN]		- 回调接口
	* 返 回 值：0--成功  !0--失败
*/
func InitSimControl(ip string,cmdport int,dataport int,ctrol inf.SimCtroler)  int{
	jsutils.InitLog("./log/simctrol","simctrol")
	jsutils.SetLogLeavel(jsutils.WARN)

	jsutils.Fatal("InitSimControl ip=",ip,"cmdport=",cmdport,"dataport",dataport)

	ret:= udpserver.UdpServerInit(ip,cmdport,dataport,ctrol)
	if ret == 0 {
		jsutils.Fatal("InitSimControl udpserver.UdpServerInit ok")
	}else{
		jsutils.Fatal("InitSimControl udpserver.UdpServerInit err, ret=",ret)
	}
	return ret
}


/******************************************************************
* 功能说明：释放卡池控制函数
* 参数说明：无
* 返 回 值：0--成功  !0--失败
******************************************************************/
func ReleaseSimCtrol()  {
	jsutils.Fatal("ReleaseSimCtrol start")
	udpserver.UdpCmdServerRelease()
	udpserver.UdpDataServerRelease()
	jsutils.Fatal("ReleaseSimCtrol ok")
}

/******************************************************************
* 功能说明：sim卡鉴权函数
* 参数说明：slot[IN]		- sim卡编号
*			imsi[IN]		- sim卡imsi
*			randstr[IN]	- 随机数
*			autn[IN]		- 挑战值,2g卡时此项为空，4g卡必填
*			clientip[IN]	- 卡池客户端ip地址
* 返 回 值：调用函数结果 0--成功  !0--失败  ,鉴权结果在回调接口函数AuthResult中返回，默认超时1秒
******************************************************************/
func Auth(slot int,imsi string,randstr string,autn string,clientip string) int{


	jsutils.Fatal("slot=",slot,"imsi=",imsi,"randstr=",randstr ,"autn=",autn,"clientip=",clientip)

	cmdClientId,_ := udpserver.GetCmdClientIdByClientIp(clientip)
	if cmdClientId >0{
		simType,ok := proc.GetSimType(uint32(slot), uint32(cmdClientId))
		if ok == false {
			jsutils.Fatal("slot=",slot,"imsi=",imsi,"randstr=",randstr ,"autn=",autn,"clientip=",clientip,"proc.GetSimType error")
			return -1
		}

		dataClientId,ok:= proc.GetDataClient(uint32(slot), uint32(cmdClientId))
		if ok == false {
			jsutils.Fatal("slot=",slot,"imsi=",imsi,"randstr=",randstr ,"autn=",autn,"clientip=",clientip,"proc.GetDataClient error")
			return -2
		}

		simInfo := proc.GetSimInfo(uint32(slot), uint32(cmdClientId))
		simInfo.AuthErrorCount++

		if simInfo.AuthErrorCount >= proc.MAX_AUTH_ERROR_FOR_CLOSE_COUNT {
			jsutils.Fatal("slot=",slot,"imsi=",imsi,"randstr=",randstr ,"autn=",autn,"clientip=",clientip,"error AuthErrorCount>=",proc.MAX_AUTH_ERROR_FOR_CLOSE_COUNT)
			return -3
		}

		sres,ik,ck,result,ret  := proc.SetRandInfo(uint32(slot), uint32(cmdClientId),randstr,autn)

		if ret <0 {	//失败
			return -5
		}else if ret ==1 {//重复鉴权,有结果
			jsutils.Fatal("repeat auth ,return resp ok  slot=",slot," imsi=",imsi,"randstr=",randstr,"autn=",autn,"clientip=",clientip," ok retturn result=",result,"sres=",sres,"ck=",ck,"ik=",ik)
			udpserver.GSimCtrlInf.AuthResult(slot,imsi, int(result),sres,ck,ik, clientip)
			simInfo.AuthErrorCount = 0
			return 0

		}else if ret == 2 {	//重复鉴权,无结果
			jsutils.Fatal("repeat auth... ,no resp , slot=",slot," imsi=",imsi,"randstr=",randstr,"autn=",autn,"clientip=",clientip,"error, repeat auth and no answer")
			return 0
		}
		//新鉴权
		dataClientInfo,ok := udpserver.FindDataClientInfo(int32(dataClientId))
		if ok != true{
			jsutils.Fatal("new auth ,slot=",slot," imsi=",imsi,"randstr=",randstr,"autn=",autn,"clientip=",clientip,"dataClientId=",dataClientId,"FindDataClientInfo error")
			return -6
		}


		if simType == inf.SIM_TYPE_2G {

			dataClientInfo.Seq++
			if dataClientInfo.Seq == 0{
				dataClientInfo.Seq = 1
			}
			udpserver.SetDataClientInfoSeq(int32(dataClientId), uint32(dataClientInfo.Seq))
			authdata,ok := proc.EncodeAuth2GStep1(randstr, uint16(dataClientInfo.Seq))
			if ok == true{
				udpserver.AddDataFifo(int(dataClientId),authdata)
				proc.SetDataStepState(proc.CMD_STATE_AUTH_1, uint32(slot), uint32(cmdClientId))
				proc.SetDataState(proc.DATA_STATE_WAITING, uint32(slot), uint32(cmdClientId))
				return 0
			}

		}else if simType == inf.SIM_TYPE_4G {
			if len(autn) <=0{
				jsutils.Fatal("slot=",slot," imsi=",imsi,"randstr=",randstr,"autn=",autn,"clientip=",clientip,"force 2G auth")
				dataClientInfo.Seq++
				if dataClientInfo.Seq == 0{
					dataClientInfo.Seq = 1
				}
				udpserver.SetDataClientInfoSeq(int32(dataClientId), uint32(dataClientInfo.Seq))
				authdata,ok := proc.EncodeAuth2GStep1(randstr, uint16(dataClientInfo.Seq))
				if ok == true{
					udpserver.AddDataFifo(int(dataClientId),authdata)
					proc.SetDataStepState(proc.CMD_STATE_AUTH_1, uint32(slot), uint32(cmdClientId))
					proc.SetDataState(proc.DATA_STATE_WAITING, uint32(slot), uint32(cmdClientId))
					return 0
				}else{
					jsutils.Fatal("slot=",slot," imsi=",imsi,"randstr=",randstr,"autn=",autn,"clientip=",clientip,"proc.EncodeAuth2GStep1 error")
					return -7
				}
			}

			dataClientInfo.Seq++
			if dataClientInfo.Seq == 0{
				dataClientInfo.Seq = 1
			}
			udpserver.SetDataClientInfoSeq(int32(dataClientId), uint32(dataClientInfo.Seq))

			authdata,ok := proc.EncodeAuth4GStep1(randstr,autn, uint16(dataClientInfo.Seq))
			if ok == true{
				udpserver.AddDataFifo(int(dataClientId),authdata)

				proc.SetDataStepState(proc.CMD_STATE_AUTH_1, uint32(slot), uint32(cmdClientId))
				proc.SetDataState(proc.DATA_STATE_WAITING, uint32(slot), uint32(cmdClientId))
				return 0
			}else{
				jsutils.Fatal("slot=",slot," imsi=",imsi,"randstr=",randstr,"autn=",autn,"clientip=",clientip,"proc.EncodeAuth4GStep1 error")
				return -8
			}
		}
	}else{
		jsutils.Fatal("slot=",slot," imsi=",imsi,"randstr=",randstr,"autn=",autn,"clientip=",clientip,"error, not found clientid by clientip")
	}
	return 0
}

/******************************************************************
* 功能说明：读取sim卡信息
* 参数说明：slot[IN]		- sim卡编号
*			clientip[IN]	- 卡池客户端ip地址
*			imsi[OUT]		- sim卡imsi
*			simtype[OUT]	- sim卡类型 参考SIM_TYPE定义
* 返 回 值：0--成功  !0--失败
******************************************************************/
func ReadImsi(slot int,clientip string) (imsi string,simtype int) {
	cmdClientId,ok := udpserver.GetCmdClientIdByClientIp(clientip)
	if ok == true {

		simInfo := proc.GetSimInfo(uint32(slot), uint32(cmdClientId))
		if simInfo == nil{
			jsutils.Fatal("ReadImsi proc.GetSimInfo error, slot=",slot,"clientip=",clientip)
			return "",0
		}
		jsutils.Fatal("ReadImsi slot=",slot,"clientip=",clientip,"imsi=",simInfo.Imsi,"simtype=",simInfo.SimType)
		return simInfo.Imsi, int(simInfo.SimType)
	}else{
		jsutils.Fatal("ReadImsi udpserver.GetCmdClientIdByClientIp error, slot=",slot,"clientip=",clientip)
	}

	return "",0
}

/******************************************************************
* 功能说明：初始化卡池控制函数
* 参数说明：slot[IN]			- 重置sim卡编号
*			clientip[IN]		- 卡池客户端ip地址
* 返 回 值：调用函数结果  0--成功  !0--失败  ,重置结果在回调接口函数SimState中上报卡状态
******************************************************************/

func ReSet(slot int,clientip string)  int{
	cmdClientId,ok := udpserver.GetCmdClientIdByClientIp(clientip)
	if ok == true {

		simInfo := proc.GetSimInfo(uint32(slot), uint32(cmdClientId))
		if simInfo == nil{
			jsutils.Fatal("ReSet proc.GetSimInfo error, slot=",slot,"clientip=",clientip)
			return -1
		}

		cmdClientInfo,ok1 := udpserver.FindCmdClientInfo(int32(cmdClientId))
		if ok1 == false {
			jsutils.Fatal("ReSet udpserver.FindCmdClientInfo error, slot=",slot,"clientip=",clientip,"cmdClientId=",cmdClientId)
			return -2
		}

		cmdClientInfo.Tid++
		cmdClientInfo.Seq++
		if cmdClientInfo.Seq == 0{
			cmdClientInfo.Seq = 1
		}

		from := simInfo.From
		to := simInfo.To
		retstr,_,ok3 :=proc.EncodeInfoReset(from,to, int32(cmdClientInfo.Seq),slot, int32(cmdClientInfo.Tid),cmdClientInfo.Src)
		if ok3 == true{
			udpserver.AddCmdFifo(cmdClientId,retstr)
			return 0
		}else{
			jsutils.Fatal("ReSet proc.EncodeInfoReset error, slot=",slot,"clientip=",clientip,"cmdClientId=",cmdClientId,"from=",from,"to=",to,"seq=",cmdClientInfo.Seq,"tid=",cmdClientInfo.Tid,"src=",cmdClientInfo.Src)
			return -3
		}
	}else{
		jsutils.Fatal("ReSet udpserver.GetCmdClientIdByClientIp error, slot=",slot,"clientip=",clientip)
	}

	return -4
}
