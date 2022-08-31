package webinspector

import (
	"context"
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
	senderID             string
	ctx                  context.Context
	conn                 *websocket.Conn
	closeSendWS          context.Context
}

var isProtocolDebug = false

func SetProtocolDebug(flag bool) {
	isProtocolDebug = flag
}

func NewWebkitDebugService(device *giDevice.Device, ctx context.Context) *WebkitDebugService {
	return &WebkitDebugService{
		device:    device,
		connectID: strings.ToUpper(uuid.New().String()),
		ctx:       ctx,
	}
}

func (w *WebkitDebugService) ConnectInspector() (context.CancelFunc, error) {
	if w.device == nil {
		return nil, fmt.Errorf("device is null")
	}
	webInspectorService, err := (*w.device).WebInspectorService()
	if err != nil {
		return nil, err
	}

	// init
	w.inspector = &webInspectorService
	w.rpcService = NewRPCServer(*w.inspector)
	w.applicationPages = w.rpcService.ApplicationPages
	w.connectedApplication = w.rpcService.ConnectedApplication

	if len(w.rpcService.ApplicationPages) == 0 {
		err = w.rpcService.SendReportIdentifier(&w.connectID)
		if err != nil {
			return nil, err
		}
	}

	ctx, cancel := context.WithCancel(w.ctx)

	go func() {
		for {
			select {
			case <-ctx.Done():
				w.Close()
			default:
				err = w.rpcService.ReceiveAndProcess()
				if err != nil {
					if strings.Contains(err.Error(), "timeout") {
						continue
					}
					fmt.Println(err)
					return
				}
			}
		}
	}()

	return cancel, err
}

func (w *WebkitDebugService) Close() {
	if w.rpcService.WirEvent != nil {
		close(w.rpcService.WirEvent)
		w.rpcService.WirEvent = nil
	}
}

func (w *WebkitDebugService) StartCDP(appID *string, pageID *int, conn *websocket.Conn) error {
	senderID := strings.ToUpper(uuid.New().String())
	w.senderID = senderID
	w.conn = conn
	var closeSendProtocol context.CancelFunc
	w.closeSendWS, closeSendProtocol = context.WithCancel(w.ctx)

	w.conn.SetCloseHandler(func(code int, text string) error {
		log.Println("try close ws")
		conn = nil
		closeSendProtocol()
		return w.rpcService.SendForwardDidClose(&w.connectID, appID, *pageID, &senderID)
	})
	//w.connectID = strings.ToUpper(uuid.New().String())
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

func (w *WebkitDebugService) SendProtocolCommand(applicationID *string, pageID *int, message []byte) {
	if isProtocolDebug {
		log.Println(fmt.Sprintf("protocol send command:%s\n", string(message)))
	}
	err := w.rpcService.SendForwardSocketData(&w.connectID, applicationID, *pageID, &w.senderID, message)
	if err != nil {
		log.Fatal(err)
	}
}

func (w *WebkitDebugService) ReceiveProtocolData(conn *websocket.Conn) error {
	if w.rpcService.WirEvent != nil {
		select {
		case message, ok := <-w.rpcService.WirEvent:
			if ok {
				if isProtocolDebug {
					log.Println(fmt.Sprintf("protocol receive command:%s\n", string(message)))
				}
				if conn != nil {
					err := conn.WriteMessage(websocket.TextMessage, message)
					if err != nil {
						log.Println("ws send error:", err)
						return err
					}
				}
			}
		case <-w.closeSendWS.Done():
			return fmt.Errorf("close send protocol")
		}
	}
	return nil
}
