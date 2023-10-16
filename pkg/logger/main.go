package logger

import (
	"fmt"

	"github.com/Escape-Technologies/repeater/pkg/fifo"
	proto "github.com/Escape-Technologies/repeater/proto/repeater/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var queue = fifo.NewQueue[proto.Log]()

func Log(level proto.LogLevel, message string, v ...any) {
	log := proto.Log{
		Level:     level,
		Message:   fmt.Sprintf(message, v...),
		Timestamp: timestamppb.Now(),
	}
	fmt.Printf("%v %v %v\n", log.Timestamp.AsTime(), log.Level, log.Message)
	queue.Add(&log)
}

func Debug(message string, v ...any) { Log(proto.LogLevel_DEBUG, message, v...) }
func Info(message string, v ...any)  { Log(proto.LogLevel_INFO, message, v...) }
func Warn(message string, v ...any)  { Log(proto.LogLevel_WARN, message, v...) }
func Error(message string, v ...any) { Log(proto.LogLevel_ERROR, message, v...) }
