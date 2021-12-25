package main

import (
	"github.com/lissdx/json-control-pipeline/internal/pkg/process"
	"github.com/lissdx/json-control-pipeline/internal/pkg/process/json"
	"github.com/lissdx/json-control-pipeline/pkg/pipeline"
	"log"
	"time"
)

var steps process.Conf = `{
  "steps":[
    {
      "processor":"AddField",
      "configuration":{
        "fieldName":"firstName",
        "fieldValue":"George"
      }
    },
    {
      "processor":"AddField",
      "configuration":{
        "fieldName":"unexpected",
        "fieldValue":"unknown"
      }
    },
    {
      "processor":"CountNumOfFields",
      "configuration":{
        "targetFieldName":"numOfFields"
      }
    },
    {
      "processor":"RemoveField",
      "configuration":{
        "fieldName":"age"
      }
    },
    {
      "processor":"CountNumOfFields",
      "configuration":{
        "targetFieldName":"numOfFields"
      }
    }
  ]
}`

func main() {
	// Create control channel
	done := make(chan interface{})
	defer close(done)

	// Data to change it
	//data := []interface{}{`{"age": 1}`,`{"age":2}`,`{"age":3}`,`{"age":4}`,`{"age":5, "someString": "string"}`}
	//data := []interface{}{`{"NUM": 1, "age": 1, "someString": "string", "someField":"someValue"}`, `{"NUM": 2, "age":5, "someString": "string"}`, `{"NUM": 3}`, `{"NUM": 4}`}
	data := []interface{}{`{"age": 1, "someString": "string", "someField":"someValue"}`}

	// Create new JSON Process line
	jsonProcess := json.NewJsonProcess()

	// Init process with  configuration
	jsonProcess.InitProcess(steps)

	// Run pipeline
	jsonProcess.Run()

	//inStream := jsonProcess.InStream()

	// Send data to process line
	// It could be never ended story
	for v := range pipeline.Generator(done, data...) {
		jsonProcess.InStream() <- v
	}

	time.Sleep(time.Second * 5)
	jsonProcess.Stop()

	log.Println("Main process END")
}
