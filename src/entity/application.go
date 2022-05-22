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

type Application struct {
	CFBundleVersion     string `json:"version"`
	CFBundleDisplayName string `json:"name"`
	CFBundleIdentifier  string `json:"bundleId"`
}

type AppList struct {
	ApplicationList []Application `json:"appList"`
}

func (appList AppList) ToString() string {
	var s strings.Builder
	for i, e := range appList.ApplicationList {
		if i != len(appList.ApplicationList)-1 {
			s.WriteString(fmt.Sprintf("%s %s %s\n", e.CFBundleDisplayName, e.CFBundleIdentifier, e.CFBundleVersion))
		} else {
			s.WriteString(fmt.Sprintf("%s %s %s", e.CFBundleDisplayName, e.CFBundleIdentifier, e.CFBundleVersion))
		}
	}
	return s.String()
}

func (appList AppList) ToJson() string {
	result, _ := json.Marshal(appList)
	return string(result)
}

func (appList AppList) ToFormat() string {
	result, _ := json.MarshalIndent(appList, "", "\t")
	return string(result)
}
