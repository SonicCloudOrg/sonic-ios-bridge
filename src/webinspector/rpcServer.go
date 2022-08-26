package webinspector

import (
	"encoding/json"
	"fmt"
	"github.com/SonicCloudOrg/sonic-ios-bridge/src/entity"
	giDevice "github.com/electricbubble/gidevice"
	"github.com/valyala/fastjson"
	"howett.net/plist"
	"log"
	"strconv"
	"strings"
)

var isDebug = false

type RPCServer struct {
	inspector            giDevice.WebInspector
	state                entity.AutomationAvailabilityType
	connectedApplication map[string]*entity.WebInspectorApplication
	applicationPages     map[string]map[string]*entity.WebInspectorPage
}

func SetRPCDebug(flag bool) {
	isDebug = flag
}

func (r *RPCServer) rpcSendMessage(selector entity.WebInspectorSelectorEnum, args entity.WIRArgument) error {
	if isDebug {
		log.Println("----->")
		log.Println(fmt.Sprintf("selector:%s", selector))
		log.Println(args)
		fmt.Println()
	}
	err := r.inspector.SendWebkitMsg(string(selector), args)
	if err != nil {
		return err
	}
	return nil
}

func (r *RPCServer) SendReportIdentifier(connectionID *string) error {
	if connectionID == nil {
		return fmt.Errorf("SendReportIdentifier func params connectionID is null")
	}
	var selector = entity.SEND_REPORT_ID
	var argument = entity.WIRArgument{
		WIRConnectionIdentifierKey: connectionID,
	}
	return r.rpcSendMessage(selector, argument)
}

func (r *RPCServer) SendGetConnectedApplications(connectionID *string) error {
	if connectionID == nil {
		return fmt.Errorf("SendGetConnectedApplications func params connectionID is null")
	}
	var selector = entity.SEND_GET_CONNECT_APP
	var argument = entity.WIRArgument{
		WIRConnectionIdentifierKey: connectionID,
	}
	return r.rpcSendMessage(selector, argument)
}

func (r *RPCServer) SendForwardGetListing(connectionID *string, appID *string) error {
	if connectionID == nil || appID == nil {
		return fmt.Errorf("SendForwardGetListing func params is null")
	}
	var selector = entity.SEND_FORWARD_GET_LISTING
	var argument = entity.WIRArgument{
		WIRConnectionIdentifierKey:  connectionID,
		WIRApplicationIdentifierKey: appID,
	}
	return r.rpcSendMessage(selector, argument)
}

func (r *RPCServer) SendForwardIndicateWebView(connectionID *string, appID *string, pageID int, isEnabled bool) error {
	if connectionID == nil || appID == nil {
		return fmt.Errorf("SendForwardIndicateWebView func params is null")
	}
	var selector = entity.SEND_FORWARD_INDICATE_WEBVIEW
	var argument = entity.WIRArgument{
		WIRConnectionIdentifierKey:  connectionID,
		WIRApplicationIdentifierKey: appID,
		WIRPageIdentifierKey:        &pageID,
		WIRIndicateEnabledKey:       &isEnabled,
	}
	return r.rpcSendMessage(selector, argument)
}

func (r *RPCServer) SendForwardSocketSetup(connectionID *string, appID *string, pageID int, senderID *string, pause bool) error {
	if connectionID == nil || appID == nil || senderID == nil {
		return fmt.Errorf("SendForwardSocketSetup func params is null")
	}
	var selector = entity.SEND_FORWARD_SOCKET_SETUP
	var argument = entity.WIRArgument{
		WIRConnectionIdentifierKey:  connectionID,
		WIRApplicationIdentifierKey: appID,
		WIRPageIdentifierKey:        &pageID,
		WIRSenderKey:                senderID,
	}
	if !pause {
		argument.WIRAutomaticallyPause = &pause
	}
	return r.rpcSendMessage(selector, argument)
}

func (r *RPCServer) SendForwardSocketData(connectionID *string, appID *string, pageID int, senderID *string, data *string) error {
	if connectionID == nil || appID == nil || senderID == nil || data == nil {
		return fmt.Errorf("SendForwardSocketData func params is null")
	}
	var selector = entity.SEND_FORWARD_SOCKET_DATA
	var argument = entity.WIRArgument{
		WIRConnectionIdentifierKey:  connectionID,
		WIRApplicationIdentifierKey: appID,
		WIRPageIdentifierKey:        &pageID,
		WIRSenderKey:                senderID,
		WIRSocketDataKey:            data,
	}
	return r.rpcSendMessage(selector, argument)
}

