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
	"os"
	"os/signal"
	"strings"
)

var syslogCmd = &cobra.Command{
	Use:   "syslog",
	Short: "Get syslog from your device.",
	Long:  "Get syslog from your device.",
	RunE: func(cmd *cobra.Command, args []string) error {
		device := util.GetDeviceByUdId(udid)
		if device == nil {
			os.Exit(0)
		}
		output, err := device.Syslog()
		if err != nil {
			return util.NewErrorPrint(util.ErrSendCommand, "syslog", err)
		}
		defer device.SyslogStop()
		done := make(chan os.Signal, 1)
		signal.Notify(done, os.Interrupt, os.Kill)

		go func() {
			for line := range output {
				if len(filter) == 0 {
					fmt.Println(line)
					continue
				} else {
					if strings.Contains(line, filter) {
						fmt.Println(line)
					}
				}
			}
			done <- os.Interrupt
		}()

		<-done

		return nil
	},
}

var filter string

func init() {
	rootCmd.AddCommand(syslogCmd)
	syslogCmd.Flags().StringVarP(&udid, "udid", "u", "", "device's serialNumber ( default first device )")
	syslogCmd.Flags().StringVarP(&filter, "filter", "f", "", "filter by some message.")
}
