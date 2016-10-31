package main

import (
	logger "github.com/alecthomas/log4go"
)

func main () {
	logMech := make(logger.Logger)
	logMech.AddFilter("stdout",logger.DEBUG,logger.NewConsoleLogWriter())

	fileLog := logger.NewFileLogWriter("log_manager.log",false)
	fileLog.SetFormat("[%D %T][%L](%S)%M")
	fileLog.SetRotate(true)
	fileLog.SetRotateSize(256)
	fileLog.SetRotateLines(20)
	logMech.AddFilter("file",logger.FINE,fileLog)

	logMech.Trace("Recieved Message: %s","All is Well")
	logMech.Info("Message Recieved: ","debug!")
	logMech.Error("Oh No!","Something Broke")

}
