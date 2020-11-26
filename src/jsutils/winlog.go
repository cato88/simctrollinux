package jsutils

import (
	"fmt"
	"os"
	"reflect"
	"time"
)

func GetCurDayInt() int {

	return time.Now().Day()
}

func GetCurDayString() string {

	return time.Now().Format("2006-01-02")
}

func GetCurTimeString() string {

	return time.Now().Format("2006-01-02 15:04:05.000")
}

//调用os.MkdirAll递归创建文件夹
func CreateMutiDir(filePath string) error {

	if !isExist(filePath) {
		err := os.MkdirAll(filePath, os.ModePerm)
		if err != nil {
			fmt.Println("创建文件夹失败,error info:", err)
			return err
		}
		return err
	}
	return nil
}

// 判断所给路径文件/文件夹是否存在(返回true是存在)
func isExist(path string) bool {

	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

type LogInfo struct {
	gFilePath    string
	gPrefix      string
	gFileMaxSize uint32
	gCurday      uint32
	gFile        *os.File
	gLogLeavel   int //日志等级
	gLogFifo     *LimitedEntryFifo
	gIndex       int
}

const DEFAULT_FILE_MAX_SIZE = 50 << 20
const WRITE_CONTINUE__COUNT = 100

var gLogInfo = &LogInfo{}

func CreateLogFile(filename string) (*os.File, bool) {

	ret, error := os.OpenFile(filename, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0777)
	if error != nil {
		fmt.Println("CreateLogFile error", error)
		return nil, false
	}
	return ret, true
}

func WriteLog(logfile *os.File, str []byte) bool {
	if logfile != nil {
		logfile.Write(str)
		return true
	}
	return false
}

func CloseLogFile(logfile *os.File) bool {
	if logfile != nil {
		logfile.Close()
	}
	return true
}

func GetNewFileIndex(filepath string, prefix string) (index int, bret bool) {
	for n := 1; n < 10000; n++ {
		ttt := fmt.Sprintf("_%04d", n)
		retfile := filepath + "/" + prefix + GetCurDayString() + ttt + ".log"
		_, err := os.Stat(retfile)
		if os.IsNotExist(err) {
			index = n
			bret = true
			return
		}
	}
	return
}

func GetCurFileSize(filename string) (size int64, bret bool) {

	fi, err := os.Stat(filename)
	if os.IsNotExist(err) {
		size = fi.Size()
		bret = true
	}

	return
}

func GetCurLogFileName(filepath string, prefix string, index int) string {
	var ss string
	ss = fmt.Sprintf("_%04d", index)
	return filepath + "/" + prefix + GetCurDayString() + ss + ".log"
}

func GetParamString(v ...interface{}) string {

	var ret string
	for _, param := range v {
		switch reflect.TypeOf(param).Kind().String() {
		case "string":
			ret = ret + fmt.Sprintf(" %s", param)
		default:
			ret = ret + fmt.Sprintf(" %v", param)
		}
	}
	ret += "\n"
	return ret
}

func InitLog(filepath string, prefix string, filesize int) bool {

	gLogInfo.gFilePath = filepath
	ret := CreateMutiDir(filepath)
	if ret != nil {
		fmt.Println("InitLog CreateMutiDir error,", ret)
		return false
	}

	gLogInfo.gIndex, _ = GetNewFileIndex(filepath, prefix)

	tempsize := filesize
	if tempsize > 100 {
		tempsize = 100
	} else if tempsize < 10 {
		tempsize = 10
	}

	gLogInfo.gFileMaxSize = uint32(tempsize) * 1024 * 1024
	gLogInfo.gPrefix = prefix
	gLogInfo.gCurday = uint32(GetCurDayInt())
	if filesize == 0 {
		gLogInfo.gFileMaxSize = DEFAULT_FILE_MAX_SIZE
	}
	gLogInfo.gLogFifo = NewLimitedEntryFifo(10000)
	go LogProcess()
	return true
}

func LogProcess() {

	var ok bool
	defer CloseLogFile(gLogInfo.gFile)

	fmt.Println("LogProcess ...")
	curDay := uint32(GetCurDayInt())
	for {
		if gLogInfo.gLogFifo.EntryFifoLen() <= 0 {
			time.Sleep(100 * time.Millisecond)
			continue
		}

		if curDay != gLogInfo.gCurday {
			gLogInfo.gCurday = curDay
			gLogInfo.gIndex = 1
		}
		//fmt.Println("have log")

		filename := GetCurLogFileName(gLogInfo.gFilePath, gLogInfo.gPrefix, gLogInfo.gIndex)
		gLogInfo.gFile, ok = CreateLogFile(filename)
		if ok == false {
			fmt.Println("!!!LogProcess CreateLogFile error ", filename)

			for n := 0; n < WRITE_CONTINUE__COUNT; n++ {
				_, ok := gLogInfo.gLogFifo.GetEntryFifo()
				if ok == true {
				} else {
					time.Sleep(10 * time.Millisecond)
					break
				}
			}
			continue
		}

		for n := 0; n < WRITE_CONTINUE__COUNT; n++ {
			msg, ok := gLogInfo.gLogFifo.GetEntryFifo()
			if ok == true {
				WriteLog(gLogInfo.gFile, msg.Values)
			} else {
				time.Sleep(10 * time.Millisecond)
				break
			}
		}
		CloseLogFile(gLogInfo.gFile)

		filelen, ok := GetCurFileSize(filename)
		if ok {
			if filelen >= int64(gLogInfo.gFileMaxSize) {
				gLogInfo.gIndex, _ = GetNewFileIndex(gLogInfo.gFilePath, gLogInfo.gPrefix)
			}
		}
	}
}

//
//func main() {
//	InitLog("./ttt/log","sds",1)
//
//
//}

const (
	FATAL = 0
	ALERT = 1
	ERROR = 2
	WARN  = 3
	INFO  = 4
	DEBUG = 5
)

func SetLogLeavel(l int) {

	gLogInfo.gLogLeavel = l
}

func Fatal(v ...interface{}) {

	var ret string
	if gLogInfo.gLogLeavel >= FATAL {
		ret = "[" + GetCurTimeString() + "]:" + GetParamString(v)
		gLogInfo.gLogFifo.PutEntryFifo(NewEntryFifo(int32(1), []byte(ret)))
	}
}

func Alert(v ...interface{}) {

	var ret string
	if gLogInfo.gLogLeavel >= ALERT {
		ret = "[" + GetCurTimeString() + "]:" + GetParamString(v)
		gLogInfo.gLogFifo.PutEntryFifo(NewEntryFifo(int32(1), []byte(ret)))
	}
}

func Warn(v ...interface{}) {

	var ret string
	if gLogInfo.gLogLeavel >= WARN {
		ret = "[" + GetCurTimeString() + "]:" + GetParamString(v)
		gLogInfo.gLogFifo.PutEntryFifo(NewEntryFifo(int32(1), []byte(ret)))
	}
}

func Info(v ...interface{}) {

	var ret string
	if gLogInfo.gLogLeavel >= INFO {
		ret = "[" + GetCurTimeString() + "]:" + GetParamString(v)
		gLogInfo.gLogFifo.PutEntryFifo(NewEntryFifo(int32(1), []byte(ret)))
	}
}

func Debug(v ...interface{}) {

	var ret string
	if gLogInfo.gLogLeavel >= DEBUG {
		ret = "[" + GetCurTimeString() + "]:" + GetParamString(v)
		gLogInfo.gLogFifo.PutEntryFifo(NewEntryFifo(int32(1), []byte(ret)))
	}
}
