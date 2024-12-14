package panylecapplog

import (
	"github.com/RangelReale/ecapplog-go"
	"github.com/RangelReale/panyl/v2"
)

type LogOption func(*Log)

type OutputOption func(*Output)

func WithSourceCategory(sourceCategory string) LogOption {
	return func(log *Log) {
		log.sourceCategory = sourceCategory
	}
}

func WithSourcePriority(sourcePriority ecapplog.Priority) LogOption {
	return func(log *Log) {
		log.sourcePriority = sourcePriority
	}
}

func WithProcessCategory(processCategory string) LogOption {
	return func(log *Log) {
		log.processCategory = processCategory
	}
}

func WithProcessPriority(processPriority ecapplog.Priority) LogOption {
	return func(log *Log) {
		log.processPriority = processPriority
	}
}

// CustomizeOutputFunc allows customization of data being passed to ECAppLog, by modifying
// the outdata struct.
// Return false if you don't want this process to be output.
type CustomizeOutputFunc func(p *panyl.Process, outdata *OutputData) (doOutput bool)

func WithCustomizeOutput(f CustomizeOutputFunc) OutputOption {
	return func(output *Output) {
		output.customizeOutput = f
	}
}

func WithApplicationAsCategory(applicationAsCategory bool) OutputOption {
	return func(output *Output) {
		output.applicationAsCategory = applicationAsCategory
	}
}

func WithAppendCategoryToApplication(appendCategoryToApplication bool) OutputOption {
	return func(output *Output) {
		output.appendCategoryToApplication = appendCategoryToApplication
	}
}
