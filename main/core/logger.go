package core

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"time"
)

const ()

type Logger struct {
	level Level

	writeToFile bool
	writeToTerm bool

	fileNamePrefix string
	filePath       string

	//prefix20041002-1.log
	//prefix20041002-2.log
	//TODO: it should be protected by a mutex
	maxSize int64

	//prefix20041002-1.log
	//prefix20041003-1.log
	logInterval time.Duration

	file    *os.File
	errFile *os.File
}

// -----------> options and constructor

type LoggerOption func(*Logger)

func WithLevel(level Level) LoggerOption {
	return func(l *Logger) {
		l.level = level
	}
}

func WithFilePrefixName(fileName string) LoggerOption {
	return func(l *Logger) {
		l.fileNamePrefix = fileName
	}
}

func WithFilePath(filePath string) LoggerOption {
	return func(l *Logger) {
		l.filePath = filePath
	}
}

func WithMaxSize(maxSize int64) LoggerOption {
	return func(l *Logger) {
		l.maxSize = maxSize
	}
}

func WithLogInterval(logInterval time.Duration) LoggerOption {
	return func(l *Logger) {
		l.logInterval = logInterval
	}
}

func WithWriteToTerm(writeToTerm bool) LoggerOption {
	return func(l *Logger) {
		l.writeToTerm = writeToTerm
	}

}

/*
if not providing the fileName and filePath, the logger won't generate log info to the disk.

	prefix!="", filePath=="": that means generate log file in the current directory.

	prefix=="", filePath!="": that means generate log file in the specified directory without the custom prefix.

	prefix!="", filePath!="": that means generate log file in the specified directory with the custom prefix.

	prefix=="", filePath=="": that means don't generate log file.
*/
func NewFileLoggerWithOption(options ...LoggerOption) *Logger {

	loggerPtr := new(Logger)

	for _, option := range options {
		option(loggerPtr)
	}

	//level not set by options
	if loggerPtr.level == 0 {
		loggerPtr.level = DEBUG
	}

	if loggerPtr.fileNamePrefix == "" && loggerPtr.filePath == "" {
		loggerPtr.writeToFile = false
	} else {
		loggerPtr.writeToFile = true
	}

	if !loggerPtr.writeToTerm {
		loggerPtr.writeToTerm = true
	}

	if loggerPtr.maxSize == 0 {
		// loggerPtr.maxSize = 1024 * 1024 * 10 // 10MB
		loggerPtr.maxSize = 1024 // 1kB
	}
	if loggerPtr.logInterval == 0 {
		loggerPtr.logInterval = DAY
	}

	if loggerPtr.writeToFile {

		var timeSuffix string
		switch loggerPtr.logInterval {

		case DAY:
			timeSuffix = time.Now().Format("20060102")

		case MONTH:
			timeSuffix = time.Now().Format("200601")

		}

		fullPathPrefix := loggerPtr.fileNamePrefix + timeSuffix

		fullPathNameLog := filepath.Join(loggerPtr.filePath, fullPathPrefix+"-1.log")

		fullPathNameErr := filepath.Join(loggerPtr.filePath, fullPathPrefix+"-1.err")

		logFile := openFile(fullPathNameLog)
		errFile := openFile(fullPathNameErr)

		loggerPtr.file = logFile
		loggerPtr.errFile = errFile
	}

	if !loggerPtr.writeToFile && !loggerPtr.writeToTerm {
		panic("no log output specified, why use the logger?")
	}

	return loggerPtr
}

// <---------- options and constructor ends

func (fl *Logger) Debug(format string, args ...any) {
	fl.writeLog(DEBUG, format, args...)
}

func (fl *Logger) Info(format string, args ...any) {
	fl.writeLog(INFO, format, args...)
}

func (fl *Logger) Warn(format string, args ...any) {
	fl.writeLog(WARN, format, args...)
}

func (fl *Logger) Error(format string, args ...any) {

	fl.writeLog(ERROR, format, args...)
}
func (fl *Logger) Fatal(format string, args ...any) {
	fl.writeLog(FATAL, format, args...)
}

