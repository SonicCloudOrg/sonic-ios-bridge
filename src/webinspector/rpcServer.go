package webinspector

import (
	"fmt"
	giDevice "github.com/SonicCloudOrg/sonic-gidevice"
	"github.com/SonicCloudOrg/sonic-ios-bridge/src/entity"
	"howett.net/plist"
	"log"
	"strconv"
	"strings"
)

type RPCService struct {
	inspector            giDevice.WebInspector
	state                entity.AutomationAvailabilityType
	ConnectedApplication map[string]*entity.WebInspectorApplication
	ApplicationPages     map[string]map[string]*entity.WebInspectorPage
	WirEvent             chan []byte
}

func NewRPCServer(inspector giDevice.WebInspector) *RPCService {
	var rpc = &RPCService{
		inspector: inspector,
	}
	rpc.ConnectedApplication = make(map[string]*entity.WebInspectorApplication)
	rpc.ApplicationPages = make(map[string]map[string]*entity.WebInspectorPage)
	rpc.WirEvent = make(chan []byte)
	return rpc
}

func (r *RPCService) rpcSendMessage(selector entity.WebInspectorSelectorEnum, args entity.WIRArgument) error {
	err := r.inspector.SendWebkitMsg(string(selector), args)
	if err != nil {
		return err
	}
	return nil
}

func (r *RPCService) SendReportIdentifier(connectionID *string) error {
	if connectionID == nil {
		return fmt.Errorf("SendReportIdentifier func params connectionID is null")
	}
	var selector = entity.SEND_REPORT_ID
	var argument = entity.WIRArgument{
		WIRConnectionIdentifierKey: connectionID,
	}
	return r.rpcSendMessage(selector, argument)
}

func (r *RPCService) SendGetConnectedApplications(connectionID *string) error {
	if connectionID == nil {
		return fmt.Errorf("SendGetConnectedApplications func params connectionID is null")
	}
	var selector = entity.SEND_GET_CONNECT_APP
	var argument = entity.WIRArgument{
		WIRConnectionIdentifierKey: connectionID,
	}
	return r.rpcSendMessage(selector, argument)
}

func (r *RPCService) SendForwardGetListing(connectionID *string, appID *string) error {
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

func (r *RPCService) SendForwardIndicateWebView(connectionID *string, appID *string, pageID int, isEnabled bool) error {
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

func (r *RPCService) SendForwardSocketSetup(connectionID *string, appID *string, pageID int, senderID *string, pause bool) error {
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

func (r *RPCService) SendForwardSocketData(connectionID *string, appID *string, pageID int, senderID *string, data []byte) error {
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

func (r *RPCService) SendForwardDidClose(connectionID *string, appID *string, pageID int, senderID *string) error {
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

func (r *RPCService) parseDataToWIRMessageStruct(plistRaw interface{}) (*entity.WIRMessageStruct, error) {
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

func (r *RPCService) ReceiveAndProcess() error {
	plistRaw, err := r.inspector.ReceiveWebkitMsg()
	if err != nil {
		return err
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
		return r.ReceiveApplicationSentData(wirMessageStruct.Argument)
	case entity.ON_APP_DISCONNECTED:
		return r.ReceiveApplicationDisconnected(wirMessageStruct.Argument)
	case entity.ON_REPORT_SETUP:
		return nil
	}
	return fmt.Errorf("not the selector:" + string(wirMessageStruct.Selector))
}

// ON_REPORT_CURRENT_STATE
func (r *RPCService) ReceiveReportCurrentState(arg entity.WIRArgument) (entity.AutomationAvailabilityType, error) {
	if arg.WIRIsApplicationReadyKey == nil {
		return "", fmt.Errorf("selector:%s argumentKey: %s is nil", entity.ON_REPORT_CURRENT_STATE, "WIRIsApplicationReadyKey")
	}
	return arg.WIRAutomationAvailabilityKey, nil
}

// ON_REPORT_CONNECTED_APP_LIST
func (r *RPCService) ReceiveReportConnectedApplicationList(arg entity.WIRArgument) error {
	if arg.WIRApplicationDictionaryKey == nil {
		return fmt.Errorf("selector:%s argumentKey: %s is nil", entity.ON_REPORT_CONNECTED_APP_LIST, "WIRApplicationDictionaryKey")
	}
	for key, applicationInfo := range arg.WIRApplicationDictionaryKey {
		if app, err1 := r.parseApp(applicationInfo); err1 != nil {
			log.Println(err1)
			continue
		} else {
			r.ConnectedApplication[key] = app
		}
	}
	return nil
}

// ON_APP_SENT_LISTING
func (r *RPCService) ReceiveApplicationSentListing(arg entity.WIRArgument) error {
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
		r.ApplicationPages[*appid] = pages
	}
	return nil
}

// ON_APP_CONNECTED
func (r *RPCService) ReceiveApplicationConnected(arg entity.WIRArgument) error {
	// todo check
	if appPage, err1 := r.parseApp(arg); err1 != nil {
		return err1
	} else {
		r.ConnectedApplication[*appPage.ApplicationID] = appPage
		return nil
	}
}

// ON_APP_SENT_DATA
func (r *RPCService) ReceiveApplicationSentData(arg entity.WIRArgument) error {
	var data = arg.WIRMessageDataKey
	if data == nil {
		return fmt.Errorf("selector:%s argumentKey: %s is nil", entity.ON_APP_SENT_DATA, "WIRMessageDataKey")
	}
	if r.WirEvent != nil {
		r.WirEvent <- data
	}
	return nil
}

// ON_APP_UPDATED
func (r *RPCService) ReceiveApplicationUpdated(arg entity.WIRArgument) error {
	if appPage, err1 := r.parseApp(arg); err1 != nil {
		return err1
	} else {
		r.ConnectedApplication[*appPage.ApplicationID] = appPage
		return nil
	}
}

// ON_APP_DISCONNECTED
func (r *RPCService) ReceiveApplicationDisconnected(arg entity.WIRArgument) error {
	// todo
	return nil
}

func (r *RPCService) parseApp(args entity.WIRArgument) (appPage *entity.WebInspectorApplication, err error) {
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
