package proc

const DEFAULT_UPDATE_TIME_INTERVAL = 20

const MAX_AUTH_ERROR_FOR_CLOSE_COUNT = 5
const CMD_STATE_NULL = 0
const CMD_STATE_IMSI_1 = 1
const CMD_STATE_IMSI_2 = 2
const CMD_STATE_AUTH_1 = 3
const CMD_STATE_AUTH_2 = 4

const DATA_STATE_NULL = 0
const DATA_STATE_WAITING = 0

const SIM_PROTO = "ESPP/2.0"
const LOCAL_FROM = 123456

const CMD_CLIENT_MAX_COUNT = 100
const DATA_CLIENT_MAX_COUNT = 2000
const SLOT_MAX_COUNT = 1025

const ESCC_MSG_PROTO = "ESPP/1.1"

const ESCC_MSG_REG = "REG"
const ESCC_MSG_OPEN = "OPEN"
const ESCC_MSG_CLOSE = "CLOSE"
const ESCC_MSG_UPDATE = "UPDATE"
const ESCC_MSG_INFO = "INFO"
const ESCC_MSG_MESSAGE = "MESSAGE"

const ESCC_MSG_REG_ACK = "REG-ACK"
const ESCC_MSG_OPEN_ACK = "OPEN-ACK"
const ESCC_MSG_CLOSE_ACK = "CLOSE-ACK"
const ESCC_MSG_UPDATE_ACK = "UPDATE-ACK"
const ESCC_MSG_INFO_ACK = "INFO-ACK"
const ESCC_MSG_MSG_ACK = "MESSAGE-ACK"
const ESCC_MSG_PUBLISH = "PUBLISH"
const ESCC_MSG_PUBLISH_ACK = "PUBLISH-ACK"

/* common parameters for normal request and response */
const ESCC_MSG_PAR_PROTO = "PROTO"
const ESCC_MSG_PAR_MSG = "MSG"
const ESCC_MSG_PAR_CSEQ = "CSEQ"
const ESCC_MSG_PAR_TYPE = "TYPE"
const ESCC_MSG_PAR_CHARSET = "CHARSET"
const ESCC_MSG_PAR_DATA = "DATA"
const ESCC_MSG_PAR_CODE = "CODE"
const ESCC_MSG_PAR_REASON = "REASON"
const ESCC_MSG_PAR_STATUS = "STATUS"

const ESCC_MSG_PAR_CTT_TYPE = "CONTENT-TYPE"
const ESCC_MSG_PAR_CTT_LEN = "CONTENT-LENGTH"
const ESCC_MSG_PAR_CTT_ENC = "CONTENT-ENCODING"

/* parameters for REG and REG_ACK */
const ESCC_MSG_PAR_DID = "DID"
const ESCC_MSG_PAR_NAME = "NAME"
const ESCC_MSG_PAR_EXPIRES = "EXPIRES"
const ESCC_MSG_PAR_VER = "VERSION"
const ESCC_MSG_PAR_MAX_PORTS = "MAX-PORTS"
const ESCC_MSG_PAR_MAC = "MAC"
const ESCC_MSG_PAR_IP = "IP"
const ESCC_MSG_PAR_CHN_GRPS = "CHN-GRPS"

const ESCC_MSG_PAR_USERNAME = "USERNAME"
const ESCC_MSG_PAR_NONCE = "NONCE"
const ESCC_MSG_PAR_CNONCE = "CNONCE"
const ESCC_MSG_PAR_NC = "NC"
const ESCC_MSG_PAR_RESP = "RESPONSE"

/* parameters for OPEN/UPDATE message */
const ESCC_MSG_PAR_FROM = "FROM"
const ESCC_MSG_PAR_TO = "TO"
const ESCC_MSG_PAR_CHN_ID = "CHN-ID"
const ESCC_MSG_PAR_MOD_TYPE = "MOD-TYPE"
const ESCC_MSG_PAR_MOD_VER = "MOD-VER"
const ESCC_MSG_PAR_MOD_IMEI = "IMEI"
const ESCC_MSG_PAR_GET = "GET"

