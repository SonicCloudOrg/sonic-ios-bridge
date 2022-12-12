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
)

type NetworkInfo struct {
	Mac  string `json:"mac"`
	IPv4 string `json:"ipv4"`
	IPv6 string `json:"ipv6"`
}

func (networkInfo NetworkInfo) ToString() string {
	return fmt.Sprintf("%s %s %s", networkInfo.Mac, networkInfo.IPv4, networkInfo.IPv6)
}

func (networkInfo NetworkInfo) ToJson() string {
	result, _ := json.Marshal(networkInfo)
	return string(result)
}

func (networkInfo NetworkInfo) ToFormat() string {
	result, _ := json.MarshalIndent(networkInfo, "", "\t")
	return string(result)
}