func (r *RPCServer) SendForwardDidClose(connectionID *string, appID *string, pageID int, senderID *string) error {
	if connectionID == nil || appID == nil || senderID == nil {
		return fmt.Errorf("SendForwardDidClose func params is null")
	}
	var selector = entity.SEND_FORWARD_DID_CLOSE
	var argument = entity.WIRArgument{
		WIRConnectionIdentifierKey:  connectionID,
		WIRApplicationIdentifierKey: appID,
		WIRPageIdentifierKey:        &pageID,
		WIRSenderKey:                senderID,
	}
	return r.rpcSendMessage(selector, argument)
}

func (r *RPCServer) parseDataToWIRMessageStruct(plistRaw interface{}) (*entity.WIRMessageStruct, error) {
	arr, err := plist.Marshal(plistRaw, plist.BinaryFormat)
	if err != nil {
		return nil, err
	}
	var paresPlist = entity.WIRMessageStruct{}
	if _, err = plist.Unmarshal(arr, &paresPlist); err != nil {
		return nil, err
	}
	return &paresPlist, nil
}

// todo isDebu print error
func (r *RPCServer) ReceiveAndProcess() error {
	plistRaw, err := r.inspector.ReceiveWebkitMsg()
	if err != nil {
		return err
	}
	if isDebug {
		log.Println("<-----")
		log.Print("recv data:")
		log.Println(plistRaw)
		fmt.Println()
	}
	wirMessageStruct, err := r.parseDataToWIRMessageStruct(plistRaw)
	if err != nil {
		return err
	}
	switch wirMessageStruct.Selector {
	case entity.ON_REPORT_CURRENT_STATE:
		r.state = wirMessageStruct.Argument.WIRAutomationAvailabilityKey
		return nil
	case entity.ON_REPORT_CONNECTED_APP_LIST:
		return r.ReceiveReportConnectedApplicationList(wirMessageStruct.Argument)
	case entity.ON_APP_SENT_LISTING:
		return r.ReceiveApplicationSentListing(wirMessageStruct.Argument)
	case entity.ON_REPORT_DRIVER_LIST:
		return nil
	case entity.ON_APP_UPDATED:
		return r.ReceiveApplicationUpdated(wirMessageStruct.Argument)
	case entity.ON_APP_CONNECTED:
		return r.ReceiveApplicationConnected(wirMessageStruct.Argument)
	case entity.ON_APP_SENT_DATA:
		// todo keyCheck
		fmt.Println("start recv data")
		wirMessages, wirEvents, err1 := r.ReceiveApplicationSentData(wirMessageStruct.Argument)
		var p fastjson.Parser
		for _, v := range wirMessages {
			if data, err2 := p.Parse(v); err2 != nil {
				return err2
			} else {
				fmt.Println(data)
			}
		}

		for _, v := range wirEvents {
			if data, err2 := p.Parse(v); err2 != nil {
				return err2
			} else {
				fmt.Println(data)
			}
		}
		fmt.Println("recv end")
		fmt.Println()
		return err1
	case entity.ON_APP_DISCONNECTED:
		return r.ReceiveApplicationDisconnected(wirMessageStruct.Argument)
	case entity.ON_REPORT_SETUP:
		return nil
	}
	return fmt.Errorf("not the selector:" + string(wirMessageStruct.Selector))
}

// ON_REPORT_CURRENT_STATE
func (r *RPCServer) ReceiveReportCurrentState(arg entity.WIRArgument) (entity.AutomationAvailabilityType, error) {
	if arg.WIRIsApplicationReadyKey == nil {
		return "", fmt.Errorf("selector:%s argumentKey: %s is nil", entity.ON_REPORT_CURRENT_STATE, "WIRIsApplicationReadyKey")
	}
	return arg.WIRAutomationAvailabilityKey, nil
}

