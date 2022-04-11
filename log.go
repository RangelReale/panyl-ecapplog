package panylecapplog

import (
	"encoding/json"
	"fmt"
	"github.com/RangelReale/ecapplog-go"
	"github.com/RangelReale/panyl"
	"time"
)

var _ panyl.Log = (*Log)(nil)

type Log struct {
	client          *ecapplog.Client
	sourceCategory  string
	sourcePriority  ecapplog.Priority
	processCategory string
	processPriority ecapplog.Priority
}

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

func (l Log) LogSourceLine(n int, line, rawLine string) {
	l.client.Log(time.Now(), l.sourcePriority, l.sourceCategory,
		fmt.Sprintf("@@@ SOURCE LINE [%d]: '%s'", n, line), ecapplog.WithSource(rawLine))
}

func (l Log) LogProcess(p *panyl.Process) {
	var lineno string
	if p.LineCount > 1 {
		lineno = fmt.Sprintf("[%d-%d]", p.LineNo, p.LineNo+p.LineCount-1)
	} else {
		lineno = fmt.Sprintf("[%d]", p.LineNo)
	}

	var message string
	if msg := p.Metadata.StringValue(panyl.Metadata_Message); msg != "" {
		message = msg
	} else if len(p.Data) > 0 {
		dt, err := json.Marshal(p.Data)
		if err != nil {
			message = fmt.Sprintf("Error marshaling data to json: %s", err.Error())
		} else {
			message = string(dt)
		}
	} else if p.Line != "" {
		message = p.Line
	}

	l.client.Log(time.Now(), l.processPriority, l.processCategory,
		fmt.Sprintf("*** PROCESS LINE %s: '%s'", lineno, message), ecapplog.WithSource(p.Source))
}
