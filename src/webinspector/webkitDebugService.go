package webinspector

import (
	"fmt"
	"github.com/SonicCloudOrg/sonic-ios-bridge/src/entity"
	giDevice "github.com/electricbubble/gidevice"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"log"
	"strings"
	"sync"
)

type WebkitDebugService struct {
	connectID            string
	rpcService           *RPCService
	inspector            *giDevice.WebInspector
	device               *giDevice.Device
	connectedApplication map[string]*entity.WebInspectorApplication
	applicationPages     map[string]map[string]*entity.WebInspectorPage
	// 摘出？
	senderID string
}

func NewWebkitDebugService(device *giDevice.Device) *WebkitDebugService {
	return &WebkitDebugService{
		device:    device,
		connectID: strings.ToUpper(uuid.New().String()),
		senderID:  strings.ToUpper(uuid.New().String()),
	}
}

func (w *WebkitDebugService) ConnectInspector() error {
	if w.device == nil {
		return fmt.Errorf("device is null")
	}
	webInspectorService, err := (*w.device).WebInspectorService()
	if err != nil {
		return err
	}

	// init
	w.inspector = &webInspectorService
	w.rpcService = NewRPCServer(*w.inspector)
	w.applicationPages = w.rpcService.ApplicationPages
	w.connectedApplication = w.rpcService.ConnectedApplication

	if len(w.rpcService.ApplicationPages) == 0 {
		err = w.rpcService.SendReportIdentifier(&w.connectID)
		if err != nil {
			return err
		}
	}

	go func() {
		for {
			err = w.rpcService.ReceiveAndProcess()
			if err != nil {
				if strings.Contains(err.Error(), "timeout") {
					continue
				}
				return
			}
		}
	}()
	return err
}

func (w *WebkitDebugService) StartCDP(appID *string, pageID *int) error {
	return w.rpcService.SendForwardSocketSetup(&w.connectID, appID, *pageID, &w.senderID, false)
}

func (w *WebkitDebugService) FindPagesByID(pageId string) (application *entity.WebInspectorApplication, page *entity.WebInspectorPage, err error) {
	for appID, value := range w.applicationPages {
		for id := range value {
			if id == pageId {
				application = w.connectedApplication[appID]
				page = w.applicationPages[appID][id]
				return
			}
		}
	}
	return nil, nil, fmt.Errorf("not find page")
}

func (w *WebkitDebugService) GetOpenPages(port int) ([]entity.UrlItem, error) {
	var wg = sync.WaitGroup{}
	for key, _ := range w.connectedApplication {
		wg.Add(1)
		go func(key string) {
			err := w.rpcService.SendForwardGetListing(&w.connectID, &key)
			if err != nil {
				log.Fatal(err)
			}
			wg.Done()
		}(key)
	}
	wg.Wait()
	var pages []entity.UrlItem
	for appID, _ := range w.applicationPages {
		for pageID, page := range w.applicationPages[appID] {
			//if page.PageType != entity.WEB && page.PageType != entity.WEB_PAGE {
			//	continue
			//}
			var pageItem = &entity.UrlItem{
				Description:          "",
				ID:                   pageID,
				Title:                page.PageWebTitle,
				Type:                 "page",
				Url:                  page.PageWebUrl,
				WebSocketDebuggerUrl: fmt.Sprintf("ws://localhost:%d/devtools/page/%s", port, pageID),
				DevtoolsFrontendUrl:  fmt.Sprintf("/devtools/inspector.html?ws://localhost:%d/devtools/page/%s", port, pageID),
			}
			pages = append(pages, *pageItem)
		}
	}
	return pages, nil
}

var isProtocolDebug = false

func SetProtocolDebug(flag bool) {
	isProtocolDebug = flag
}

func (w *WebkitDebugService) SendProtocolCommand(applicationID *string, pageID *int, message []byte) {
	if isProtocolDebug {
		log.Println(fmt.Sprintf("protocol send command:%s", string(message)))
		fmt.Println()
	}
	err := w.rpcService.SendForwardSocketData(&w.connectID, applicationID, *pageID, &w.senderID, message)
	if err != nil {
		log.Fatal(err)
	}
}

func (w *WebkitDebugService) ReceiveProtocolData(conn *websocket.Conn) {
	select {
	case message, ok := <-w.rpcService.WirEvent:
		if ok {
			if isProtocolDebug {
				log.Println("protocol receive command:")
				log.Println(string(message))
				fmt.Println()
			}
			err := conn.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}
