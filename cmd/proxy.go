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
	"github.com/SonicCloudOrg/sonic-ios-bridge/src/util"
	"github.com/spf13/cobra"
	"net"
	"os"
	"os/signal"
)

var proxyCmd = &cobra.Command{
	Use:   "proxy",
	Short: "Proxy unix process or port to your pc.",
	Long:  "Proxy unix process or port to your pc.",
	RunE: func(cmd *cobra.Command, args []string) error {
		device := util.GetDeviceByUdId(udid)
		if device == nil {
			os.Exit(0)
		}
		serverListener, err := net.Listen("tcp", fmt.Sprintf(":%d", localPort))
		if err != nil {
			return err
		}
		defer serverListener.Close()

		go util.StartProxy()(serverListener, remotePort, device)

		shutdown := make(chan os.Signal, 1)
		signal.Notify(shutdown, os.Interrupt, os.Kill)

		<-shutdown
		fmt.Println("stopped.")
		return nil
	},
}

var remotePort, localPort int

func init() {
	rootCmd.AddCommand(proxyCmd)
	proxyCmd.Flags().StringVarP(&udid, "udid", "u", "", "device's serialNumber ( default first device )")
	proxyCmd.Flags().IntVarP(&remotePort, "remote-port", "r", 8100, "remote port")
	proxyCmd.Flags().IntVarP(&localPort, "local-port", "l", 9100, "local port")
}
