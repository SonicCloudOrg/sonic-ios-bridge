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
package cmd

import (
	"fmt"
	"github.com/SonicCloudOrg/sonic-ios-bridge/src/entity"
	"github.com/SonicCloudOrg/sonic-ios-bridge/src/util"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"os"

	"github.com/spf13/cobra"
)

var ipCmd = &cobra.Command{
	Use:   "ip",
	Short: "IP monitor for your device",
	Long:  "IP monitor for your device",
	RunE: func(cmd *cobra.Command, args []string) error {
		device := util.GetDeviceByUdId(udid)
		if device == nil {
			os.Exit(0)
		}
		mac, err := device.GetValue("", "WiFiAddress")
		if err != nil {
			return util.NewErrorPrint(util.ErrSendCommand, "get value", err)
		}
		macStr, _ := mac.(string)
		info := entity.NetworkInfo{}
		info.Mac = macStr
		resultBytes, err := device.Pcap()
		if err != nil {
			return util.NewErrorPrint(util.ErrSendCommand, "pcap", err)
		}
		for {
			select {
			case data, ok := <-resultBytes:
				if ok {
					err = findIP(data, &info)
					if err != nil {
						return err
					}
					if info.Mac != "" && info.IPv6 != "" && info.IPv4 != "" {
						data := util.ResultData(info)
						fmt.Println(util.Format(data, isFormat, isJson))
						return nil
					}
				}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(ipCmd)
	ipCmd.Flags().StringVarP(&udid, "udid", "u", "", "device's serialNumber ( default first device )")
	ipCmd.Flags().BoolVarP(&isJson, "json", "j", false, "convert to JSON string")
	ipCmd.Flags().BoolVarP(&isFormat, "format", "f", false, "convert to JSON string and format")
}

// findIP @link https://github.com/danielpaulus/go-ios/blob/main/ios/pcap/ipfinder.go
func findIP(p []byte, info *entity.NetworkInfo) error {
	packet := gopacket.NewPacket(p, layers.LayerTypeEthernet, gopacket.Default)
	if tcpLayer := packet.Layer(layers.LayerTypeEthernet); tcpLayer != nil {
		tcp, _ := tcpLayer.(*layers.Ethernet)
		if tcp.SrcMAC.String() == info.Mac {
			if ipv4Layer := packet.Layer(layers.LayerTypeIPv4); ipv4Layer != nil {
				ipv4, ok := ipv4Layer.(*layers.IPv4)
				if ok {
					info.IPv4 = ipv4.SrcIP.String()
				}
			}
			if ipv6Layer := packet.Layer(layers.LayerTypeIPv6); ipv6Layer != nil {
				ipv6, ok := ipv6Layer.(*layers.IPv6)
				if ok {
					info.IPv6 = ipv6.SrcIP.String()
				}
			}
		}
	}
	return nil
}
