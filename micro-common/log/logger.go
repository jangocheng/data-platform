package log

import (
	"github.com/pkg/errors"
	"os"
	logutil "platform/common/logging"
)

type Level string
var (
	Trace Level = "TRACE"
	Debug Level = "DEBUG"
	Info Level = "INFO"
	Error Level = "ERROR"
	Fatal Level = "FATAL"
	WARN  Level = "WARN"
)

type Logger struct {
	path   			string
	fileLogger 		*logutil.Logger
	streamLogger 	*logutil.Logger
}

func NewLogger(logPath ...string) (*Logger, error) {
	logger := &Logger{}
	if len(logPath) != 0 {
		fileHandler, err := logutil.NewRotatingFileHandler(logPath[0], 1024*1024*128, 100)
		if err != nil {
			return nil, errors.Wrap(err, "fail to new logger")
		}
		fileLogger := logutil.New(fileHandler, logutil.Ltime|logutil.Llevel)
		logger.path = logPath[0]
		logger.fileLogger = fileLogger
	} else {
		steamHandler, err := logutil.NewStreamHandler(os.Stdout)
		if err != nil {
			return nil, errors.Wrap(err, "fail to new logger")
		}
		streamLogger := logutil.New(steamHandler, logutil.Ltime|logutil.Llevel)
		logger.streamLogger = streamLogger
	}
	return logger, nil
}


func (l *Logger) Close() {
	if l.fileLogger != nil {
		l.fileLogger.Close()
	}
	if l.streamLogger != nil {
		l.streamLogger.Close()
	}
}
