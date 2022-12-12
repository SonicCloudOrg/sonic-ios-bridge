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
	"github.com/SonicCloudOrg/sonic-ios-bridge/src/webinspector"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
)

var webInspectorCmd = &cobra.Command{
	Use:   "webinspector",
	Short: "Enable iOS webinspector communication service",
	Long:  "Enable iOS webinspector communication service",
	Run: func(cmd *cobra.Command, args []string) {
		done := make(chan os.Signal, 1)
		signal.Notify(done, os.Interrupt, os.Kill)
		cancel := webinspector.InitWebInspectorServer(udid, port, isProtocolDebug, isDTXDebug)
		fmt.Println("service started successfully")
		go func() {
			select {
			case <-done:
				fmt.Println("force end of webinspector")
				cancel()
				os.Exit(0)
			}
		}()

		gin.SetMode(gin.ReleaseMode)
		r := gin.Default()
		r.GET("/", webinspector.PagesHandle)
		r.GET("/json", webinspector.PagesHandle)
		r.GET("/json/list", webinspector.PagesHandle)
		webinspector.SetIsAdapter(isAdapter)
		r.GET("/devtools/page/:id", webinspector.PageDebugHandle)
		r.Run(fmt.Sprintf("127.0.0.1:%d", port))
	},
}

var (
	port            int
	isProtocolDebug bool
	isDTXDebug      bool
	isAdapter       bool
	version         bool
)

func init() {
	rootCmd.AddCommand(webInspectorCmd)
	webInspectorCmd.Flags().StringVarP(&udid, "udid", "u", "", "device's serialNumber ( default first device )")
	webInspectorCmd.Flags().IntVarP(&port, "port", "p", 9222, "local proxy inspector communication port")
	webInspectorCmd.Flags().BoolVar(&isAdapter, "cdp", false, "whether to enable chrome devtool protocol compatibility mode ( experimental function to be improved )")
	webInspectorCmd.Flags().BoolVar(&isProtocolDebug, "protocol-debug", false, "whether to enable protocol debug mode")
	webInspectorCmd.Flags().BoolVar(&isDTXDebug, "dtx-debug", false, "whether to enable dtx debug mode")
}
