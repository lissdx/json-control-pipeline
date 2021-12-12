package process

import "github.com/lissdx/json-control-pipeline/pkg/pipeline"

type Conf string

type Processor interface {
	InitProcess(Conf)
	InStream() pipeline.WriteOnlyStream
	Run()
	Stop()
}
