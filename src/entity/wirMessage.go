/*
 *   sonic-ios-bridge  Connect to your iOS Devices.
 *   Copyright (C) 2022 SonicCloudOrg
 *
 *   This program is free software: you can redistribute it and/or modify
 *   it under the terms of the GNU Affero General Public License as published
 *   by the Free Software Foundation, either version 3 of the License, or
 *   (at your option) any later version.
 *
 *   This program is distributed in the hope that it will be useful,
 *   but WITHOUT ANY WARRANTY; without even the implied warranty of
 *   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 *   GNU Affero General Public License for more details.
 *
 *   You should have received a copy of the GNU Affero General Public License
 *   along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */
package entity

type (
	WIRArgument struct {
		WIRMessageDataKey          []byte                     `plist:"WIRMessageDataKey,omitempty"`
		WIRConnectionIdentifierKey *string                    `plist:"WIRConnectionIdentifierKey,omitempty"`
		WIRPageIdentifierKey       *int                       `plist:"WIRPageIdentifierKey,omitempty"`
		WIRIndicateEnabledKey      *bool                      `plist:"WIRIndicateEnabledKey,omitempty"`
		WIRSessionIdentifierKey    *string                    `plist:"WIRSessionIdentifierKey,omitempty"`
		WIRSenderKey               *string                    `plist:"WIRSenderKey,omitempty"`
		WIRAutomaticallyPause      *bool                      `plist:"WIRAutomaticallyPause,omitempty"`
		WIRSocketDataKey           []byte                     `plist:"WIRSocketDataKey,omitempty"`
		WIRSessionCapabilitiesKey  *WIRSessionCapabilitiesKey `plist:"WIRSessionCapabilitiesKey,omitempty"`
		// 单个appInfo情况下
		WIRApplicationIdentifierKey       *string                    `plist:"WIRApplicationIdentifierKey,omitempty"`
		WIRApplicationBundleIdentifierKey *string                    `plist:"WIRApplicationBundleIdentifierKey,omitempty"`
		WIRApplicationNameKey             *string                    `plist:"WIRApplicationNameKey,omitempty"`
		WIRAutomationAvailabilityKey      AutomationAvailabilityType `plist:"WIRAutomationAvailabilityKey,omitempty"`
		WIRIsApplicationActiveKey         *int                       `plist:"WIRIsApplicationActiveKey,omitempty"`
		WIRIsApplicationProxyKey          *bool                      `plist:"WIRIsApplicationProxyKey,omitempty"`
		WIRIsApplicationReadyKey          *bool                      `plist:"WIRIsApplicationReadyKey,omitempty"`
		WIRHostApplicationIdentifierKey   *string                    `plist:"WIRHostApplicationIdentifierKey,omitempty"`
		// 多个appInfo情况下
		WIRApplicationDictionaryKey map[string]WIRArgument      `plist:"WIRApplicationDictionaryKey,omitempty"`
		WIRListingKey               map[string]WebInspectorPage `plist:"WIRListingKey,omitempty"`
	}

	WIRSessionCapabilitiesKey struct {
		AllowInsecureMediaCapture    bool `plist:"org.webkit.webdriver.webrtc.allow-insecure-media-capture"`
		SppressIceCandidateFiltering bool `plist:"org.webkit.webdriver.webrtc.suppress-ice-candidate-filtering"`
	}

	WIRMessageStruct struct {
		Argument WIRArgument              `plist:"__argument"`
		Selector WebInspectorSelectorEnum `plist:"__selector"`
	}
)

type WebInspectorSelectorEnum string

const (
	// 发送消息类型
	SEND_REPORT_ID                WebInspectorSelectorEnum = "_rpc_reportIdentifier:"
	SEND_GET_CONNECT_APP          WebInspectorSelectorEnum = "_rpc_getConnectedApplications:"
	SEND_FORWARD_GET_LISTING      WebInspectorSelectorEnum = "_rpc_forwardGetListing:"
	SEND_FORWARD_SOCKET_SETUP     WebInspectorSelectorEnum = "_rpc_forwardSocketSetup:"
	SEND_FORWARD_SOCKET_DATA      WebInspectorSelectorEnum = "_rpc_forwardSocketData:"
	SEND_FORWARD_INDICATE_WEBVIEW WebInspectorSelectorEnum = "_rpc_forwardIndicateWebView:"
	SEND_FORWARD_DID_CLOSE        WebInspectorSelectorEnum = "_rpc_forwardDidClose:"

	REQUEST_APPLICATION_LAUNCH WebInspectorSelectorEnum = "_rpc_requestApplicationLaunch"

	// 接收消息类型
	ON_REPORT_CURRENT_STATE      WebInspectorSelectorEnum = "_rpc_reportCurrentState:"
	ON_REPORT_SETUP              WebInspectorSelectorEnum = "_rpc_reportSetup:"
	ON_REPORT_DRIVER_LIST        WebInspectorSelectorEnum = "_rpc_reportConnectedDriverList:"
	ON_REPORT_CONNECTED_APP_LIST WebInspectorSelectorEnum = "_rpc_reportConnectedApplicationList:"
	ON_APP_CONNECTED             WebInspectorSelectorEnum = "_rpc_applicationConnected:"
	ON_APP_UPDATED               WebInspectorSelectorEnum = "_rpc_applicationUpdated:"
	ON_APP_SENT_LISTING          WebInspectorSelectorEnum = "_rpc_applicationSentListing:"
	ON_APP_SENT_DATA             WebInspectorSelectorEnum = "_rpc_applicationSentData:"
	ON_APP_DISCONNECTED          WebInspectorSelectorEnum = "_rpc_applicationDisconnected:"
)

