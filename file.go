package logger

import (
	"fmt"
	"os"
	"path"
	"time"
)

type FileLogger struct {
	Level       LogLevel
	filePath    string
	fileName    string
	errFileName string
	maxFileSize int64
	fileObj     *os.File
	errFileObj  *os.File
	writerChan  chan *logMsg
	newFileName string
}

type logMsg struct {
	levelInt  LogLevel
	levelStr  string
	msg       string
	funcName  string
	fileName  string
	line      int
	timestamp string
}

// check file size
// func (f *FileLogger) checkSize(file *os.File) bool {
// 	fileInfo, err := file.Stat()
// 	if err != nil {
// 		fmt.Printf("gefile failed %v", err)
// 		return false
// 	}
// 	return fileInfo.Size() >= f.maxFileSize
// }

// check current time
func (f *FileLogger) checkTime() bool {
	// 1m per
	return time.Now().Format("05") == "00"
	// 1h per
	// return time.Now().Format("0405") == "0000"
	// 1day per
	// return time.Now().Format("150405") == "000000"
}

// Constructor [return filelogger object]
func NewFileLogger(levelStr, fp, fn string, maxSize int64) *FileLogger {
	LogLevel, err := levelConversion(levelStr)
	if err != nil {
		panic(err)
	}
	fl := &FileLogger{
		Level:       LogLevel,
		filePath:    fp,
		fileName:    fn,
		maxFileSize: maxSize,
		errFileName: fmt.Sprintf("%s.err", fn),
		writerChan:  make(chan *logMsg, 50000),
	}
	err = fl.initFile()
	if err != nil {
		panic(err)
	}
	return fl
}

func (f *FileLogger) writeBackground() {
	for {
		if f.checkTime() {
			nowStr := time.Now().Format("20060102150405000")
			pathName := path.Join(f.filePath, f.fileName)
			newlogName := fmt.Sprintf("%s.back%s", pathName, nowStr)
			if f.newFileName != newlogName {
				f.newFileName = newlogName
				f.fileObj.Close()
				os.Rename(pathName, newlogName)
				fileObj, err := os.OpenFile(pathName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
				if err != nil {
					fmt.Print("open new file log failed")
				}
				f.fileObj = fileObj
			}
		}

		select {
		case logTmp := <-f.writerChan:
			logInfo := fmt.Sprintf("[%s] [%s] [%s:%s:%d] %s\n", logTmp.timestamp, logTmp.levelStr, logTmp.funcName, logTmp.fileName, logTmp.line, logTmp.msg)
			fmt.Fprint(f.fileObj, logInfo)
			if logTmp.levelInt >= ERROR {
				fmt.Fprint(f.errFileObj, logInfo)
			}
		default:
			time.Sleep(time.Millisecond * 50)
		}

	}
}

// write log & cut log
func (f *FileLogger) log(level, msg string, a ...interface{}) {
	if f.enable(ERROR) {
		levelInt, err := levelConversion(level)
		if err != nil {
			fmt.Print("conversion error")
		}
		msg = fmt.Sprintf(msg, a...)
		timestamp := time.Now().Format("2006-01-02 15:04:05")
		funcName, fileName, lineNum := getInfo(3)
		logTmp := &logMsg{
			levelInt:  levelInt,
			levelStr:  level,
			msg:       msg,
			funcName:  funcName,
			fileName:  fileName,
			line:      lineNum,
			timestamp: timestamp,
		}
		select {
		case f.writerChan <- logTmp:
		default:
			// 保证业务代码的顺畅执行
		}
		// if f.checkSize(f.fileObj) {
		// 	f.fileObj.Close()
		// 	nowStr := time.Now().Format("20060102150405000")
		// 	pathName := path.Join(f.filePath, f.fileName)
		// 	newlogName := fmt.Sprintf("%s.back%s", pathName, nowStr)
		// 	os.Rename(pathName, newlogName)
		// 	fileObj, err := os.OpenFile(pathName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		// 	if err != nil {
		// 		fmt.Print("open new file log failed")
		// 	}
		// 	f.fileObj = fileObj
		// }

		// if f.checkTime() {
		// 	f.fileObj.Close()
		// 	nowStr := time.Now().Format("20060102150405000")
		// 	pathName := path.Join(f.filePath, f.fileName)
		// 	newlogName := fmt.Sprintf("%s.back%s", pathName, nowStr)
		// 	os.Rename(pathName, newlogName)
		// 	fileObj, err := os.OpenFile(pathName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		// 	if err != nil {
		// 		fmt.Print("open new file log failed")
		// 	}
		// 	f.fileObj = fileObj
		// }

		// fmt.Fprintf(f.fileObj, "[%s] [%s] [%s:%s:%d] %s\n", timestamp, level, funcName, fileName, lineNum, msg)
		// if levelInt >= ERROR {
		// 	fmt.Fprintf(f.errFileObj, "[%s] [%s] [%s:%s:%d] %s\n", timestamp, level, funcName, fileName, lineNum, msg)
		// }
	}
}

// initialize file object
func (f *FileLogger) initFile() error {
	fullFileName := path.Join(f.filePath, f.fileName)
	fileObj, err := os.OpenFile(fullFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("openfile failed err%v", err)
		return err
	}

	fullErrFileName := path.Join(f.filePath, f.errFileName)
	errFileObj, err := os.OpenFile(fullErrFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("openfile failed err%v", err)
		return err
	}
	f.fileObj = fileObj
	f.errFileObj = errFileObj
	// 开启一个后台的groutine，写日志
	for i := 0; i < 4; i++ {
		go f.writeBackground()
	}
	return nil
}

// start logging [判断级别是否写入]
func (f *FileLogger) enable(logLevel LogLevel) bool {
	return logLevel >= f.Level
}

func (f *FileLogger) Debug(msg string, arg ...interface{}) {
	if f.enable(DEBUG) {
		f.log("DEBUG", msg, arg...)
	}
}

func (f *FileLogger) Info(msg string, arg ...interface{}) {
	if f.enable(INFO) {
		f.log("INFO", msg, arg...)
	}
}

func (f *FileLogger) Warn(msg string, arg ...interface{}) {
	if f.enable(WARN) {
		f.log("WARN", msg, arg...)
	}
}

func (f *FileLogger) Error(msg string, arg ...interface{}) {
	if f.enable(ERROR) {
		f.log("ERROR", msg, arg...)
	}
}

func (f *FileLogger) Fatal(msg string, arg ...interface{}) {
	if f.enable(FATAL) {
		f.log("FATAL", msg, arg...)
	}
}

func (f *FileLogger) Close() {
	f.fileObj.Close()
	f.errFileObj.Close()
}
