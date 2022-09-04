package adapters

import (
	"encoding/json"
	"github.com/tidwall/gjson"
	"log"
)

type ProtocolAdapter struct {
	adapter    *Adapter
	lastNodeId int64
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

// todo
//func (p *ProtocolAdapter) onEvaluate(message []byte) []byte {
//	p.adapter.CallTarget("Debugger.setBreakpointsActive",map[string]interface{}{
//		"active":true,
//	}, p.defaultCallFunc)
//	return message
//}

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

func (p *ProtocolAdapter) defaultCallFunc(message []byte) {
	log.Println(string(message))
}

//// todo 真恶心
//func (p *ProtocolAdapter) onAddRule(message []byte) []byte {
//	return nil
//}

// todo 完善CSS，真他妈畜生这块
// CSSSetStyleTexts todo KeyCheck
//func (p *ProtocolAdapter) CSSSetStyleTexts(message []byte) []byte {
//	var msg = string(message)
//	resultId := gjson.Get(msg, "id").Int()
//	editsResult := gjson.Get(msg, "params.edits").Array()
//	for _, edit := range editsResult {
//		paramsGetStyleSheet := map[string]interface{}{
//			"styleSheetId": edit.Get("styleSheetId").String(),
//		}
//		p.adapter.CallTarget("CSS.getStyleSheet", paramsGetStyleSheet, func(message []byte) {
//			msg = string(message)
//			styleSheet := gjson.Get(msg, "styleSheet")
//			styleSheetRules := gjson.Get(msg, "styleSheet.rules")
//			if !styleSheet.Exists() || !styleSheetRules.Exists() {
//				log.Panic("iOS returned a value we were not expecting for getStyleSheet")
//			}
//			for index, rule := range styleSheetRules.Array() {
//				if compareRanges(rule.Get("style.range"), edit.Get("range")) {
//					params := map[string]interface{}{
//						"styleId": map[interface{}]interface{}{
//							"styleSheetId": edit.Get("styleSheetId").String(),
//							"ordinal":      index,
//						},
//						"text": edit.Get("text").String(),
//					}
//					p.adapter.CallTarget("CSS.setStyleText", params, func(message []byte) {
//
//					})
//				}
//			}
//		})
//	}
//
//	//editsResult := message["edit"]
//}

// todo KeyCheck
//func (p *ProtocolAdapter) mapStyle(cssStyle gjson.Result, ruleOrigin string) {
//	if cssStyle.Get("cssText").Exists() {
//		cssRangeStr := cssStyle.Get("range").String()
//		cssRange := &entity.IRange{}
//		err:=json.Unmarshal([]byte(cssRangeStr),cssRange)
//		if err!=nil {
//			log.Panic(err)
//		}
//		disabled := p.extractDisabledStyles(cssStyle.Get("cssText").String(),cssRange)
//		for i,value:=range disabled{
//			noSpaceStr := strings.TrimSpace(value.Content)
//			// 原版 const text = disabled[i].content.trim().replace(/^\/\*\s*/, '').replace(/;\s*\*\/$/, '');
//			reg := regexp.MustCompile(`^\\/\\*\\s*`)
//			noSpaceStr = reg.ReplaceAllString(noSpaceStr, ``)
//
//			reg = regexp.MustCompile(`;\\s*\\*\\/$`)
//			noSpaceStr = reg.ReplaceAllString(noSpaceStr, ``)
//
//			parts := strings.Split(noSpaceStr,":")
//			if cssStyle.Get("cssProperties").Exists() {
//				cssProperties := cssStyle.Get("cssProperties").Array()
//				var index = len(cssProperties)
//				for j,_:=range cssProperties{
//					if cssProperties[j].Get("range").Exists() &&
//						(cssProperties[j].Get("range.startLine").Int()> int64(disabled[i].CssRange.StartLine) ||
//							cssProperties[j].Get("range.startLine").Int()== int64(disabled[i].CssRange.StartLine )||
//							cssProperties[j].Get("range.startColumn").Int()> int64(disabled[i].CssRange.StartColumn)){
//
//						index = j
//						break
//					}
//				}
//
//				// 畜生啊
//			}
//
//		}
//	}
//}

// extractDisabledStyles todo KeyCheck
//func (p *ProtocolAdapter) extractDisabledStyles(styleText string, cssRange *entity.IRange) []entity.IDisabledStyle {
//	var startIndices []int
//	var styles []entity.IDisabledStyle
//	for index,_:= range styleText{
//		endIndexBEGINCOMMENT := index+len(BEGIN_COMMENT)
//		endIndexENDCOMMENT := index + len(END_COMMENT)
//		if endIndexBEGINCOMMENT <=len(styleText) && string([]rune(styleText)[index:endIndexBEGINCOMMENT]) == BEGIN_COMMENT {
//			startIndices = append(startIndices,index)
//			index = index + len(BEGIN_COMMENT)
//		}else if endIndexENDCOMMENT <=len(styleText) && string([]rune(styleText)[index:endIndexENDCOMMENT]) == END_COMMENT  {
//			if len(startIndices) ==0 {
//				return nil
//			}
//			startIndex := startIndices[0]
//			startIndices = startIndices[1:]
//			endIndex := index+len(END_COMMENT)
//
//			startRangeLine,startRangeColumn := p.getLineColumnFromIndex(styleText,startIndex,cssRange)
//			endRangeLine,endRangeColumn := p.getLineColumnFromIndex(styleText,endIndex,cssRange)
//
//			propertyItem := entity.IDisabledStyle{
//				Content : styleText[startIndex:endIndex],
//				CssRange: entity.IRange{
//					StartLine :startRangeLine,
//					StartColumn: startRangeColumn,
//					EndLine: endRangeLine,
//					EndColumn: endRangeColumn,
//				},
//			}
//			styles = append(styles,propertyItem)
//			index = endIndex-1
//		}
//	}
//	if len(startIndices) ==0 {
//		return nil
//	}
//	return styles
//}

// todo KeyCheck
//func (p *ProtocolAdapter) getLineColumnFromIndex(text string,index int,startRange *entity.IRange) (line int,column int) {
//	if text=="" ||index<0||index>len(text){
//		return -1,-1
//	}
//	if startRange!=nil {
//		line = startRange.StartLine
//		column = startRange.StartColumn
//	}
//	for i := 0;i<len(text)&&i<index;i++{
//		if text[i] == '\r' && i + 1 < len(text) && text[i + 1] == '\n' {
//			i++
//			line++
//			column = 0
//		} else if text[i] == '\n' || text[i] == '\r' {
//			line++
//			column = 0
//		} else {
//			column++
//		}
//	}
//	return line,column
//}

//func compareRanges(rangeLeft gjson.Result, rangeRight gjson.Result) bool {
//	return rangeLeft.Get("startLine").Int() == rangeRight.Get("startLine").Int() &&
//		rangeLeft.Get("startColumn").Int() == rangeRight.Get("startColumn").Int() &&
//		rangeLeft.Get("endLine").Int() == rangeRight.Get("endLine").Int() &&
//		rangeLeft.Get("v").Int() == rangeRight.Get("endColumn").Int()
//}

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

//var BEGIN_COMMENT = "/* "
//var END_COMMENT = " */"
//var SEPARATOR = ": "
