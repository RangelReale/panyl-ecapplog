package panylecapplog

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/RangelReale/ecapplog-go"
	"github.com/RangelReale/panyl/v2"
	"github.com/RangelReale/panyl/v2/util"
)

var _ panyl.ProcessResult = (*Output)(nil)

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
	Time             time.Time
	Priority         ecapplog.Priority
	Category         string
	OriginalCategory string
	Message          string
	ExtraCategories  []string
}

func (o *Output) OnResult(ctx context.Context, p *panyl.Item) (cont bool) {
	outdata := &OutputData{}

	// timestamp
	if ts, ok := p.Metadata[panyl.MetadataTimestamp]; ok {
		outdata.Time = ts.(time.Time)
	}

	// application
	outdata.Category = ecapplog.CategoryDEFAULT
	var application string
	if application = p.Metadata.StringValue(panyl.MetadataApplication); application != "" && o.applicationAsCategory {
		outdata.Category = application
	}

	// level
	outdata.Priority = ecapplog.Priority_INFORMATION
	if level := p.Metadata.StringValue(panyl.MetadataLevel); level != "" {
		switch level {
		case panyl.MetadataLevelTRACE:
			outdata.Priority = ecapplog.Priority_TRACE
		case panyl.MetadataLevelDEBUG:
			outdata.Priority = ecapplog.Priority_DEBUG
		case panyl.MetadataLevelINFO:
			outdata.Priority = ecapplog.Priority_INFORMATION
		case panyl.MetadataLevelWARNING:
			outdata.Priority = ecapplog.Priority_WARNING
		case panyl.MetadataLevelERROR:
			outdata.Priority = ecapplog.Priority_ERROR
		}
	}

	// category
	if dcategory := p.Metadata.StringValue(panyl.MetadataCategory); dcategory != "" {
		if o.applicationAsCategory && application != "" {
			if o.appendCategoryToApplication {
				outdata.Category = fmt.Sprintf("%s-%s", application, dcategory)
			}
		} else {
			outdata.Category = dcategory
		}
	}

	// original category
	if doriginalcategory := p.Metadata.StringValue(panyl.MetadataOriginalCategory); doriginalcategory != "" {
		outdata.OriginalCategory = doriginalcategory
	}

	// message
	if msg := p.Metadata.StringValue(panyl.MetadataMessage); msg != "" {
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
		ecapplog.WithSource(util.DoAnsiEscapeString(p.Source)),
		ecapplog.WithOriginalCategory(outdata.OriginalCategory),
		ecapplog.WithExtraCategories(outdata.ExtraCategories))
	return true
}

func (o *Output) OnFlush(ctx context.Context) {}

func (o *Output) OnClose(ctx context.Context) {
	o.client.Close()
}
