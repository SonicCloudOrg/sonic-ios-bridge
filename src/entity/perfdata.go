/*
 *  Copyright (C) [SonicCloudOrg] Sonic Project
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *         http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 *
 */
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
