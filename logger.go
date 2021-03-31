package logger

import (
	"errors"
	"fmt"
	"path"
	"runtime"
	"strings"
)

const (
	INVAILD LogLevel = iota
	DEBUG
	INFO
	WARN
	ERROR
	FATAL
)

// create your own type
type LogLevel uint16

// level conversion
func levelConversion(levelStr string) (LogLevel, error) {
	switch strings.ToUpper(levelStr) {
	case "DEBUG":
		return DEBUG, nil
	case "INFO":
		return INFO, nil
	case "WARN":
		return WARN, nil
	case "ERROR":
		return ERROR, nil
	case "FATAL":
		return FATAL, nil
	default:
		err := errors.New("无效的日志级别")
		return INVAILD, err
	}
}

// get the calling function、filename、line number
func getInfo(skip int) (funcName, fileName string, lineNu int) {
	pc, file, lineNu, ok := runtime.Caller(skip)
	if !ok {
		fmt.Println("函数调用检测失败")
		return
	}
	funcName = runtime.FuncForPC(pc).Name()
	fileName = path.Base(file)
	return
}
