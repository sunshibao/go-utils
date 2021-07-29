// Author: Sunshibao <664588619@qq.com>
package configsdk

import (
	"fmt"
)

var defaultLogger Logger = emptyLogger(0)

func GetLogger() Logger {
	return defaultLogger
}

func SetLogger(l Logger) {
	if l == nil {
		defaultLogger = emptyLogger(0)
	} else {
		defaultLogger = l
	}
}

type Logger interface {
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
}

type emptyLogger int

func (emptyLogger) Debug(args ...interface{}) { fmt.Println(fmt.Sprint(args...)) }
func (emptyLogger) Info(args ...interface{})  { fmt.Println(fmt.Sprint(args...)) }
func (emptyLogger) Warn(args ...interface{})  { fmt.Println(fmt.Sprint(args...)) }
func (emptyLogger) Error(args ...interface{}) { fmt.Println(fmt.Sprint(args...)) }

func (emptyLogger) Debugf(format string, args ...interface{}) {
	fmt.Println(fmt.Sprintf(format, args...))
}

func (emptyLogger) Infof(format string, args ...interface{}) {
	fmt.Println(fmt.Sprintf(format, args...))
}

func (emptyLogger) Warnf(format string, args ...interface{}) {
	fmt.Println(fmt.Sprintf(format, args...))
}

func (emptyLogger) Errorf(format string, args ...interface{}) {
	fmt.Println(fmt.Sprintf(format, args...))
}