func (fl *Logger) writeLog(level Level, format string, args ...any) {

	if level < fl.level {
		return
	}

	//msg is the payload
	msg := fmt.Sprintf(format, args...)

	//logMsg is the completed log msg
	logMsg := fl.prepareLogInfo(msg, level)

	if fl.writeToTerm {
		fmt.Println(logMsg)
	}

	if fl.writeToFile {
		fl.doWriteToFile(logMsg, level)
	}
}

/* logFormat: [time][fileName:lineNumber][methodName][logLevel] - logMsg */
func (fl *Logger) prepareLogInfo(msg string, level Level) string {
	nowStr := time.Now().Format("20060102-15:04:05")

	fileName, line, funcName := getCallerInfo(4)

	logMsg := fmt.Sprintf("[%s][%s:%d][%s][%s] - %s", nowStr, fileName, line, funcName, getLevelString(level), msg)

	return logMsg
}

func (fl *Logger) doWriteToFile(logMsg string, level Level) {

	var timeSuffix string

	switch fl.logInterval {
	case DAY:
		timeSuffix = time.Now().Format("20060102")
	case MONTH:
		timeSuffix = time.Now().Format("200601")
	}

	logFileStat, _ := fl.file.Stat()
	errFileStat, _ := fl.errFile.Stat()

	logMsgSize := int64(len(logMsg))

	if logFileStat.Size()+logMsgSize > fl.maxSize {

		newFileCount := getFileCount(fl.file.Name()) + 1

		fullPathNameLog :=
			filepath.Join(fl.filePath, fl.fileNamePrefix+timeSuffix+"-"+strconv.Itoa(newFileCount)+".log")

		fl.file = openFile(fullPathNameLog)
	}
	fmt.Fprintln(fl.file, logMsg)

	if level == ERROR || level == FATAL {
		if errFileStat.Size()+logMsgSize > fl.maxSize {
			newFileCount := getFileCount(fl.errFile.Name()) + 1
			fullPathNameErr :=
				filepath.Join(fl.filePath, fl.fileNamePrefix+timeSuffix+"-"+strconv.Itoa(newFileCount)+".err")

			fl.errFile = openFile(fullPathNameErr)

		}
		fmt.Fprintln(fl.errFile, logMsg)
	}

}

/*
"1-PrefxiTest20201231-1.log", -1.log, the 1 here

"-1PrefxiTest20201231-2.log", -2.log, the 2 here

"PrefxiTest20201231-3.log", -3.log, the 3 here
*/
func getFileCount(filePath string) int {

	re := regexp.MustCompile(`-(\d+).log`)

	matches := re.FindStringSubmatch(filePath)
	if len(matches) > 0 {

		num, err := strconv.Atoi(matches[1])

		if err != nil {
			panic(fmt.Errorf("failed to parse file count: %v", filePath))
		}
		return num
	}
	panic(fmt.Errorf("failed to parse file count: %v", filePath))
}

/* helpers */

/* constructor */

/*
@param level: could be null, default is DEBUG
*/
func NewFileLogger(fileName string, filePath string, level ...Level) *Logger {

	// Correctly use path.Join to combine fileName and filePath
	//full path
	fullPathName := path.Join(filePath, fileName)

	fullPathNameLog := fullPathName + ".log"
	fullPathNameErr := fullPathName + ".err"

	logFile := openFile(fullPathNameLog)
	errFile := openFile(fullPathNameErr)

	_level := DEBUG

	if len(level) > 0 {
		_level = level[0]
	}

	return &Logger{
		level:          _level,
		fileNamePrefix: fileName,
		filePath:       filePath,
		file:           logFile,
		errFile:        errFile,
	}
}

func openFile(fullPathName string) *os.File {

	fileObj, err := os.OpenFile(
		fullPathName,
		os.O_CREATE|os.O_APPEND|os.O_WRONLY,
		0664)

	if err != nil {
		panic(fmt.Errorf("failed to create log file: %v", fullPathName))
	}

	return fileObj
}
