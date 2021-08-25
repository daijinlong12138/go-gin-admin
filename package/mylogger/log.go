package mylogger

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"time"
)

//FileLogger 往文件里面写日志相关代码  文件日志结构体
type FileLogger struct {
	Level    int //目前没用，可作为分不同等级日志打印
	filePath string
	fileName string
	fileObj  *os.File
}

//NewFileLogger 构造方法
func NewFileLogger(fp string, fn string) *FileLogger {

	//给名称添加 .20210824
	fn = fn + "." + time.Now().Format("20060102")
	fl := &FileLogger{
		Level:    15,
		filePath: fp,
		fileName: fn,
	}

	err := fl.initFile() //按照文件路径和文件名将文件打开
	if err != nil {
		panic(err)
	}
	return fl
}

//根据指定的文件路径和文件名打开文件
func (f *FileLogger) initFile() error {
	fullName := path.Join(f.filePath, f.fileName)
	fileObj, err := os.OpenFile(fullName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("open log file failed,err:%v\n", err)
		return err
	}
	//日志文件都已经打开
	f.fileObj = fileObj
	return nil
}

//记录日志的方法
func (f *FileLogger) log(lv int, msg string, arg ...interface{}) {
	// fmt.Println(msg)
	fullMsg := fmt.Sprintf(msg, arg...)
	// fmt.Println(fullMsg)
	now := time.Now()
	fmt.Fprintf(f.fileObj, "[%s] [%s] %s\n", now.Format("2006-01-02 15:04:05"), getLogString(lv), fullMsg)

}

// 获取正在运行的函数名
func runFuncName() string {
	pc := make([]uintptr, 1)
	runtime.Callers(2, pc)
	f := runtime.FuncForPC(pc[0])
	return f.Name()
}

//Debug 方法
func (f *FileLogger) Debug(msg string, arg ...interface{}) {

	msg += "["

	msg += "]"
	f.log(1, msg, arg...)
}

//Info 方法
func (f *FileLogger) Info(msg string, arg ...interface{}) {
	f.log(3, msg, arg...)

}

//Warn 方法
func (f *FileLogger) Warn(msg string, arg ...interface{}) {
	f.log(4, msg, arg...)
}

//Error 方法
func (f *FileLogger) Error(msg string, arg ...interface{}) {
	f.log(5, msg, arg...)
}

//Fatal 方法
func (f *FileLogger) Fatal(msg string, arg ...interface{}) {
	f.log(6, msg, arg...)
}

func getLogString(lv int) string {
	switch lv {
	case 1:
		return "DEBUG"
	case 2:
		return "TRACE"
	case 3:
		return "INFO"
	case 4:
		return "WARING"
	case 5:
		return "ERROR"
	case 6:
		return "FATAL"
	default:
		return "UNKNOW"
	}
}
