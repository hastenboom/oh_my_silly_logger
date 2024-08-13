package core_test

import (
	"oh_my_logger/main/core"
	"testing"
)

func TestWriteFile(t *testing.T) {
	//FIXME:
	// logger := core.NewFileLogger("./log", "xxx.log")

	logger := core.NewFileLogger("./", "xxx.log")

	sb := "hello world"

	logger.Debug(sb)
	logger.Info(sb)
	logger.Fatal(sb)

}
