package logs

import (
	"fmt"

	"github.com/xuperchain/xasset-sdk-go/utils"
)

// 底层日志库约束接口
type LogDriver interface {
	Error(msg string, ctx ...interface{})
	Warn(msg string, ctx ...interface{})
	Info(msg string, ctx ...interface{})
	Trace(msg string, ctx ...interface{})
	Debug(msg string, ctx ...interface{})
}

type Logger struct {
	logDriver LogDriver
}

func NewLogger(logDriver LogDriver) *Logger {
	return &Logger{
		logDriver: logDriver,
	}
}

func (t *Logger) Error(msg string, ctx ...interface{}) {
	if t.logDriver == nil {
		return
	}

	t.logDriver.Error(t.fmtMsg(msg, ctx...))
}

func (t *Logger) Warn(msg string, ctx ...interface{}) {
	if t.logDriver == nil {
		return
	}

	t.logDriver.Warn(t.fmtMsg(msg, ctx...))
}

func (t *Logger) Info(msg string, ctx ...interface{}) {
	if t.logDriver == nil {
		return
	}

	t.logDriver.Info(t.fmtMsg(msg, ctx...))
}

func (t *Logger) Trace(msg string, ctx ...interface{}) {
	if t.logDriver == nil {
		return
	}

	t.logDriver.Trace(t.fmtMsg(msg, ctx...))
}

func (t *Logger) Debug(msg string, ctx ...interface{}) {
	if t.logDriver == nil {
		return
	}

	t.logDriver.Debug(t.fmtMsg(msg, ctx...))
}

func (t *Logger) genBaseMsg() string {
	call, _ := utils.GetFuncCall(4)
	return fmt.Sprintf("[sdk_call:%s]", call)
}

func (t *Logger) fmtMsg(msg string, ctx ...interface{}) string {
	if msg == "" {
		return t.genBaseMsg()
	}

	return fmt.Sprintf(msg+" "+t.genBaseMsg(), ctx...)
}
