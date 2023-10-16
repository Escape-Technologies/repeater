package logger

import (
	"fmt"

	proto "github.com/Escape-Technologies/repeater/proto/repeater/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var logSink = make(chan *proto.Log)

func Log(level proto.LogLevel, message string, v ...any) {
	log := proto.Log{
		Level: level,
		Message: fmt.Sprintf(message, v...),
		Timestamp: timestamppb.Now(),
	}
	fmt.Printf("%v %v %v", log.Timestamp, log.Level, log.Message)
	logSink <- &log
}

func Debug(message string, v ...any) { Log(proto.LogLevel_DEBUG, message, v...) }
func Info(message string, v ...any) { Log(proto.LogLevel_INFO, message, v...) }
func Warn(message string, v ...any) { Log(proto.LogLevel_WARN, message, v...) }
func Error(message string, v ...any) { Log(proto.LogLevel_ERROR, message, v...) }
