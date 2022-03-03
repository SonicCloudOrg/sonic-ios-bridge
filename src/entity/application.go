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
