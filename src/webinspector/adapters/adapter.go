package adapters

import (
	"encoding/json"
	"github.com/SonicCloudOrg/sonic-ios-bridge/src/entity"
	"github.com/SonicCloudOrg/sonic-ios-bridge/src/webinspector"
	"log"
	"strings"
)

type Adapter struct {
	targetID          string
	messageFilters    map[string]MessageAdapters
	webDebugService   *webinspector.WebkitDebugService
	isTargetBased     bool
	applicationID     *string
	pageID            *int
	waitingForID      int
	adapterRequestMap map[int]func(message []byte)
}

func (a *Adapter) AddMessageFilter(method string, filter MessageAdapters) {
	if a.messageFilters == nil {
		a.messageFilters = make(map[string]MessageAdapters)
	}
	a.messageFilters[method] = filter
}

func (a *Adapter) CallTarget(method string, params interface{}, callFunc func(message []byte)) {
	a.waitingForID -= 1
	var message = &entity.TargetProtocol{}
	arr, err := json.Marshal(params)
	if err != nil {
		log.Fatal(err)
	}
	println(string(arr))
	message.ID = a.waitingForID
	message.Method = method
	message.Params = params
	a.adapterRequestMap[a.waitingForID] = callFunc
	a.sendToTarget(message)
}

func (a *Adapter) sendToTarget(message *entity.TargetProtocol) {
	log.Println("origin send message:")
	arr, err := json.Marshal(message)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(string(arr))
	if a.isTargetBased {
		if !strings.Contains(message.Method, "Target") {
			var newMessage = &entity.TargetProtocol{}

			newMessage.ID = message.ID
			newMessage.Method = "Target.sendMessageToTarget"
			newMessage.Params = &entity.TargetParams{
				TargetId: a.targetID,
				Message:  string(arr),
				ID:       message.ID,
			}
			message = newMessage
			log.Println("new send message:")
		}
	}
	arr, err = json.Marshal(message)
	if err != nil {
		log.Fatal(err)
	}
	a.webDebugService.SendProtocolCommand(a.applicationID, a.pageID, arr)
}

func (a *Adapter) FireResultToTools(id int, params interface{}) {
	response := map[string]interface{}{
		"id":     id,
		"result": params,
	}
	arr, err := json.Marshal(response)
	if err != nil {
		log.Panic(err)
	}
	a.webDebugService.SendMessageTool(arr)
}
