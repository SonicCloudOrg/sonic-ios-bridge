package entity

import (
	"encoding/json"
	"log"
)

type PerfData struct {
	PerfDataBytes []byte
}

func (p *PerfData) ToString() string {
	return string(p.PerfDataBytes)
}

func (p *PerfData) ToJson() string {
	return string(p.PerfDataBytes)
}

func (p *PerfData) ToFormat() string {
	data := make(map[string]interface{})
	err := json.Unmarshal(p.PerfDataBytes, &data)
	if err != nil {
		log.Println(err)
	}
	result, _ := json.MarshalIndent(data, "", "\t")
	return string(result)
}
