/*
 *   sonic-ios-bridge  Connect to your iOS Devices.
 *   Copyright (C) 2022 SonicCloudOrg
 *
 *   This program is free software: you can redistribute it and/or modify
 *   it under the terms of the GNU Affero General Public License as published
 *   by the Free Software Foundation, either version 3 of the License, or
 *   (at your option) any later version.
 *
 *   This program is distributed in the hope that it will be useful,
 *   but WITHOUT ANY WARRANTY; without even the implied warranty of
 *   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 *   GNU Affero General Public License for more details.
 *
 *   You should have received a copy of the GNU Affero General Public License
 *   along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */
package entity

import (
	"encoding/json"
	"fmt"
	"strings"
)

type Battery struct {
	Serial                string      `json:"Serial,omitempty"`
	CurrentCapacity       interface{} `json:"CurrentCapacity,omitempty"`
	CycleCount            interface{} `json:"CycleCount"`
	AbsoluteCapacity      interface{} `json:"AbsoluteCapacity"`
	NominalChargeCapacity interface{} `json:"NominalChargeCapacity"`
	DesignCapacity        interface{} `json:"DesignCapacity"`
	Voltage               interface{} `json:"Voltage"`
	BootVoltage           interface{} `json:"BootVoltage"`
	AdapterDetailsVoltage interface{} `json:"AdapterDetailsVoltage,omitempty"`
	AdapterDetailsWatts   interface{} `json:"AdapterDetailsWatts,omitempty"`
	InstantAmperage       interface{} `json:"InstantAmperage"`
	Temperature           interface{} `json:"Temperature"`
}

func (battery *Battery) AnalyzeBatteryData(batteryData map[string]interface{}) error {
	DiagnosticsData := batteryData["Diagnostics"].(map[string]interface{})
	IORegistryData := DiagnosticsData["IORegistry"].(map[string]interface{})

	AdapterDetailsData := IORegistryData["AdapterDetails"].(map[string]interface{})
	battery.AdapterDetailsVoltage = AdapterDetailsData["Voltage"]
	battery.AdapterDetailsWatts = AdapterDetailsData["Watts"]

	registryDataBytes, err := json.Marshal(IORegistryData)
	if err != nil {
		return err
	}
	err = json.Unmarshal(registryDataBytes, battery)
	if err != nil {
		return err
	}
	return nil
}

func (battery Battery) ToString() string {
	var s strings.Builder
	s.WriteString(fmt.Sprintf("Serial:%s\n", battery.Serial))
	s.WriteString(fmt.Sprintf("Temperature:%v\n", battery.Temperature))
	s.WriteString(fmt.Sprintf("CycleCount:%v\n", battery.CycleCount))

	s.WriteString(fmt.Sprintf("NominalChargeCapacity:%vmAh\n", battery.NominalChargeCapacity))
	s.WriteString(fmt.Sprintf("DesignCapacity:%vmAh\n", battery.DesignCapacity))
	s.WriteString(fmt.Sprintf("AbsoluteCapacity:%vmAh\n", battery.AbsoluteCapacity))
	s.WriteString(fmt.Sprintf("CurrentCapacity:%v\n", battery.CurrentCapacity))

	s.WriteString(fmt.Sprintf("Voltage:%vmV\nBootVoltage:%vmV\n", battery.Voltage, battery.BootVoltage))
	s.WriteString(fmt.Sprintf("InstantAmperage:%vmA\nAdapterDetailsVoltage:%vmV\n", battery.InstantAmperage, battery.AdapterDetailsVoltage))
	s.WriteString(fmt.Sprintf("AdapterDetailsWatts:%vW", battery.AdapterDetailsWatts))
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

type BatteryList struct {
	DeviceBatteryInfo map[string]Battery
}

func (battery *BatteryList) Put(key string, value Battery) {
	if battery.DeviceBatteryInfo == nil {
		battery.DeviceBatteryInfo = make(map[string]Battery)
	}
	battery.DeviceBatteryInfo[key] = value
}

func (battery BatteryList) ToString() string {
	if battery.DeviceBatteryInfo == nil {
		return ""
	}
	for key, e := range battery.DeviceBatteryInfo {
		fmt.Sprintf("udId:%s\n", key)
		fmt.Println(e.ToString())
		fmt.Println("")
	}
	return ""
}

func (battery BatteryList) ToJson() string {
	if battery.DeviceBatteryInfo == nil {
		return ""
	}
	result, _ := json.Marshal(battery.DeviceBatteryInfo)
	return string(result)
}

func (battery BatteryList) ToFormat() string {
	if battery.DeviceBatteryInfo == nil {
		return ""
	}
	result, _ := json.MarshalIndent(battery.DeviceBatteryInfo, "", "\t")
	return string(result)
}
