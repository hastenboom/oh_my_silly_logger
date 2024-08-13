package core_test

import (
	"fmt"
	"oh_my_logger/main/core"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"testing"
	"time"
)

func TestPath(t *testing.T) {
	/* filePath := "D:\\dev\\go_learning\\syntax\\assignment\\03_log_sys\\oh_my_logger\\test\\core\\log\\" */
	filePath := ""
	fileNamePrefix := "PrefixTest" // 保持为空字符串

	timeSuffix := time.Now().Format("20060102")
	//use the filepath not the filePath or path
	fullPathName := filepath.Join(filePath, fileNamePrefix+timeSuffix+"-1.log")

	// 打开文件
	fileObj, err := os.OpenFile(
		fullPathName,
		os.O_CREATE|os.O_APPEND|os.O_WRONLY,
		0664)

	if err != nil {
		panic(fmt.Errorf("failed to create log file: %v", fullPathName))
	}

	fileObj.Stat()

	defer fileObj.Close() // 确保文件关闭

	// 写入文件
	fmt.Fprintln(fileObj, "hello world")
}

func TestRegex(t *testing.T) {
	var strs []string = []string{
		"1-PrefxiTest20201231-1.log",
		"-1PrefxiTest20201231-2.log",
		"PrefxiTest20201231-3.log",
	}
	re := regexp.MustCompile(`-(\d+).log`)

	for _, str := range strs {
		matches := re.FindStringSubmatch(str)
		if len(matches) > 0 {
			num, err := strconv.Atoi(matches[1])
			if err != nil {
				t.Error(err)
			}
			fmt.Printf("num: %v\n", num)
		}

	}

}

func TestLogger(t *testing.T) {
	filePath := "D:\\dev\\go_learning\\syntax\\assignment\\03_log_sys\\oh_my_logger\\test\\core\\log\\"
	// filePath := ""
	fileNamePrefix := "PrefixTest"

	loggerPtr := core.NewFileLoggerWithOption(
		core.WithFilePrefixName(fileNamePrefix),
		core.WithFilePath(filePath),
	)

	//!WARN: for debug purpose
	for i := 0; i < 1000000; i++ {
		loggerPtr.Debug("debug message")
	}

}
