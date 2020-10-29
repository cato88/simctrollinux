package inf

//sim卡状态定义
const SIM_STATE_IDLE = 0		//无卡
const SIM_STATE_DISABLE = 1		//有卡不可用
const SIM_STATE_USABLE = 2		//有卡可用


//鉴权接口定义
const AUTH_RESULT_SUC = 0 	//鉴权成功
const AUTH_RESULT_FAIL = 1	//鉴权失败


//卡类型定义
const SIM_TYPE_2G = 0		//2G卡
const SIM_TYPE_4G = 1			//4G卡


//回调接口定义
type SimCtroler interface {

	/******************************************************************
	* 功能说明：卡池连接上报
	* 参数说明：
	*			clientid[IN]	- 卡池硬件客户端句柄编号
	*			clientip[IN]	- 卡池硬件客户端地址
	*			clientport[IN]	- 卡池硬件客户端端口
	* 返 回 值：无
	 */
	SimClientConnectNotify(clientid int,clientip string,clientport int)

	/******************************************************************
	* 功能说明：卡池断开连接上报
	* 参数说明：
	*			clientid[IN]	- 卡池硬件客户端句柄编号
	*			clientip[IN]	- 卡池硬件客户端地址
	*			clientport[IN]	- 卡池硬件客户端端口
	* 返 回 值：无
	 */
	SimClientDisConnectNotify (clientid int,clientip string,clientport int)

	/******************************************************************
	* 功能说明：卡池卡状态报告
	* 参数说明：
	*			slot[IN]	- 卡池槽位编号
	*			imsi[IN]	- imsi
	*			state[IN]	- 卡状态，参看sim卡状态定义
	*			simtype[IN]	- 卡类型，参看sim卡类型定义
	*			clientip[IN]- 客户端地址
	* 返 回 值：无
	 */
	SimState(slot int,imsi string,state int,simtype int,clientip string)

	/******************************************************************
	* 功能说明：鉴权结果报告
	* 参数说明：
	*			slot[IN]	- 卡池槽位编号
	*			imsi[IN]	- imsi
	*			result[IN]	- 鉴权结果，参看鉴权结果定义
	*			sres[IN]	- sres值
	*			kc[IN]		- kc值
	*			ik[IN]		- ik值
	*			clientip[IN]- 客户端地址
	* 返 回 值：无
	 */
	AuthResult(slot int,imsi string,result int,sres string,kc string,ik string,clientip string)
}

