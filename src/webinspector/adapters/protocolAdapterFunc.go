package adapters

import (
	"encoding/json"
	"fmt"
	"github.com/SonicCloudOrg/sonic-ios-bridge/src/entity"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"log"
	"regexp"
	"strings"
)

type ProtocolAdapter struct {
	adapter                    *Adapter
	lastNodeId                 int64
	lastPageExecutionContextId int64
	styMap                     map[string]interface{}
	lastScriptEval             interface{}
}

type MessageAdapters func(message []byte) []byte

var PageSetOverlay MessageAdapters = func(message []byte) []byte {
	method := "Debugger.setOverlayMessage"
	return ReplaceMethodNameAndOutputBinary(message, method)
}

var PageConfigureOverlay MessageAdapters = func(message []byte) []byte {
	return PageSetOverlay(message)
}

var DOMSetInspectedNode MessageAdapters = func(message []byte) []byte {
	method := "Console.addInspectedNode"
	return ReplaceMethodNameAndOutputBinary(message, method)
}

var EmulationSetTouchEmulationEnabled MessageAdapters = func(message []byte) []byte {
	method := "Page.setTouchEmulationEnabled"
	return ReplaceMethodNameAndOutputBinary(message, method)
}

var EmulationSetScriptExecutionDisabled MessageAdapters = func(message []byte) []byte {
	method := "Page.setScriptExecutionDisabled"
	return ReplaceMethodNameAndOutputBinary(message, method)
}

var EmulationSetEmulatedMedia MessageAdapters = func(message []byte) []byte {
	method := "Page.setEmulatedMedia"
	return ReplaceMethodNameAndOutputBinary(message, method)
}

var RenderingSetShowPaintRects MessageAdapters = func(message []byte) []byte {
	method := "Page.setShowPaintRects"
	return ReplaceMethodNameAndOutputBinary(message, method)
}

var LogClear MessageAdapters = func(message []byte) []byte {
	method := "Console.clearMessages"
	return ReplaceMethodNameAndOutputBinary(message, method)
}

var LogDisable MessageAdapters = func(message []byte) []byte {
	method := "Console.disable"
	return ReplaceMethodNameAndOutputBinary(message, method)
}

var LogEnable MessageAdapters = func(message []byte) []byte {
	method := "Console.enable"
	return ReplaceMethodNameAndOutputBinary(message, method)
}
var NetworkGetCookies MessageAdapters = func(message []byte) []byte {
	method := "Page.getCookies"
	return ReplaceMethodNameAndOutputBinary(message, method)
}

var NetworkDeleteCookie MessageAdapters = func(message []byte) []byte {
	method := "Page.deleteCookie"
	return ReplaceMethodNameAndOutputBinary(message, method)
}

var NetworkSetMonitoringXHREnabled MessageAdapters = func(message []byte) []byte {
	method := "Console.setMonitoringXHREnabled"
	return ReplaceMethodNameAndOutputBinary(message, method)
}

func (p *ProtocolAdapter) onGetMatchedStylesForNode(message []byte) []byte {
	p.lastNodeId = gjson.Get(string(message), "params.nodeId").Int()
	return message
}

func (p *ProtocolAdapter) onCanEmulate(message []byte) []byte {
	result := map[string]interface{}{
		"result": true,
	}
	p.adapter.FireResultToTools(int(gjson.Get(string(message), "id").Int()), result)
	return nil
}

func (p *ProtocolAdapter) onGetPlatformFontsForNode(message []byte) []byte {
	result := map[string]interface{}{
		"fonts": []string{},
	}
	p.adapter.FireResultToTools(int(gjson.Get(string(message), "id").Int()), result)
	return nil
}

func (p *ProtocolAdapter) onGetBackgroundColors(message []byte) []byte {
	result := map[string]interface{}{
		"backgroundColors": []string{},
	}
	p.adapter.FireResultToTools(int(gjson.Get(string(message), "id").Int()), result)
	return nil
}

func (p *ProtocolAdapter) onCanSetScriptSource(message []byte) []byte {
	result := map[string]interface{}{
		"result": false,
	}
	p.adapter.FireResultToTools(int(gjson.Get(string(message), "id").Int()), result)
	return nil
}

func (p *ProtocolAdapter) onSetBlackboxPatterns(message []byte) []byte {
	result := map[string]interface{}{}
	p.adapter.FireResultToTools(int(gjson.Get(string(message), "id").Int()), result)
	return nil
}

