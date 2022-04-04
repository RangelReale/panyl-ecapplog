package panylecapplog

import (
	"encoding/json"
	"fmt"
	"github.com/RangelReale/ecapplog-go"
	"github.com/RangelReale/panyl"
	"github.com/RangelReale/panyl/util"
	"time"
)

type Output struct {
	client                      *ecapplog.Client
	customizeOutput             CustomizeOutputFunc
	applicationAsCategory       bool
	appendCategoryToApplication bool
}

func NewOutput(client *ecapplog.Client, options ...OutputOption) *Output {
	ret := &Output{
		client: client,
	}
	for _, opt := range options {
		opt(ret)
	}
	return ret
}

type OutputData struct {
	Time     time.Time
	Priority ecapplog.Priority
	Category string
	Message  string
}

func (o *Output) OnResult(p *panyl.Process) (cont bool) {
	outdata := &OutputData{}

	// timestamp
	if ts, ok := p.Metadata[panyl.Metadata_Timestamp]; ok {
		outdata.Time = ts.(time.Time)
	}

	// application
	outdata.Category = "ALL"
	var application string
	if application = p.Metadata.StringValue(panyl.Metadata_Application); application != "" && o.applicationAsCategory {
		outdata.Category = application
	}

	// level
	outdata.Priority = ecapplog.Priority_INFORMATION
	if level := p.Metadata.StringValue(panyl.Metadata_Level); level != "" {
		switch level {
		case panyl.MetadataLevel_TRACE:
			outdata.Priority = ecapplog.Priority_TRACE
		case panyl.MetadataLevel_DEBUG:
			outdata.Priority = ecapplog.Priority_DEBUG
		case panyl.MetadataLevel_INFO:
			outdata.Priority = ecapplog.Priority_INFORMATION
		case panyl.MetadataLevel_WARNING:
			outdata.Priority = ecapplog.Priority_WARNING
		case panyl.MetadataLevel_ERROR:
			outdata.Priority = ecapplog.Priority_ERROR
		case panyl.MetadataLevel_CRITICAL:
			outdata.Priority = ecapplog.Priority_CRITICAL
		case panyl.MetadataLevel_FATAL:
			outdata.Priority = ecapplog.Priority_FATAL
		}
	}

	// category
	if dcategory := p.Metadata.StringValue(panyl.Metadata_Category); dcategory != "" {
		if o.applicationAsCategory && o.appendCategoryToApplication && application != "" {
			outdata.Category = fmt.Sprintf("%s-%s", application, dcategory)
		} else {
			outdata.Category = dcategory
		}
	}

	// message
	if msg := p.Metadata.StringValue(panyl.Metadata_Message); msg != "" {
		outdata.Message = msg
	} else if len(p.Data) > 0 {
		dt, err := json.Marshal(p.Data)
		if err != nil {
			outdata.Message = fmt.Sprintf("Error marshaling data to json: %s", err.Error())
		} else {
			outdata.Message = string(dt)
		}
	} else if p.Line != "" {
		outdata.Message = p.Line
	}

	// output customization
	if o.customizeOutput != nil {
		if !o.customizeOutput(p, outdata) {
			return
		}
	}

	o.client.Log(outdata.Time, outdata.Priority, outdata.Category, outdata.Message,
		ecapplog.WithSource(util.DoAnsiEscapeString(p.Source)))
	return true
}
