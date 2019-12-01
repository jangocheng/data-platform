package log

import (
	"encoding/json"
	"github.com/pkg/errors"
	logutil "platform/common/logging"
)

func (l *Logger) Log(msg interface{}, level ...Level) error {
	errMsg := "fail to log msg"
	logLevel := Info
	if len(level) != 0 {
		logLevel = level[0]
	}
	var msgBytes []byte
	var err error
	if msgBytes, ok := msg.([]byte); ok {
		msg = string(msgBytes)
	}
	if msgStr, ok := msg.(string); !ok {
		msgBytes, err = json.Marshal(msg)
		if err != nil {
			return errors.Wrap(err, errMsg)
		}
	} else {
		msgBytes = []byte(msgStr)
	}
	var msgMap map[string]interface{}
	msgMap = make(map[string]interface{})
	err = json.Unmarshal(msgBytes, &msgMap)
	if err == nil {
		for k, v := range msgMap {
			var vLen int
			switch v.(type) {
			case string:
				vStr, _ := v.(string)
				vLen = len([]byte(vStr))
			case []byte:
				vBytes, _ := v.([]byte)
				vLen = len(vBytes)
			default:
				_vBytes, err := json.Marshal(v)
				if err == nil {
					var childMap map[string]interface{}
					childMap = make(map[string]interface{})
					err = json.Unmarshal(_vBytes, &childMap)
					if err != nil {
						_vLen := len(_vBytes) - 2
						if _vLen >= 1*1024*1024 {
							msgMap[k] = ""
						}
					} else {
						for _, _v := range childMap {
							_vBytes, err := json.Marshal(_v)
							if err != nil {
								return errors.Wrap(err, errMsg)
							}
							_vLen := len(_vBytes) - 2
							if _vLen >= 1*1024*1024 {
								childMap[k] = ""
							}
						}
						msgMap[k] = childMap
					}
				}
				vBytes, err := json.Marshal(v)
				if err != nil {
					return errors.Wrap(err, errMsg)
				}
				vLen = len(vBytes) - 2
			}

			if vLen >= 1*1024*1024 {
				msgMap[k] = ""
			}
		}
		msgBytes, err = json.Marshal(&msgMap)
		if err != nil {
			return errors.Wrap(err, errMsg)
		}
	}
	msgStr := string(msgBytes)
	if l.fileLogger != nil {
		doLog(l.fileLogger, logLevel, msgStr)
	} else {
		doLog(l.streamLogger, logLevel, msgStr)
	}
	return nil
}

func doLog(logger *logutil.Logger, level Level, msg string) {
	switch level {
	case Info:
		logger.Info(msg)
	case Debug:
		logger.Debug(msg)
	case Trace:
		logger.Trace(msg)
	case Error:
		logger.Error(msg)
	case Fatal:
		logger.Fatal(msg)
	case WARN:
		logger.Fatal(msg)

	}
}
