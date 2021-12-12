package main

import (
	"encoding/json"
	"log"
)

var steps = `{
  "steps":[
    {
      "processor":"AddField",
      "configuration":{
        "fieldName":"firstName",
        "fieldValue":"George"
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

type T struct {
	Steps []struct {
		Processor     string `json:"processor"`
		Configuration string `json:"configuration"`
	} `json:"steps"`
}

//type ProcessSteps struct {
//	Steps []struct {
//		Processor     string `json:"processor"`
//		Configuration map[string]string `json:"configuration"`
//	} `json:"steps"`
//}
func main() {
	t := T{}
 	if err := json.Unmarshal([]byte(steps), &t); err != nil {
		log.Printf("%v", err)
	}

	log.Printf("%v", t)
}