func (p *ProtocolAdapter) onSetAsyncCallStackDepth(message []byte) []byte {
	result := map[string]interface{}{
		"result": true,
	}
	p.adapter.FireResultToTools(int(gjson.Get(string(message), "id").Int()), result)
	return nil
}

func (p *ProtocolAdapter) onDebuggerEnable(message []byte) []byte {
	p.adapter.CallTarget("Debugger.setBreakpointsActive", map[string]interface{}{
		"active": true,
	}, p.defaultCallFunc)
	return message
}

func (p *ProtocolAdapter) onGetMatchedStylesForNodeResult(message []byte) []byte {
	msg := string(message)
	result := gjson.Get(msg, "result")
	if result.Exists() {
		for _, matchedCSSRule := range result.Get("matchedCSSRules").Array() {
			p.mapRule(matchedCSSRule.Get("rule"), &msg)
		}
		for _, inherited := range result.Get("inherited").Array() {
			if inherited.Get("matchedCSSRules").Exists() {
				for _, matchedCSSRule := range result.Get("matchedCSSRules").Array() {
					p.mapRule(matchedCSSRule.Get("rule"), &msg)
				}
			}
		}
	}
	return []byte(msg)
}

func (p *ProtocolAdapter) onEvaluate(message []byte) []byte {
	msg := string(message)
	var err error
	result := gjson.Get(msg, "result")
	if result.Exists() && result.Get("wasThrown").Exists() {
		msg, err = sjson.Set(msg, "result.result.subtype", "error")
		if err != nil {
			return nil
		}
		arr, err1 := json.Marshal(map[string]interface{}{
			"text":     gjson.Get(msg, "result.result.description").Value(),
			"url":      "",
			"scriptId": p.lastScriptEval,
			"line":     1,
			"column":   0,
			"stack": map[string]interface{}{
				"callFrames": []map[string]interface{}{
					{
						"functionName": "",
						"scriptId":     p.lastScriptEval,
						"url":          "",
						"lineNumber":   1,
						"columnNumber": 1,
					},
				},
			},
		})
		if err1 != nil {
			log.Panic(err)
		}
		msg, err = sjson.Set(msg, "result.exceptionDetails", string(arr))
		if err != nil {
			log.Panic(err)
		}
	} else if result.Exists() && result.Get("result").Exists() && result.Get("result.preview").Exists() {
		msg, err = sjson.Set(msg, "result.result.preview.description", gjson.Get(msg, "result.result.description").Value())
		if err != nil {
			log.Panic(err)
		}
		msg, err = sjson.Set(msg, "result.result.preview.type", "object")
		if err != nil {
			log.Panic(err)
		}
	}
	return []byte(msg)
}

func (p *ProtocolAdapter) onRuntimeOnCompileScript(message []byte) []byte {
	params := map[string]interface{}{
		"expression": gjson.Get(string(message), "params.expression").String(),
		"contextId":  gjson.Get(string(message), "params.executionContextId").Int(),
	}
	p.adapter.CallTarget("Runtime.evaluate", params, func(message []byte) {
		result := map[string]interface{}{
			"scriptId":         nil,
			"exceptionDetails": nil,
		}
		p.adapter.FireResultToTools(int(gjson.Get(string(message), "id").Int()), result)
	})
	return nil
}

func (p *ProtocolAdapter) onExecutionContextCreated(message []byte) []byte {
	msg := string(message)
	var err error
	if gjson.Get(msg, "params").Exists() && gjson.Get(msg, "params.context").Exists() {
		if !gjson.Get(msg, "params.context.origin").Exists() {
			msg, err = sjson.Set(msg, "params.context.origin", gjson.Get(msg, "params.context.name").String())
			if err != nil {
				log.Panic(err)
			}
			if gjson.Get(msg, "params.context.isPageContext").Exists() {
				p.lastPageExecutionContextId = gjson.Get(msg, "params.context.id").Int()
			}
			if gjson.Get(msg, "params.context.frameId").Exists() {
				msg, err = sjson.Set(msg, "params.context.auxData", map[string]interface{}{
					"frameId":   gjson.Get(msg, "params.context.frameId").String(),
					"isDefault": true,
				})
				if err != nil {
					log.Panic(err)
				}
				msg, err = sjson.Delete(msg, "params.context.frameId")
				if err != nil {
					log.Panic(err)
				}
			}
		}
	}

	return []byte(msg)
}

