package utils

import (
	"errors"
	"path"
	"runtime"
)

type CallerInfo struct {
	FuncName string
	File     string
	LineNo   int
}

func GetCallerInfo(skip int) (CallerInfo, error) {
	pc, file, lineNo, ok := runtime.Caller(skip)
	if !ok {
		return CallerInfo{}, errors.New("runtime.Caller failed")
	}
	funcName := runtime.FuncForPC(pc).Name()
	fileName := path.Base(file)
	return CallerInfo{
		FuncName: funcName,
		File:     fileName,
		LineNo:   lineNo,
	}, nil
}
