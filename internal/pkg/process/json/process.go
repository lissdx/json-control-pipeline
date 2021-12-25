package json

import (
	"encoding/json"
	"fmt"
	"github.com/lissdx/json-control-pipeline/internal/pkg/process"
	"github.com/lissdx/json-control-pipeline/pkg/pipeline"
	"log"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

const sleepTimeMs = 100

// Process main process obj
type Process struct {
	doneChannel           chan interface{}
	inStream              pipeline.WriteOnlyStream
	jsonTransformPipeline pipeline.Pipeline
	activeStatus          atomic.Value
	processWg             *sync.WaitGroup
}

// NewJsonProcess Process constructor
func NewJsonProcess() process.Processor {
	return &Process{}
}

// InitProcess get conf. file and create main process data-pipeline
func (jp *Process) InitProcess(conf process.Conf) {

	// Create process step parser
	processSteps := NewProcessSteps()
	// Init ProcessSteps
	processSteps.Init(string(conf))

	// Add some Start stage to process
	// Not mandatory, just for visual control
	jp.jsonTransformPipeline.AddStage(jp.stubStart(), jp.errorHandler("stubStart"))

	// Fetch every step from ProcessSteps obj.
	// and add it into our process
	// use factory pattern
	for processSteps.HasNext() {
		pNameConf := processSteps.Next()
		// Add process stage
		// The stage is created by ProcessName and conf. Factory pattern
		jp.jsonTransformPipeline.AddStage(jp.ProcessFactory(pNameConf.ProcessName, pNameConf.ProcessConf), jp.errorHandler(pNameConf.ProcessName))
	}
	//jp.jsonTransformPipeline.AddStage(jp.addField(`{"fieldName":"firstName","fieldValue":"George"}`), jp.errorHandler("addField"))
	//jp.jsonTransformPipeline.AddStage(jp.removeField(`{"fieldName":"age"}`), jp.errorHandler("removeField"))
	//jp.jsonTransformPipeline.AddStage(jp.countNumOfField(`{"targetFieldName":"numOfFields"}`), jp.errorHandler("countNumOfField"))

	// Add some Final stage to process
	// Not mandatory, just for visual control
	jp.jsonTransformPipeline.AddStage(jp.stubFinal(), jp.errorHandler("stubFinal"))
}

func (jp *Process) Run() {
	inStream := make(chan interface{})
	jp.doneChannel = make(chan interface{})
	jp.inStream = inStream
	jp.activeStatus.Store(uint8(1))
	jp.processWg = &sync.WaitGroup{}

	jp.processWg.Add(1)

	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		defer close(inStream)

		log.Println("JsonProcess RUN")

		processDone := jp.jsonTransformPipeline.RunPlug(jp.doneChannel, inStream)

		<-processDone
	}(jp.processWg)
}

func (jp *Process) Stop() {
	if jp.isActive() {
		jp.setActiveOff()
		close(jp.doneChannel)
		jp.processWg.Wait()
		log.Println("JsonProcess STOP")
	} else {
		log.Println("JsonProcess already not active")
	}
}

func (jp *Process) isActive() bool {
	return jp.activeStatus.Load() != nil && jp.activeStatus.Load().(uint8) > 0
}

func (jp *Process) setActiveOff() {
	jp.activeStatus.Store(uint8(0))
}

func (jp *Process) InStream() pipeline.WriteOnlyStream {
	if !jp.isActive() {
		log.Fatal("JsonProcess not active")
	}
	return jp.inStream
}

// ProcessFactory create and init new stage in process
func (jp *Process) ProcessFactory(processName string, processConf string) pipeline.ProcessFn {
	switch strings.ToLower(processName) {
	case "removefield":
		return jp.removeField(processConf)
	case "countnumoffields":
		return jp.countNumOfField(processConf)
	case "addfield":
		return jp.addField(processConf)
	default:
		log.Panic("Unknown Process")
	}

	return nil
}

// Process methods -----------------------------------------------------------------------------------------------
func (jp *Process) addField(configuration string) pipeline.ProcessFn {
	var processName = "addField"
	var config addConfiguration

	if err := json.Unmarshal([]byte(configuration), &config); err != nil {
		panic(err)
	}
	log.Printf("Process %s params: %+v", processName, config)

	return func(inObj interface{}) (interface{}, error) {

		var dat map[string]interface{}
		if err := json.Unmarshal([]byte(fmt.Sprintf("%v", inObj)), &dat); err != nil {
			panic(err)
		}

		dat[config.FieldName] = config.FieldValue
		result, _ := json.Marshal(dat)

		//log.Printf("JsonProcess.%s data: %v, configuration: %+v", processName, inObj, config)
		log.Printf("JsonProcess.%s data result: %s", processName, string(result))

		time.Sleep(time.Millisecond * sleepTimeMs)
		return string(result), nil
	}
}

func (jp *Process) removeField(configuration string) pipeline.ProcessFn {
	var processName = "removeField"
	var config removeFieldConfiguration

	if err := json.Unmarshal([]byte(configuration), &config); err != nil {
		panic(err)
	}

	log.Printf("Process %s params: %+v", processName, config)

	return func(inObj interface{}) (interface{}, error) {

		var dat map[string]interface{}
		if err := json.Unmarshal([]byte(fmt.Sprintf("%v", inObj)), &dat); err != nil {
			panic(err)
		}

		if _, ok := dat[config.FieldName]; ok {
			delete(dat, config.FieldName)
		}

		result, _ := json.Marshal(dat)

		//log.Printf("JsonProcess.%s data: %v, configuration: %+v", processName, inObj, config)
		log.Printf("JsonProcess.%s data result: %s", processName, string(result))

		time.Sleep(time.Millisecond * sleepTimeMs)
		return string(result), nil
	}
}

func (jp *Process) countNumOfField(configuration string) pipeline.ProcessFn {
	var processName = "countNumOfFields"
	var config countNumOfFieldConfiguration

	if err := json.Unmarshal([]byte(configuration), &config); err != nil {
		panic(err)
	}

	log.Printf("Process %s params: %+v", processName, config)

	return func(inObj interface{}) (interface{}, error) {

		var dat map[string]interface{}
		if err := json.Unmarshal([]byte(fmt.Sprintf("%v", inObj)), &dat); err != nil {
			panic(err)
		}

		if _, ok := dat[config.TargetFieldName]; ok {
			delete(dat, config.TargetFieldName)
		}

		dat[config.TargetFieldName] = len(dat)
		result, _ := json.Marshal(dat)

		//log.Printf("JsonProcess.%s data: %v, configuration: %+v", processName, inObj, config)
		log.Printf("JsonProcess.%s data result: %s", processName, string(result))

		time.Sleep(time.Millisecond * sleepTimeMs)
		return string(result), nil
	}
}

func (jp *Process) stubStart() pipeline.ProcessFn {
	return func(inObj interface{}) (interface{}, error) {
		log.Printf("JsonProcess Start data: %v", inObj)
		time.Sleep(time.Millisecond * sleepTimeMs)
		return inObj, nil
	}
}

func (jp *Process) stubFinal() pipeline.ProcessFn {
	return func(inObj interface{}) (interface{}, error) {
		log.Printf("JsonProcess Final data: %v", inObj)
		time.Sleep(time.Millisecond * sleepTimeMs)
		return inObj, nil
	}
}

func (jp *Process) errorHandler(stgName string) pipeline.ErrorProcessFn {
	return func(err error) {
		errorFrom := fmt.Sprintf("stg %s error: %v", stgName, err)
		log.Printf(errorFrom)
	}
}