type (
	WebInspectorPage struct {
		PageID                     *int                 `plist:"WIRPageIdentifierKey,omitempty"`
		PageType                   WebInspectorPageType `plist:"WIRTypeKey,omitempty"`
		PageWebUrl                 *string              `plist:"WIRURLKey,omitempty"`
		PageWebTitle               *string              `plist:"WIRTitleKey,omitempty"`
		PageAutoationIsPairedKey   *bool                `plist:"WIRAutomationTargetIsPairedKey,omitempty"`
		PageAutomationName         *string              `plist:"WIRAutomationTargetNameKey,omitempty"`
		PageAutomationVersion      *string              `plist:"WIRAutomationTargetVersionKey,omitempty"`
		PageAutomationSessionID    *string              `plist:"WIRSessionIdentifierKey,omitempty"`
		PageAutomationConnectionID *string              `plist:"WIRConnectionIdentifierKey,omitempty"`
	}

	WebInspectorApplication struct {
		ApplicationID           *string `plist:"WIRApplicationIdentifierKey,omitempty"`
		ApplicationBundle       *string `plist:"WIRApplicationBundleIdentifierKey,omitempty"`
		ApplicationPID          *int
		ApplicationName         *string                    `plist:"WIRApplicationNameKey,omitempty"`
		ApplicationAvailability AutomationAvailabilityType `plist:"WIRAutomationAvailabilityKey,omitempty"`
		ApplicationActive       *int                       `plist:"WIRIsApplicationActiveKey,omitempty"`
		ApplicationProxy        *bool                      `plist:"WIRIsApplicationProxyKey,omitempty"`
		ApplicationReady        *bool                      `plist:"WIRIsApplicationReadyKey,omitempty"`
		ApplicationHost         *string                    `plist:"WIRHostApplicationIdentifierKey,omitempty"`
	}
)

type WebInspectorPageType string
type AutomationAvailabilityType string

const (
	AUTOMATION          WebInspectorPageType = "WIRTypeAutomation"
	ITML                WebInspectorPageType = "WIRTypeITML"
	JAVASCRIPT          WebInspectorPageType = "WIRTypeJavaScript"
	PAGE                WebInspectorPageType = "WIRTypePage"
	SERVICE_WORKER      WebInspectorPageType = "WIRTypeServiceWorker"
	WEB                 WebInspectorPageType = "WIRTypeWeb"
	WEB_PAGE            WebInspectorPageType = "WIRTypeWebPage"
	AUTOMATICALLY_PAUSE WebInspectorPageType = "WIRAutomaticallyPause"

	NOT_AVAILABLE       AutomationAvailabilityType = "WIRAutomationAvailabilityNotAvailable"
	AVAILABLE           AutomationAvailabilityType = "WIRAutomationAvailabilityAvailable"
	AvailabilityUNKNOWN AutomationAvailabilityType = "WIRAutomationAvailabilityUnknown"
)

type UrlItem struct {
	Description          string  `json:"description"`
	ID                   string  `json:"id"`
	Port                 int     `json:"port"`
	Title                *string `json:"title"`
	Type                 string  `json:"type"`
	Url                  *string `json:"url"`
	WebSocketDebuggerUrl string  `json:"webSocketDebuggerUrl"`
	DevtoolsFrontendUrl  string  `json:"devtoolsFrontendUrl"`
}

type BundleItem struct {
	PID      string    `json:"pid,omitempty"`
	BundleId string    `json:"bundleId,omitempty"`
	Name     string    `json:"name,omitempty"`
	Pages    []UrlItem `json:"pages,omitempty"`
}