func (p *ProtocolAdapter) defaultCallFunc(message []byte) {
	log.Println(string(message))
}

func (p *ProtocolAdapter) onAddRule(message []byte) []byte {
	var selector = gjson.Get(string(message), "params.ruleText").String()
	selector = strings.TrimSpace(selector)
	selector = strings.Replace(selector, "{}", "", -1)
	params := map[string]interface{}{
		"contextNodeId": p.lastNodeId,
		"selector":      selector,
	}
	p.adapter.CallTarget("CSS.addRule", params, func(message []byte) {
		var msg = string(message)
		var param interface{}
		err := json.Unmarshal(message, param)
		if err != nil {
			log.Panic(err)
		}
		p.mapRule(gjson.Get(msg, "rule"), &msg)
		p.adapter.FireResultToTools(int(gjson.Get(msg, "id").Int()), param)
	})
	return nil
}

func (p *ProtocolAdapter) mapRule(cssRule gjson.Result, message *string) {
	var err error
	if cssRule.Get("ruleId").Exists() {
		path := cssRule.Get("styleSheetId").Path(*message)
		*message, err = sjson.Set(*message, path, cssRule.Get("ruleId.styleSheetId").Value())
		if err != nil {
			log.Panic(err)
		}
		*message, err = sjson.Delete(*message, path)
		if err != nil {
			log.Panic(err)
		}
		// todo
		//p.mapSelectorList(nil)

		p.mapStyle(cssRule.Get("style"), cssRule.Get("origin").String(), message)

		path = cssRule.Get("sourceLine").Path(*message)
		*message, err = sjson.Delete(*message, path)
		if err != nil {
			log.Panic(err)
		}
	}
}

func (p *ProtocolAdapter) mapSelectorList(result gjson.Result) {

}

// todo 完善CSS，真他妈畜生这块
// CSSSetStyleTexts todo KeyCheck
func (p *ProtocolAdapter) CSSSetStyleTexts(message []byte) []byte {
	var msg = string(message)
	var allStyleText []interface{}
	resultId := gjson.Get(msg, "id").Int()
	editsResult := gjson.Get(msg, "params.edits").Array()
	var whetherToContinueTheCycle = true

	for _, edit := range editsResult {
		if !whetherToContinueTheCycle {
			break
		}
		paramsGetStyleSheet := map[string]interface{}{
			"styleSheetId": edit.Get("styleSheetId").String(),
		}
		p.adapter.CallTarget("CSS.getStyleSheet", paramsGetStyleSheet, func(message []byte) {
			msg = string(message)
			styleSheet := gjson.Get(msg, "styleSheet")
			styleSheetRules := gjson.Get(msg, "styleSheet.rules")
			if !styleSheet.Exists() || !styleSheetRules.Exists() {
				log.Panic("iOS returned a value we were not expecting for getStyleSheet")
			}
			for index, rule := range styleSheetRules.Array() {
				if compareRanges(rule.Get("style.range"), edit.Get("range")) {
					params := map[string]interface{}{
						"styleId": map[interface{}]interface{}{
							"styleSheetId": edit.Get("styleSheetId").String(),
							"ordinal":      index,
						},
						"text": edit.Get("text").String(),
					}
					p.adapter.CallTarget("CSS.allStyleText", params, func(message []byte) {
						msg = string(message)
						p.mapStyle(gjson.Get(string(message), "style"), "", &msg)
						allStyleText = append(allStyleText, gjson.Get(msg, "style").Value())
						whetherToContinueTheCycle = false
					})
				}
			}
		})
	}
	result := map[string]interface{}{
		"styles": allStyleText,
	}
	p.adapter.FireResultToTools(int(resultId), result)
	return nil
}

