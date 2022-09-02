package adapters

import "github.com/valyala/fastjson"

type MessageAdapters func(message *fastjson.Value) *fastjson.Value

var PageSetOverlay MessageAdapters = func(message *fastjson.Value) *fastjson.Value {
	method := "Debugger.setOverlayMessage"
	return setMethod(message, method)
}

var PageConfigureOverlay MessageAdapters = func(message *fastjson.Value) *fastjson.Value {
	return PageSetOverlay(message)
}

var DOMSetInspectedNode MessageAdapters = func(message *fastjson.Value) *fastjson.Value {
	method := "Console.addInspectedNode"
	return setMethod(message, method)
}

var EmulationSetTouchEmulationEnabled MessageAdapters = func(message *fastjson.Value) *fastjson.Value {
	method := "Page.setTouchEmulationEnabled"
	return setMethod(message, method)
}

var EmulationSetScriptExecutionDisabled MessageAdapters = func(message *fastjson.Value) *fastjson.Value {
	method := "Page.setScriptExecutionDisabled"
	return setMethod(message, method)
}

var EmulationSetEmulatedMedia MessageAdapters = func(message *fastjson.Value) *fastjson.Value {
	method := "Page.setEmulatedMedia"
	return setMethod(message, method)
}

var RenderingSetShowPaintRects MessageAdapters = func(message *fastjson.Value) *fastjson.Value {
	method := "Page.setShowPaintRects"
	return setMethod(message, method)
}

var LogClear MessageAdapters = func(message *fastjson.Value) *fastjson.Value {
	method := "Console.clearMessages"
	return setMethod(message, method)
}

var LogDisable MessageAdapters = func(message *fastjson.Value) *fastjson.Value {
	method := "Console.disable"
	return setMethod(message, method)
}

var LogEnable MessageAdapters = func(message *fastjson.Value) *fastjson.Value {
	method := "Console.enable"
	return setMethod(message, method)
}
var NetworkGetCookies MessageAdapters = func(message *fastjson.Value) *fastjson.Value {
	method := "Page.getCookies"
	return setMethod(message, method)
}

var NetworkDeleteCookie MessageAdapters = func(message *fastjson.Value) *fastjson.Value {
	method := "Page.deleteCookie"
	return setMethod(message, method)
}

var NetworkSetMonitoringXHREnabled MessageAdapters = func(message *fastjson.Value) *fastjson.Value {
	method := "Console.setMonitoringXHREnabled"
	return setMethod(message, method)
}

//var CSSSetStyleTexts MessageAdapters = func(message *fastjson.Value) *fastjson.Value {
//	edits := message.Get("params").GetArray("edits")
//}

func setMethod(message *fastjson.Value, methodName string) *fastjson.Value {
	message.Set("method", fastjson.MustParse(methodName))
	return message
}
