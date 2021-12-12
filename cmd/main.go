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
	done := make(chan interface{})
	defer close(done)
	//data := []interface{}{`{"age": 1}`,`{"age":2}`,`{"age":3}`,`{"age":4}`,`{"age":5, "someString": "string"}`}
	//data := []interface{}{`{"NUM": 1, "age": 1, "someString": "string", "someField":"someValue"}`, `{"NUM": 2, "age":5, "someString": "string"}`, `{"NUM": 3}`, `{"NUM": 4}`}
	data := []interface{}{`{"age": 1, "someString": "string", "someField":"someValue"}`}

	jsonProcess := json.NewJsonProcess()

	jsonProcess.InitProcess(steps)

	jsonProcess.Run()

	//inStream := jsonProcess.InStream()

	for v := range pipeline.Generator(done, data...){
		jsonProcess.InStream() <- v
	}

	time.Sleep(time.Second * 5)
	jsonProcess.Stop()

	log.Println("Main process END")
}