func (p *ProtocolAdapter) mapStyle(cssStyle gjson.Result, ruleOrigin string, message *string) {
	var err error
	if cssStyle.Get("cssText").Exists() {
		disabled := p.extractDisabledStyles(cssStyle.Get("cssText").String(), cssStyle.Get("range"))
		for i, value := range disabled {
			noSpaceStr := strings.TrimSpace(value.Content)
			// 原版 const text = disabled[i].content.trim().replace(/^\/\*\s*/, '').replace(/;\s*\*\/$/, '');
			reg := regexp.MustCompile(`^\\/\\*\\s*`)
			noSpaceStr = reg.ReplaceAllString(noSpaceStr, ``)

			reg = regexp.MustCompile(`;\\s*\\*\\/$`)
			noSpaceStr = reg.ReplaceAllString(noSpaceStr, ``)

			parts := strings.Split(noSpaceStr, ":")
			if cssStyle.Get("cssProperties").Exists() {
				cssProperties := cssStyle.Get("cssProperties").Array()
				var index = len(cssProperties)
				for j, _ := range cssProperties {
					if cssProperties[j].Get("range").Exists() &&
						(cssProperties[j].Get("range.startLine").Int() > int64(disabled[i].CssRange.StartLine) ||
							cssProperties[j].Get("range.startLine").Int() == int64(disabled[i].CssRange.StartLine) ||
							cssProperties[j].Get("range.startColumn").Int() > int64(disabled[i].CssRange.StartColumn)) {
						index = j
						break
					}
				}

				cssPropertiesObjects := cssStyle.Get("cssProperties").Value()
				path := cssStyle.Get("cssProperties").Path(*message)
				if cssPropertiesArrays, ok := cssPropertiesObjects.([]interface{}); ok {
					var cssPropertiesFinal []interface{}
					cssPropertiesLeft := cssPropertiesArrays[:index+1]
					cssPropertiesRight := cssPropertiesArrays[index+1:]

					cssPropertiesFinal = append(cssPropertiesLeft, map[string]interface{}{
						"implicit": false,
						"name":     parts[0],
						"range":    disabled[i].CssRange,
						"status":   "disabled",
						"text":     disabled[i].Content,
						"value":    parts[1],
					})
					cssPropertiesFinal = append(cssPropertiesFinal, cssPropertiesRight...)
					arr, err1 := json.Marshal(cssPropertiesFinal)
					if err1 != nil {
						log.Panic(err1)
					}
					*message, err = sjson.Set(*message, path, string(arr))
					if err != nil {
						log.Panic(err)
					}
				} else {
					log.Panic(fmt.Errorf("failed to convert object"))
				}

			}

		}
	}

	for _, cssProperty := range cssStyle.Get("cssProperties").Array() {
		p.mapCssProperty(cssProperty, message)
	}
	if ruleOrigin != "user-agent" {
		path := cssStyle.Get("styleSheetId").Path(*message)
		*message, err = sjson.Set(*message, path, cssStyle.Get("styleId.styleSheetId").String())
		if err != nil {
			log.Panic(err)
		}
		cssStyleRangeArr, err1 := json.Marshal(cssStyle.Get("range").Value())
		if err1 != nil {
			log.Panic(err1)
		}
		var styleKey = fmt.Sprintf("%s_%s", cssStyle.Get("styleId.styleSheetId").String(), string(cssStyleRangeArr))
		if p.styMap == nil {
			p.styMap = make(map[string]interface{})
			p.styMap[styleKey] = cssStyle.Get("styleId.styleSheetId").String()
		}
		// delete
		path = cssStyle.Get("styleId").Path(*message)
		*message, err = sjson.Delete(*message, path)
		if err != nil {
			log.Panic(err)
		}
		path = cssStyle.Get("sourceLine").Path(*message)
		*message, err = sjson.Delete(*message, path)
		if err != nil {
			log.Panic(err)
		}
		path = cssStyle.Get("sourceURL").Path(*message)
		*message, err = sjson.Delete(*message, path)
		if err != nil {
			log.Panic(err)
		}
		path = cssStyle.Get("width").Path(*message)
		*message, err = sjson.Delete(*message, path)
		if err != nil {
			log.Panic(err)
		}
		path = cssStyle.Get("height").Path(*message)
		*message, err = sjson.Delete(*message, path)
		if err != nil {
			log.Panic(err)
		}
	}
}

