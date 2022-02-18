package conn

var generationMap = map[string]string{
	"AirPods1,1":    "AirPods (1st generation)",
	"AirPods1,2":    "AirPods (2nd generation)",
	"AirPods2,1":    "AirPods (2nd generation)",
	"AirPods1,3":    "AirPods (3rd generation)",
	"Audio2,1":      "AirPods (3rd generation)",
	"AirPods2,2":    "AirPods Pro",
	"AirPodsPro1,1": "AirPods Pro",
	"iProd8,1":      "AirPods Pro",
	"AirTag1,1":     "AirTag",
	"AppleTV1,1":    "Apple TV (1st generation)",
	"AppleTV2,1":    "Apple TV (2nd generation)",
	"AppleTV3,1":    "Apple TV (3rd generation)",
	"AppleTV3,2":    "Apple TV (3rd generation)",
	"AppleTV5,3":    "Apple TV (4th generation)",
	"AppleTV6,2":    "Apple TV 4K",
	"AppleTV11,1":   "Apple TV 4K (2nd generation)",
}

func (deviceDetail *DeviceDetail) GetGenerationName() string {
	if len(deviceDetail.ProductType) > 0 {
		return generationMap[deviceDetail.ProductType]
	} else {
		return ""
	}
}
