// 日志
package log

import (
	"io/ioutil"
	"log"
	"os"
	"sync"
)

const (
	InfoLevel = iota
	ErrorLevel
	Disabled
)

var (
	errorLog = log.New(os.Stdout, "\033[31m[error]\033[0m", log.LstdFlags|log.Lshortfile)
	infoLog  = log.New(os.Stdout, "\033[34m[info ]\033[0m", log.LstdFlags|log.Lshortfile)
	loggers  = []*log.Logger{errorLog, infoLog}
	mu       sync.Mutex
)

// log methods
var (
	Error  = errorLog.Println
	Errorf = errorLog.Printf
	Info   = infoLog.Println
	Infof  = infoLog.Printf
)

func SetLevel(level int) {
	mu.Lock()
	defer mu.Unlock()

	for _, logger := range loggers {
		logger.SetOutput(os.Stdout)
	}

	if ErrorLevel < level {
		errorLog.SetOutput((ioutil.Discard))
	}
	if InfoLevel < level { // level > InfoLevel(0), infoLog 的输出被重定向到ioutil.Discard，即不输出。
		infoLog.SetOutput(ioutil.Discard)
	}
}
