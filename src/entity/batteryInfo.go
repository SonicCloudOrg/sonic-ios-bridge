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
	"fmt"
	"strings"
)

type BatteryInter struct {
	BatteryCurrentCapacity int
}

type Battery struct {
	SerialNumber string `json:"serialNumber,omitempty"`
	Level        int    `json:"level"`
	Temperature  int    `json:"temperature"`
}

type BatteryList struct {
	BatteryInfo []Battery `json:"batteryList"`
}

func (batteryList BatteryList) ToString() string {
	var s strings.Builder
	for i, e := range batteryList.BatteryInfo {
		if i != len(batteryList.BatteryInfo)-1 {
			s.WriteString(fmt.Sprintf("%s %d %d\n", e.SerialNumber, e.Level, e.Temperature))
		} else {
			s.WriteString(fmt.Sprintf("%s %d %d", e.SerialNumber, e.Level, e.Temperature))
		}
	}
	return s.String()
}

func (batteryList BatteryList) ToJson() string {
	result, _ := json.Marshal(batteryList)
	return string(result)
}

func (batteryList BatteryList) ToFormat() string {
	result, _ := json.MarshalIndent(batteryList, "", "\t")
	return string(result)
}
