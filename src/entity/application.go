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

type Application struct {
	CFBundleShortVersionString string `json:"shortVersion,omitempty"`
	CFBundleVersion            string `json:"version"`
	CFBundleDisplayName        string `json:"name"`
	CFBundleIdentifier         string `json:"bundleId"`
	IconBase64                 string `json:"iconBase64,omitempty"`
}

type AppList struct {
	ApplicationList []Application `json:"appList"`
}

func (appList AppList) ToString() string {
	var s strings.Builder
	for i, e := range appList.ApplicationList {
		if i != len(appList.ApplicationList)-1 {
			s.WriteString(fmt.Sprintf("%s %s %s %s\n", e.CFBundleDisplayName, e.CFBundleIdentifier, e.CFBundleVersion, e.CFBundleShortVersionString))
		} else {
			s.WriteString(fmt.Sprintf("%s %s %s %s", e.CFBundleDisplayName, e.CFBundleIdentifier, e.CFBundleVersion, e.CFBundleShortVersionString))
		}
	}
	return s.String()
}

func (appList AppList) ToJson() string {
	for _, a := range appList.ApplicationList {
		result, _ := json.Marshal(a)
		fmt.Println(string(result))
	}
	return ""
}

func (appList AppList) ToFormat() string {
	result, _ := json.MarshalIndent(appList, "", "\t")
	return string(result)
}