const ESCC_MSG_PAR_GOIP_SLOT = "GOIP-SLOT"
const ESCC_MSG_PAR_SIM_SLOT = "SIM-SLOT"
const ESCC_MSG_PAR_RTMODE = "ROUTE-MODE"

/* parameters for OPEN_ACK */
const ESCC_MSG_PAR_CONN_PROTO = "CONN-PROTO" /* UDP|TCP, default is udp */
const ESCC_MSG_PAR_CONN_IP = "CONN-IP"
const ESCC_MSG_PAR_CONN_PORT = "CONN-PORT"
const ESCC_MSG_PAR_CONN_ACKTIME = "CONN-ACKTIME"
const ESCC_MSG_PAR_CONN_RETRANSTIME = "CONN-RETRANSTIME"
const ESCC_MSG_PAR_CONN_RETRANSINTVL = "CONN-RETRANSINTVL"
const ESCC_MSG_PAR_CONN_RETRANSCOUNT = "CONN-RETRANSCOUNT"
const ESCC_MSG_PAR_RETRY_AFTER = "RETRY-AFTER"
const ESCC_MSG_PAR_HOT_NUM = "HOT-NUM"

const ESCC_MSG_PAR_REQUIRE = "REQUIRE"
const ESCC_MSG_PAR_SUPPORTED = "SUPPORTED"

/* SIM card information parameters */
const ESCC_MSG_PAR_SIM_ICCID = "ICCID"
const ESCC_MSG_PAR_SIM_IMSI = "IMSI"
const ESCC_MSG_PAR_SIM_IMEI = "BINDED-IMEI"
const ESCC_MSG_PAR_SIM_PROVIDER = "PROVIDER" /* "id name" */
const ESCC_MSG_PAR_SIM_NUM = "SIM-NUM"
const ESCC_MSG_PAR_SIM_STATUS = "SIM-STATUS"
const ESCC_MSG_PAR_SIM_REASON = "SIM-REASON"
const ESCC_MSG_PAR_SIM_SIGNAL = "SIM_SIG"
const ESCC_MSG_PAR_SIM_BAL = "SIM-BAL"
const ESCC_MSG_PAR_SIM_BLKTIME = "BLOCK-TIME"
const ESCC_MSG_PAR_SIM_LEDACT = "LED-ACTION"
const ESCC_MSG_PAR_SIM_DATA = "SIM-DATA"

/* CDR information parameters */
const ESCC_MSG_PAR_CDR_CID = "CID" /* call-id = sip_call_id */
const ESCC_MSG_PAR_CDR_DIR = "DIR"
const ESCC_MSG_PAR_CDR_CALLER = "CALLER"
const ESCC_MSG_PAR_CDR_CALLEE = "CALLEE"
const ESCC_MSG_PAR_CDR_BEGIN = "BEGIN" /* the call start time in seconds */
const ESCC_MSG_PAR_CDR_ALERT = "ALERT"
const ESCC_MSG_PAR_CDR_ANSWER = "ANSWER"
const ESCC_MSG_PAR_CDR_UPDATE = "UPDATE"
const ESCC_MSG_PAR_CDR_END = "END"
const ESCC_MSG_PAR_CDR_HUPRSN = "REASON"
const ESCC_MSG_PAR_CDR_IMEI = "IMEI"
const ESCC_MSG_PAR_CDR_ICCID = "ICCID"
const ESCC_MSG_PAR_CDR_IMSI = "IMSI"
const ESCC_MSG_PAR_CDR_SIMNUM = "SIM-NUM"

/* SMS information parameters */
const ESCC_MSG_PAR_SMS_ID = "SMSID" /* for SMS_REQ/SMS_QUERY/SMS_RESP */
const ESCC_MSG_PAR_SMS_TIME = "TIME"
const ESCC_MSG_PAR_SMS_FROM = "SENDER"
const ESCC_MSG_PAR_SMS_TO = "RECVER"
const ESCC_MSG_PAR_SMS_SMSC = "SMSC"
const ESCC_MSG_PAR_SMS_SCTS = "SCTS"
const ESCC_MSG_PAR_SMS_CHARSET = "CHARSET" /* encoded as UTF-8 */
const ESCC_MSG_PAR_SMS_CONTENT = "CONTENT" /* encoded as UTF-8 */
const ESCC_MSG_PAR_SMS_SIMNUM = "SIM_NUM"
const ESCC_MSG_PAR_SMS_IMEI = "IMEI"
const ESCC_MSG_PAR_SMS_ICCID = "ICCID"
const ESCC_MSG_PAR_SMS_IMSI = "IMSI"

