/*
 *   sonic-ios-bridge  Connect to your iOS Devices.
 *   Copyright (C) 2022 SonicCloudOrg
 *
 *   This program is free software: you can redistribute it and/or modify
 *   it under the terms of the GNU General Public License as published by
 *   the Free Software Foundation, either version 3 of the License, or
 *   (at your option) any later version.
 *
 *   This program is distributed in the hope that it will be useful,
 *   but WITHOUT ANY WARRANTY; without even the implied warranty of
 *   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 *   GNU General Public License for more details.
 *
 *   You should have received a copy of the GNU General Public License
 *   along with this program.  If not, see <https://www.gnu.org/licenses/>.
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

func (battery Battery) ToString() string {
	var s strings.Builder
	s.WriteString(fmt.Sprintf("%s %d %d", battery.SerialNumber, battery.Level, battery.Temperature))
	return s.String()
}

func (battery Battery) ToJson() string {
	result, _ := json.Marshal(battery)
	return string(result)
}

func (battery Battery) ToFormat() string {
	result, _ := json.MarshalIndent(battery, "", "\t")
	return string(result)
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
