package core

import (
	"path"
	"runtime"
)

/**
这里的skip参数表示要跳过的调用栈层数，
*/
func getCallerInfo(skip int) (string, int, string) {
	pc, fileFullPath, line, ok := runtime.Caller(skip)

	if !ok {
		return "unknown", -1, "unknown"
	}

	fileName := path.Base(fileFullPath)

	funcName := path.Base(runtime.FuncForPC(pc).Name())

	return fileName, line, funcName

}
