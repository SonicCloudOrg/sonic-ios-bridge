package entity

import (
	"encoding/json"
	"fmt"
	giDevice "github.com/SonicCloudOrg/sonic-gidevice"
	"strings"
)

type PerformanceData struct {
	CPUInfo *giDevice.CPUInfo        `json:"CPUInfo,omitempty"`
	GPUInfo *giDevice.GPUInfo        `json:"GPUInfo,omitempty"`
	MEMInfo *giDevice.MEMInfo        `json:"MemoryInfo,omitempty"`
	FPSInfo *giDevice.FPSInfo        `json:"FPSInfo,omitempty"`
	Network *giDevice.NetWorkingInfo `json:"NetworkingInfo,omitempty"`
}

func CreatePerformanceData(data interface{}) PerformanceData {
	var perfmonData PerformanceData
	if bUser, ok := data.(giDevice.CPUInfo); ok {
		perfmonData.CPUInfo = &bUser
	}
	if bUser, ok := data.(giDevice.MEMInfo); ok {
		perfmonData.MEMInfo = &bUser
	}
	if bUser, ok := data.(giDevice.FPSInfo); ok {
		perfmonData.FPSInfo = &bUser
	}
	if bUser, ok := data.(giDevice.GPUInfo); ok {
		perfmonData.GPUInfo = &bUser
	}
	if bUser, ok := data.(giDevice.NetWorkingInfo); ok {
		perfmonData.Network = &bUser
	}
	return perfmonData
}

func (perfData PerformanceData) ToString() string {
	var s strings.Builder
	s.WriteString(toStringCPUInfo(perfData.CPUInfo))
	s.WriteString(toStringMEMInfo(perfData.MEMInfo))
	s.WriteString(toGPUInfoString(perfData.GPUInfo))
	s.WriteString(toStringFPSInfo(perfData.FPSInfo))
	s.WriteString(toStringNetwork(perfData.Network))
	return s.String()
}

func (perfData PerformanceData) ToJson() string {
	result, _ := json.Marshal(perfData)
	return string(result)
}

func (perfData PerformanceData) ToFormat() string {
	result, _ := json.MarshalIndent(perfData, "", "\t")
	return string(result)
}

func toStringCPUInfo(CPUInfo *giDevice.CPUInfo) string {
	var s strings.Builder
	if CPUInfo == nil {
		return ""
	}
	s.WriteString("CPUInfo:\n")
	if CPUInfo.Mess != "" {
		s.WriteString(fmt.Sprintf("SystemCpuCount: %d, SystemCpuUsage:%f, ProcessInfo:%s", CPUInfo.CPUCount,
			CPUInfo.SysCpuUsage, CPUInfo.Mess))
	} else {
		s.WriteString(
			fmt.Sprintf("PID:%s, SystemCpuCount:%d, SystemCpuUsage:%f, ProcessAttrCtxSwitch:%d, ProcessAttrIntWakeups:%d, ProcessCpuUsage:%f",
				CPUInfo.Pid, CPUInfo.CPUCount, CPUInfo.SysCpuUsage,
				CPUInfo.AttrCtxSwitch,
				CPUInfo.AttrIntWakeups,
				CPUInfo.CPUUsage))
	}
	s.WriteString(fmt.Sprintf(", TimeStamp:%d", CPUInfo.TimeStamp))
	return s.String()
}

func toGPUInfoString(GPUInfo *giDevice.GPUInfo) string {
	var s strings.Builder
	if GPUInfo == nil {
		return ""
	}
	s.WriteString(fmt.Sprintf("GPUInfo:\nTilerUtilization:%d, DeviceUtilization:%d, RendererUtilization:%d",
		GPUInfo.TilerUtilization, GPUInfo.DeviceUtilization, GPUInfo.RendererUtilization))
	s.WriteString(fmt.Sprintf(", TimeStamp:%d", GPUInfo.TimeStamp))
	return s.String()
}

func toStringMEMInfo(MEMInfo *giDevice.MEMInfo) string {
	var s strings.Builder
	if MEMInfo == nil {
		return ""
	}
	s.WriteString("MemoryInfo:\n")
	if MEMInfo.Mess != "" {
		s.WriteString(fmt.Sprintf("ProcessInfo:%s", MEMInfo.Mess))
	} else {
		s.WriteString(fmt.Sprintf("AnonMemory:%d, PhysMemory:%d, Rss:%d, Vss:%d",
			MEMInfo.Anon, MEMInfo.PhysMemory, MEMInfo.Rss, MEMInfo.Vss))
	}
	s.WriteString(fmt.Sprintf(", TimeStamp:%d", MEMInfo.TimeStamp))
	return s.String()
}

func toStringFPSInfo(FPSInfo *giDevice.FPSInfo) string {
	var s strings.Builder
	if FPSInfo == nil {
		return ""
	}
	s.WriteString(fmt.Sprintf("FPSInfo:\nFPS:%d, TimeStamp:%d", FPSInfo.FPS, FPSInfo.TimeStamp))
	return s.String()
}

func toStringNetwork(Network *giDevice.NetWorkingInfo) string {
	var s strings.Builder
	if Network == nil {
		return ""
	}
	s.WriteString(fmt.Sprintf("NetworkingInfo:\nRxBytes:%d, RxPackets:%d, TxBytes:%d, TxPackets:%d",
		Network.RxBytes, Network.RxPackets, Network.TxBytes, Network.TxPackets))
	s.WriteString(fmt.Sprintf(", TimeStamp:%d", Network.TimeStamp))
	return s.String()
}
