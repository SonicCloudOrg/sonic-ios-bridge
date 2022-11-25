/*
 *   sonic-ios-bridge  Connect to your iOS Devices.
 *   Copyright (C) 2022 SonicCloudOrg
 *
 *   This program is free software: you can redistribute it and/or modify
 *   it under the terms of the GNU General Public License as published by
 *   the Free Software Foundation, either version 3 of the License, or
 *   (at your option) any later version.
 *
 *   This program is distributed in the hope that it will be useful,
 *   but WITHOUT ANY WARRANTY; without even the implied warranty of
 *   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 *   GNU General Public License for more details.
 *
 *   You should have received a copy of the GNU General Public License
 *   along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */
package cmd

import (
	"fmt"
	"github.com/SonicCloudOrg/sonic-ios-bridge/src/util"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"time"

	"github.com/spf13/cobra"
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

		shutdown := make(chan os.Signal, 1)
		signal.Notify(shutdown, os.Interrupt, os.Kill)

		go func(listener net.Listener) {
			for {
				var accept net.Conn
				var err error
				if accept, err = listener.Accept(); err != nil {
					log.Println("accept:", err)
				}
				fmt.Println("accept", accept.RemoteAddr())
				rInnerConn, err := device.NewConnect(remotePort)
				if err != nil {
					fmt.Println("connect to device fail")
					continue
				}
				rConn := rInnerConn.RawConn()
				rConn.SetDeadline(time.Time{})
				go func(lConn net.Conn) {
					go func(lConn, rConn net.Conn) {
						if _, err := io.Copy(lConn, rConn); err != nil {
							log.Println("local -> remote failed:", err)
						}
					}(lConn, rConn)
					go func(lConn, rConn net.Conn) {
						if _, err := io.Copy(rConn, lConn); err != nil {
							log.Println("local <- remote failed:", err)
						}
					}(lConn, rConn)
				}(accept)
			}
		}(serverListener)
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
