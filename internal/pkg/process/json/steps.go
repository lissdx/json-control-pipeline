package json

import (
	"encoding/json"
	"log"
)

type ProcessSteps struct {
	processes []ProcessStep
}

func NewProcessSteps() *ProcessSteps {
	return &ProcessSteps{}
}

func (ps *ProcessSteps)Init(steps string)  {

	ps.processes = make([]ProcessStep, 0, 10)
	t := UnmarshalT{}

	if err := json.Unmarshal([]byte(steps), &t); err != nil {
		log.Printf("%v", err)
	}

	for _, v := range t.Steps {
		newStep := ProcessStep{}
		newStep.ProcessName = v.Processor
		res, _ := json.Marshal(v.Configuration)
		newStep.ProcessConf = string(res)
		ps.processes = append(ps.processes, newStep)
	}
}

func (ps *ProcessSteps)HasNext() bool {
	return len(ps.processes) > 0
}

func (ps *ProcessSteps)GetNext() ProcessStep   {
	if !ps.HasNext(){
		//log.Printf("Empty Steps")
		log.Panic("Empty Steps")
	}
	defer func() {
		ps.processes = ps.processes[1:]
	}()
	return ps.processes[0]
}
