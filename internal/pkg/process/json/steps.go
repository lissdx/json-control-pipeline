package json

import (
	"encoding/json"
	"log"
)

// ProcessSteps struct to keep
// array of steps. Implements Iterator pattern
type ProcessSteps struct {
	processes []ProcessStep
}

// NewProcessSteps is ProcessSteps constructor
func NewProcessSteps() *ProcessSteps {
	return &ProcessSteps{}
}

// Init takes JSON steps conf and transforms it
// to process steps obj.
func (ps *ProcessSteps) Init(steps string) {

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

// HasNext Iterator pattern implementation
// Return true if there are next node in ProcessSteps obj
func (ps *ProcessSteps) HasNext() bool {
	return len(ps.processes) > 0
}

// Next Iterator pattern implementation
// Return next node of ProcessSteps obj and delete it
func (ps *ProcessSteps) Next() ProcessStep {
	if !ps.HasNext() {
		//log.Printf("Empty Steps")
		log.Panic("Empty Steps")
	}
	defer func() {
		ps.processes = ps.processes[1:]
	}()
	return ps.processes[0]
}