func (p *ProtocolAdapter) mapCssProperty(cssProperty gjson.Result, message *string) {
	var err error
	path := cssProperty.Get("status.disabled").Path(*message)
	if cssProperty.Get("status").String() == "disabled" {
		*message, err = sjson.Set(*message, path, true)
		if err != nil {
			log.Panic(err)
		}
	} else if cssProperty.Get("status").String() == "active" {
		*message, err = sjson.Set(*message, path, false)
		if err != nil {
			log.Panic(err)
		}
	}
	path = cssProperty.Get("status").Path(*message)
	*message, err = sjson.Delete(*message, path)
	if err != nil {
		log.Panic(err)
	}
	priority := cssProperty.Get("priority")

	path = cssProperty.Path(*message) + ".important"

	if cssProperty.Get("priority").Exists() {
		if priority.String() == "" {
			*message, err = sjson.Set(*message, path, false)
			if err != nil {
				log.Panic(err)
			}
		} else {
			*message, err = sjson.Set(*message, path, true)
			if err != nil {
				log.Panic(err)
			}
		}
	} else {
		*message, err = sjson.Set(*message, path, false)
		if err != nil {
			log.Panic(err)
		}
	}

	path = priority.Path(*message)
	*message, err = sjson.Delete(*message, path)
	if err != nil {
		log.Panic(err)
	}
}

// extractDisabledStyles todo KeyCheck
func (p *ProtocolAdapter) extractDisabledStyles(styleText string, cssRange gjson.Result) []entity.IDisabledStyle {
	var startIndices []int
	var styles []entity.IDisabledStyle
	for index, _ := range styleText {
		endIndexBEGINCOMMENT := index + len(BEGIN_COMMENT)
		endIndexENDCOMMENT := index + len(END_COMMENT)
		if endIndexBEGINCOMMENT <= len(styleText) && string([]rune(styleText)[index:endIndexBEGINCOMMENT]) == BEGIN_COMMENT {
			startIndices = append(startIndices, index)
			index = index + len(BEGIN_COMMENT)
		} else if endIndexENDCOMMENT <= len(styleText) && string([]rune(styleText)[index:endIndexENDCOMMENT]) == END_COMMENT {
			if len(startIndices) == 0 {
				return nil
			}
			startIndex := startIndices[0]
			startIndices = startIndices[1:]
			endIndex := index + len(END_COMMENT)

			startRangeLine, startRangeColumn := p.getLineColumnFromIndex(styleText, startIndex, cssRange)
			endRangeLine, endRangeColumn := p.getLineColumnFromIndex(styleText, endIndex, cssRange)

			propertyItem := entity.IDisabledStyle{
				Content: styleText[startIndex:endIndex],
				CssRange: entity.IRange{
					StartLine:   startRangeLine,
					StartColumn: startRangeColumn,
					EndLine:     endRangeLine,
					EndColumn:   endRangeColumn,
				},
			}
			styles = append(styles, propertyItem)
			index = endIndex - 1
		}
	}
	if len(startIndices) == 0 {
		return nil
	}
	return styles
}

// todo KeyCheck
func (p *ProtocolAdapter) getLineColumnFromIndex(text string, index int, startRange gjson.Result) (line int, column int) {
	if text == "" || index < 0 || index > len(text) {
		return 0, 0
	}
	if startRange.Exists() {
		line = int(startRange.Get("StartLine").Int())
		column = int(startRange.Get("StartColumn").Int())
	}
	for i := 0; i < len(text) && i < index; i++ {
		if text[i] == '\r' && i+1 < len(text) && text[i+1] == '\n' {
			i++
			line++
			column = 0
		} else if text[i] == '\n' || text[i] == '\r' {
			line++
			column = 0
		} else {
			column++
		}
	}
	return line, column
}

func compareRanges(rangeLeft gjson.Result, rangeRight gjson.Result) bool {
	return rangeLeft.Get("startLine").Int() == rangeRight.Get("startLine").Int() &&
		rangeLeft.Get("startColumn").Int() == rangeRight.Get("startColumn").Int() &&
		rangeLeft.Get("endLine").Int() == rangeRight.Get("endLine").Int() &&
		rangeLeft.Get("endColumn").Int() == rangeRight.Get("endColumn").Int()
}

func ReplaceMethodNameAndOutputBinary(message []byte, method string) []byte {
	var msg = make(map[string]interface{})
	err := json.Unmarshal(message, &msg)
	if err != nil {
		log.Panic(err)
	}
	// todo Regular?
	msg["method"] = method

	arr, err1 := json.Marshal(message)
	if err1 != nil {
		log.Panic(err1)
	}
	return arr
}

var BEGIN_COMMENT = "/* "
var END_COMMENT = " */"
var SEPARATOR = ": "
