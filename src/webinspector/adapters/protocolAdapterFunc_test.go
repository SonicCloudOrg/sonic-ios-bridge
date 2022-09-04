package adapters

import (
	"encoding/json"
	"fmt"
	"github.com/SonicCloudOrg/sonic-ios-bridge/src/entity"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"log"
	"testing"
)

var protocolAdapter *ProtocolAdapter

func setProtocolAdapterFunc() {
	protocolAdapter = &ProtocolAdapter{}
}

func TestOnExecutionContextCreated(t *testing.T) {
	setProtocolAdapterFunc()
	simulationInformation := map[string]interface{}{
		"params": map[string]interface{}{
			"context": map[string]interface{}{
				"id":      2,
				"type":    "normal",
				"name":    "test",
				"frameId": "testFrameId",
			},
		},
	}
	arr, _ := json.Marshal(simulationInformation)
	log.Println(string(protocolAdapter.onExecutionContextCreated(arr)))
}

func TestSJsonArrayInsertion(t *testing.T) {
	var cssProperties []interface{}
	cssProperties = append(cssProperties, map[string]interface{}{
		"implicit": false,
		"name":     "cccc",
		"range": entity.IRange{
			StartLine:   6,
			StartColumn: 7,
			EndLine:     2,
			EndColumn:   2,
		},
		"status": "disabled",
		"text":   "wwwwwwwww",
		"value":  "cxccccc",
	})
	cssProperties = append(cssProperties, map[string]interface{}{
		"implicit": false,
		"name":     "aaaa",
		"range": entity.IRange{
			StartLine:   6,
			StartColumn: 7,
			EndLine:     2,
			EndColumn:   2,
		},
		"status": "aaa",
		"text":   "aaa",
		"value":  "aaa",
	})
	arr, _ := json.Marshal(map[string]interface{}{
		"cssProperties": cssProperties,
	})
	cssPropertiesObjects := gjson.Get(string(arr), "cssProperties").Value()
	var index = 1
	if cssPropertiesArrays, ok := cssPropertiesObjects.([]interface{}); ok {
		var cssPropertiesFinal []interface{}
		cssPropertiesLeft := cssPropertiesArrays[:index+1]
		cssPropertiesRight := cssPropertiesArrays[index+1:]

		//cssPropertiesFinal = append(cssPropertiesFinal,cssPropertiesLeft...)
		cssPropertiesFinal = append(cssPropertiesLeft, map[string]interface{}{
			"implicit": false,
			"name":     "parts[0]",
			"range":    "disabled[i].CssRange",
			"status":   "disabled",
			"text":     "disabled[i].Content",
			"value":    "parts[1]",
		})
		cssPropertiesFinal = append(cssPropertiesFinal, cssPropertiesRight...)
		arr1, err := json.Marshal(cssPropertiesFinal)
		if err != nil {
			log.Panic(err)
		}
		value, err := sjson.Set(string(arr), "cssProperties", string(arr1))
		if err != nil {
			log.Panic(err)
		}
		log.Println(value)
	} else {
		log.Panic(fmt.Errorf("failed to convert object"))
	}
}
