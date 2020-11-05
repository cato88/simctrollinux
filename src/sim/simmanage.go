package sim

import (
	"fmt"
	"inf"
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
	return udpserver.UdpServerInit(ip,cmdport,dataport,ctrol)

}


/******************************************************************
* 功能说明：释放卡池控制函数
* 参数说明：无
* 返 回 值：0--成功  !0--失败
******************************************************************/
func ReleaseSimCtrol()  {
	udpserver.UdpCmdServerRelease()
	udpserver.UdpDataServerRelease()
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

	fmt.Printf("slot=%d imsi=%s randstr=%s autn=%s clientip=%s\n",slot,imsi,randstr,autn,clientip)

	cmdClientId,_ := udpserver.GetCmdClientIdByClientIp(clientip)
	if cmdClientId >0{
		simType,ok := proc.GetSimType(uint32(slot), uint32(cmdClientId))
		if ok == false {
			fmt.Printf("slot=%d imsi=%s randstr=%s autn=%s clientip=%s proc.GetSimType error\n",slot,imsi,randstr,autn,clientip)
			return -1
		}

		dataClientId,ok:= proc.GetDataClient(uint32(slot), uint32(cmdClientId))
		if ok == false {
			fmt.Printf("slot=%d imsi=%s randstr=%s autn=%s clientip=%s proc.GetDataClient error\n",slot,imsi,randstr,autn,clientip)
			return -2
		}

		simInfo := proc.GetSimInfo(uint32(slot), uint32(cmdClientId))
		simInfo.AuthErrorCount++

		if simInfo.AuthErrorCount >= proc.MAX_AUTH_ERROR_FOR_CLOSE_COUNT {
			fmt.Printf("slot=%d imsi=%s randstr=%s autn=%s clientip=%s AuthErrorCount>=%d error\n",slot,imsi,randstr,autn,clientip,proc.MAX_AUTH_ERROR_FOR_CLOSE_COUNT)
			return -3
		}

		sres,ik,ck,result,ret  := proc.SetRandInfo(uint32(slot), uint32(cmdClientId),randstr,autn)

		if ret <0 {	//失败
			return -5
		}else if ret ==1 {//重复鉴权,有结果
			fmt.Printf("slot=%d imsi=%s randstr=%s autn=%s clientip=%s ok retturn result=%d sres=%s ck=%s ik=%s\n",slot,imsi,randstr,autn,clientip,int(result),sres,ck,ik)
			udpserver.GSimCtrlInf.AuthResult(slot,imsi, int(result),sres,ck,ik, clientip)
			simInfo.AuthErrorCount = 0
			return 0

		}else if ret == 2 {	//重复鉴权,无结果
			fmt.Printf("slot=%d imsi=%s randstr=%s autn=%s clientip=%s error, repeat auth and no answer\n",slot,imsi,randstr,autn,clientip)
			return 0
		}
		//新鉴权
		dataClientInfo,ok := udpserver.FindDataClientInfo(int32(dataClientId))
		if ok != true{
			fmt.Printf("slot=%d imsi=%s randstr=%s autn=%s clientip=%s dataClientId=%d FindDataClientInfo error\n",slot,imsi,randstr,autn,clientip,int32(dataClientId))
			return -6
		}


AUTH2:	if simType == inf.SIM_TYPE_2G {
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
				goto AUTH2
				return 0
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
			}
		}
	}else{
		fmt.Printf("slot=%d imsi=%s randstr=%s autn=%s clientip=%s error, not found clientid by clientip\n",slot,imsi,randstr,autn,clientip)
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

		return simInfo.Imsi, int(simInfo.SimType)
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

		cmdClientInfo,ok1 := udpserver.FindCmdClientInfo(int32(cmdClientId))
		if ok1 == false {
			return -1
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
		}
	}

	return 1
}