/* SMS information parameters */
const ESCC_MSG_PAR_USSD_TIME = "TIME"
const ESCC_MSG_PAR_USSD_CONTENT = "CONTENT" /* encoded as UTF-8 */
const ESCC_MSG_PAR_USSD_SIMNUM = "SIM-NUM"
const ESCC_MSG_PAR_USSD_IMEI = "IMEI"
const ESCC_MSG_PAR_USSD_ICCID = "ICCID"
const ESCC_MSG_PAR_USSD_IMSI = "IMSI"

/* the values of param type */
const ESCC_MSG_VAL_SIM = "SIM"
const ESCC_MSG_VAL_CDR = "CDR"
const ESCC_MSG_VAL_SMS = "SMS"             /* DEV -> SVR */
const ESCC_MSG_VAL_SMS_REQ = "SMS-REQ"     /* SVR -> DEV, need to wait response */
const ESCC_MSG_VAL_SMS_QUERY = "SMS-QUERY" /* SVR -> DEV, need to wait response */
const ESCC_MSG_VAL_SMS_RESP = "SMS-RESP"   /* DEV -> SVR, response for SMS_REQ/SMS_QUERY */
const ESCC_MSG_VAL_USSD_REQ = "USSD-REQ"   /* SVR -> DEV, need to wait response */
const ESCC_MSG_VAL_USSD_RESP = "USSD-RESP" /* DEV -> SVR, response for USSD_REQ */
const ESCC_MSG_VAL_CMD = "CMD"             /* DEV <-> SVR <-> DEV, user request message */
const ESCC_MSG_VAL_CMD_RESP = "CMD_RESP"   /* DEV <-> SVR <-> DEV, user response message */

const ESCC_MSG_VAL_MT = "MT" /* the call is Terminated at Mobile */
const ESCC_MSG_VAL_MO = "MO" /* the call is Originated from Mobile */

/* chksum: means the data need to do udp chksum */
const ESCC_MSG_VAL_REQUIRE = "chksum"
const ESCC_MSG_VAL_SUPPORTED = "chksum"

const ESCC_VAL_CONN_MPORT = "CONN-MPORT"
const ESCC_VAL_MOD_VALUE = "MOD-VALUE"

const ESCC_VAL_HAS_SIM = "HAS-SIM"

/* the code difinition for parameter ESCC_MSG_PAR_CODE */
const ESCC_MSG_CODE_OK = 200

const ESCC_MSG_CODE_BAD_REQUEST = 400
const ESCC_MSG_CODE_UNAUTHORIZED = 401
const ESCC_MSG_CODE_FORBIDDEN = 403
const ESCC_MSG_CODE_NOT_FOUND = 404
const ESCC_MSG_CODE_REQ_TIMEOUT = 408
const ESCC_MSG_CODE_GONE = 410
const ESCC_MSG_CODE_TEMP_UNAVAIL = 480
const ESCC_MSG_CODE_TR_NOT_EXIST = 481
const ESCC_MSG_CODE_BUSY_HERE = 486
const ESCC_MSG_CODE_REQ_PENDING = 491

const ESCC_MSG_CODE_SERVER_ERROR = 500
const ESCC_MSG_CODE_NOT_IMPLEMENTED = 501
const ESCC_MSG_CODE_SERVICE_UNVAIL = 503
const ESCC_MSG_CODE_SERVER_TIMEOUT = 508

type CmdHeadInfo struct {
	Ver    uint8
	Flag   uint8
	Magic  uint16
	Ssrc   uint32
	Ttid   uint32
	Status uint16
	Length uint16
}

type DataReqHead struct {
	Index      uint16
	CmdLen     byte
	ParamLen   byte
	ExpRespLen uint16
	PrefixNum  uint16
	PrefixFid  [4][2]byte
}

type DataRespHead struct {
	Index uint16
	Len   uint16
	Temp  [6]byte
}

const MAX_SAVE_RANDOM_COUNT = 10
