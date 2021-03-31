package logger

import (
	"fmt"
	"time"
)

// return logger obj struct
type Logger struct {
	Level    LogLevel
	LevelStr string
}

// 构造函数
func Initlalize(levelStr string) Logger {
	levelNum, err := levelConversion(levelStr)
	if err != nil {
		panic(err)
	}
	return Logger{
		Level:    levelNum,
		LevelStr: levelStr,
	}
}

func (l Logger) log(lv LogLevel, msg string, a ...interface{}) {
	if l.enable(ERROR) {
		msg = fmt.Sprintf(msg, a...)
		timestamp := time.Now().Format("2006-01-02 15:04:05")
		funcName, fileName, lineNu := getInfo(3)
		info := fmt.Sprintf("[%s] [%d] [DEBUG] [%s:%s:%d]%s\n", timestamp, lv, funcName, fileName, lineNu, msg)
		fmt.Print(info)
	}
}

func (l Logger) enable(logLevel LogLevel) bool {
	return logLevel >= l.Level
}

func (l Logger) Debug(msg string, arg ...interface{}) {
	l.log(DEBUG, msg, arg...)

}

func (l Logger) INFO(msg string) {
	fmt.Println(msg)
}

func (l Logger) WARN(msg string) {
	fmt.Println(msg)
}

func (l Logger) ERROR(msg string) {
	fmt.Println(msg)
}

func (l Logger) FATAL(msg string) {
	fmt.Println(msg)
}
