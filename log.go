package panylecapplog

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/RangelReale/ecapplog-go"
	"github.com/RangelReale/panyl/v2"
)

type Log struct {
	client          *ecapplog.Client
	sourceCategory  string
	sourcePriority  ecapplog.Priority
	processCategory string
	processPriority ecapplog.Priority
}

var _ panyl.DebugLog = (*Log)(nil)

func NewLog(client *ecapplog.Client, options ...LogOption) *Log {
	ret := &Log{
		client:          client,
		sourceCategory:  "panyl-log",
		sourcePriority:  ecapplog.Priority_TRACE,
		processCategory: "panyl-log",
		processPriority: ecapplog.Priority_INFORMATION,
	}
	for _, opt := range options {
		opt(ret)
	}
	return ret
}

func (l Log) LogSourceLine(ctx context.Context, n int, line, rawLine string) {
	l.client.Log(time.Now(), l.sourcePriority, l.sourceCategory,
		fmt.Sprintf("@@@ SOURCE LINE [%d]: '%s'", n, line), ecapplog.WithSource(rawLine))
}

func (l Log) LogItem(ctx context.Context, item *panyl.Item) {
	var lineno string
	if item.LineCount > 1 {
		lineno = fmt.Sprintf("[%d-%d]", item.LineNo, item.LineNo+item.LineCount-1)
	} else {
		lineno = fmt.Sprintf("[%d]", item.LineNo)
	}

	var message string
	if msg := item.Metadata.StringValue(panyl.MetadataMessage); msg != "" {
		message = msg
	} else if len(item.Data) > 0 {
		dt, err := json.Marshal(item.Data)
		if err != nil {
			message = fmt.Sprintf("Error marshaling data to json: %s", err.Error())
		} else {
			message = string(dt)
		}
	} else if item.Line != "" {
		message = item.Line
	}

	l.client.Log(time.Now(), l.processPriority, l.processCategory,
		fmt.Sprintf("*** PROCESS LINE %s: '%s'", lineno, message), ecapplog.WithSource(item.Source))
}
