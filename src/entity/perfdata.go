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
