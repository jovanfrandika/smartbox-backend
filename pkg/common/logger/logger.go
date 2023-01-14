package logger

import (
	"os"

	"github.com/jovanfrandika/smartbox-backend/pkg/common/utils"
	log "github.com/sirupsen/logrus"
)

type Log struct {
	utils.CallerInfo
}

func Init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}

func Info(message string, callerBackStack int) {
	callerInfo, err := utils.GetCallerInfo(callerBackStack + 2)
	if err != nil {
		log.Info(message)
	} else {
		log.WithFields(log.Fields{
			"func_name": callerInfo.FuncName,
			"file_name": callerInfo.File,
			"line_no":   callerInfo.LineNo,
		}).Info(message)
	}
}

func Error(message string, callerBackStack int) {
	callerInfo, err := utils.GetCallerInfo(callerBackStack + 2)
	if err != nil {
		log.Error(message)
	} else {
		log.WithFields(log.Fields{
			"func_name": callerInfo.FuncName,
			"file_name": callerInfo.File,
			"line_no":   callerInfo.LineNo,
		}).Error(message)
	}
}

func Fatal(message string, callerBackStack int) {
	callerInfo, err := utils.GetCallerInfo(callerBackStack + 2)
	if err != nil {
		log.Fatal(message)
	} else {
		log.WithFields(log.Fields{
			"func_name": callerInfo.FuncName,
			"file_name": callerInfo.File,
			"line_no":   callerInfo.LineNo,
		}).Fatal(message)
	}
}