// ON_REPORT_CONNECTED_APP_LIST
func (r *RPCServer) ReceiveReportConnectedApplicationList(arg entity.WIRArgument) error {
	if arg.WIRApplicationDictionaryKey == nil {
		return fmt.Errorf("selector:%s argumentKey: %s is nil", entity.ON_REPORT_CONNECTED_APP_LIST, "WIRApplicationDictionaryKey")
	}
	for key, applicationInfo := range arg.WIRApplicationDictionaryKey {
		if app, err1 := r.parseApp(applicationInfo); err1 != nil {
			if isDebug {
				log.Fatal(err1)
			}
			continue
		} else {
			r.connectedApplication[key] = app
		}
	}
	return nil
}

// ON_APP_SENT_LISTING
func (r *RPCServer) ReceiveApplicationSentListing(arg entity.WIRArgument) error {
	var item = arg.WIRListingKey
	if item == nil {
		return fmt.Errorf("selector:%s argumentKey: %s is nil", entity.ON_APP_SENT_LISTING, "WIRListingKey")
	}
	var appid = arg.WIRApplicationIdentifierKey
	if appid == nil {
		return fmt.Errorf("selector:%s argumentKey: %s is nil", entity.ON_APP_SENT_LISTING, "WIRApplicationIdentifierKey")
	}
	pages := make(map[string]*entity.WebInspectorPage)
	for id, page := range item {
		pages[id] = &page
	}
	if len(pages) > 0 {
		r.applicationPages[*appid] = pages
	}
	return nil
}

// ON_APP_CONNECTED
func (r *RPCServer) ReceiveApplicationConnected(arg entity.WIRArgument) error {
	// todo check
	if appPage, err1 := r.parseApp(arg); err1 != nil {
		return err1
	} else {
		r.connectedApplication[*appPage.ApplicationID] = appPage
		return nil
	}
}

// ON_APP_SENT_DATA
func (r *RPCServer) ReceiveApplicationSentData(arg entity.WIRArgument) (map[int]string, []string, error) {
	var data = arg.WIRMessageDataKey
	if data == nil {
		return nil, nil, fmt.Errorf("selector:%s argumentKey: %s is nil", entity.ON_APP_SENT_DATA, "WIRMessageDataKey")
	}
	arr, err := json.Marshal(data)
	if err != nil {
		return nil, nil, err
	}
	wirMessageResult := make(map[int]string)
	var wirEvent []string
	// todo check
	if id := data["id"]; id != nil {
		wirMessageResult[id.(int)] = string(arr)
	} else {
		wirEvent = append(wirEvent, string(arr))
	}
	return wirMessageResult, wirEvent, nil
}

// ON_APP_UPDATED
func (r *RPCServer) ReceiveApplicationUpdated(arg entity.WIRArgument) error {
	if appPage, err1 := r.parseApp(arg); err1 != nil {
		return err1
	} else {
		r.connectedApplication[*appPage.ApplicationID] = appPage
		return nil
	}
}

// ON_APP_DISCONNECTED
func (r *RPCServer) ReceiveApplicationDisconnected(arg entity.WIRArgument) error {
	// todo
	return nil
}

func (r *RPCServer) parseApp(args entity.WIRArgument) (appPage *entity.WebInspectorApplication, err error) {
	if args.WIRApplicationIdentifierKey == nil {
		return nil, fmt.Errorf("parse app is fail")
	}
	var page = &entity.WebInspectorApplication{
		ApplicationID:           args.WIRApplicationIdentifierKey,
		ApplicationBundle:       args.WIRApplicationBundleIdentifierKey,
		ApplicationPID:          keyToPID(*args.WIRApplicationIdentifierKey),
		ApplicationName:         args.WIRApplicationNameKey,
		ApplicationAvailability: args.WIRAutomationAvailabilityKey,
		ApplicationActive:       args.WIRIsApplicationActiveKey,
		ApplicationProxy:        args.WIRIsApplicationProxyKey,
		ApplicationReady:        args.WIRIsApplicationReadyKey,
		ApplicationHost:         args.WIRHostApplicationIdentifierKey,
	}
	return page, nil
}

func keyToPID(key string) *int {
	var pidStr = strings.Split(key, ":")[1]
	pid, err := strconv.Atoi(pidStr)
	if err != nil {
		pid = -1
		fmt.Println(err)
		return &pid
	}
	return &pid
}
