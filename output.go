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

var _ panyl.Output = (*Output)(nil)

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

func (o *Output) OnItem(ctx context.Context, item *panyl.Item) (cont bool) {
	outdata := &OutputData{}

	logOptions := []ecapplog.LogOption{
		ecapplog.WithSource(util.DoAnsiEscapeString(item.Source)),
	}

	// timestamp
	if ts, ok := item.Metadata[panyl.MetadataTimestamp]; ok {
		outdata.Time = ts.(time.Time)
	}

	// application
	outdata.Category = ecapplog.CategoryDEFAULT
	var application string
	if application = item.Metadata.StringValue(panyl.MetadataApplication); application != "" && o.applicationAsCategory {
		outdata.Category = application
	}

	// level
	outdata.Priority = ecapplog.Priority_INFORMATION
	if level := item.Metadata.StringValue(panyl.MetadataLevel); level != "" {
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
	if dcategory := item.Metadata.StringValue(panyl.MetadataCategory); dcategory != "" {
		if o.applicationAsCategory && application != "" {
			if o.appendCategoryToApplication {
				outdata.Category = fmt.Sprintf("%s-%s", application, dcategory)
			}
		} else {
			outdata.Category = dcategory
		}
	}

	// original category
	if doriginalcategory := item.Metadata.StringValue(panyl.MetadataOriginalCategory); doriginalcategory != "" {
		logOptions = append(logOptions, ecapplog.WithOriginalCategory(doriginalcategory))
	}

	// message
	if msg := item.Metadata.StringValue(panyl.MetadataMessage); msg != "" {
		outdata.Message = msg
	} else if len(item.Data) > 0 {
		dt, err := json.Marshal(item.Data)
		if err != nil {
			outdata.Message = fmt.Sprintf("Error marshaling data to json: %s", err.Error())
		} else {
			outdata.Message = string(dt)
		}
	} else if item.Line != "" {
		outdata.Message = item.Line
	}

	// extra categories
	if item.Metadata.HasValue(panyl.MetadataExtraCategories) {
		logOptions = append(logOptions, ecapplog.WithExtraCategories(item.Metadata.ListValue(panyl.MetadataExtraCategories)))
	}

	// color
	if color := item.Metadata.StringValue(MetadataColor); color != "" {
		logOptions = append(logOptions, ecapplog.WithColor(color))
	}

	// bgcolor
	if color := item.Metadata.StringValue(MetadataBGColor); color != "" {
		logOptions = append(logOptions, ecapplog.WithBgColor(color))
	}

	// output customization
	if o.customizeOutput != nil {
		if !o.customizeOutput(item, outdata) {
			return
		}
	}

	o.client.Log(outdata.Time, outdata.Priority, outdata.Category, outdata.Message, logOptions...)
	return true
}

func (o *Output) OnFlush(ctx context.Context) {}

func (o *Output) OnClose(ctx context.Context) {
	o.client.Close()
}
