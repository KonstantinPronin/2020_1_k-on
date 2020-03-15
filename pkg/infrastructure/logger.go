package infrastructure

import (
	"encoding/json"
	"go.uber.org/zap"
)

var logConf = []byte(`{
		"level": "debug",
		"encoding": "console",
		"outputPaths": ["stdout", "/tmp/logs"],
		"errorOutputPaths": ["stderr"],
		"encoderConfig": {
		  "messageKey": "message",
          "callerKey": "caller",
		  "callerEncoder": "short",
		  "levelKey": "level",
		  "levelEncoder": "capital",
		  "timeKey": "time",
 		  "timeEncoder": "ISO8601"
		}
	  }`)

func InitLog() (*zap.Logger, error) {
	var conf zap.Config
	if err := json.Unmarshal(logConf, &conf); err != nil {
		return nil, err
	}

	return conf.Build()
}
