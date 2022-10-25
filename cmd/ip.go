/*
 *  Copyright (C) [SonicCloudOrg] Sonic Project
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *         http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 *
 */
package cmd

import (
	giDevice "github.com/SonicCloudOrg/sonic-gidevice"
	"github.com/SonicCloudOrg/sonic-ios-bridge/src/util"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var ipCmd = &cobra.Command{
	Use:   "ip",
	Short: "IP monitor for your device",
	Long:  "IP monitor for your device",
	Run: func(cmd *cobra.Command, args []string) {
		device := util.GetDeviceByUdId(udid)
		if device == nil {
			os.Exit(0)
		}
		GetNetworkIP(device)
	},
}

func init() {
	rootCmd.AddCommand(ipCmd)
	ipCmd.Flags().StringVarP(&udid, "udid", "u", "", "device's serialNumber ( default first device )")
}

type NetworkInfo struct {
	Mac  string
	IPv4 string
	IPv6 string
}

// GetNetworkIP @link https://github.com/danielpaulus/go-ios/blob/main/ios/pcap/ipfinder.go
func GetNetworkIP(device giDevice.Device) error {
	mac, err := device.GetValue("", "WiFiAddress")
	if err != nil {
		return util.NewErrorPrint(util.ErrSendCommand, "get mac value", err)
	}
	macStr, _ := mac.(string)
	info := NetworkInfo{}
	info.Mac = macStr
	resultBytes, err := device.Pcap()
	if err != nil {
		return err
	}
	for {
		select {
		case data, ok := <-resultBytes:
			if ok {
				err = findIP(data, &info)
				if err != nil {
					return err
				}
			}
			log.Println(info)
		}
	}
}

func findIP(p []byte, info *NetworkInfo) error {
	packet := gopacket.NewPacket(p, layers.LayerTypeEthernet, gopacket.Default)
	if tcpLayer := packet.Layer(layers.LayerTypeEthernet); tcpLayer != nil {
		tcp, _ := tcpLayer.(*layers.Ethernet)
		if tcp.SrcMAC.String() == info.Mac {
			for _, layer := range packet.Layers() {
				log.Printf("layer: %s", layer.LayerType().String())
			}
			if ipv4Layer := packet.Layer(layers.LayerTypeIPv4); ipv4Layer != nil {
				ipv4, ok := ipv4Layer.(*layers.IPv4)
				if ok {
					info.IPv4 = ipv4.SrcIP.String()
					log.Printf("ip4 found: %s", info.IPv4)
				}
			}
			if ipv6Layer := packet.Layer(layers.LayerTypeIPv6); ipv6Layer != nil {
				ipv6, ok := ipv6Layer.(*layers.IPv6)
				if ok {
					info.IPv6 = ipv6.SrcIP.String()
					log.Printf("ip6 found: %s", info.IPv6)
				}
			}
		}
	}
	return nil
}
