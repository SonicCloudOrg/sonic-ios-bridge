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
package remote

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/SonicCloudOrg/sonic-ios-bridge/src/entity"
	"github.com/SonicCloudOrg/sonic-ios-bridge/src/util"
	"github.com/spf13/cobra"
	"io/ioutil"
	"log"
	"os"
)

var disConnectCmd = &cobra.Command{
	Use:   "disconnect",
	Short: "Disconnect remote device",
	Long:  "Disconnect remote device",
	RunE: func(cmd *cobra.Command, args []string) error {

		_, err := os.Stat(".sib")
		if err != nil {
			fmt.Println("success")
			return nil
		}

		file, err := os.OpenFile(util.RemoteInfoFilePath, os.O_RDWR, os.ModePerm)
		defer file.Close()

		if err != nil {
			fmt.Println("success")
			return nil
		}
		jsonData, err1 := ioutil.ReadAll(file)
		if err1 != nil {
			fmt.Println("success")
			return nil
		}

		remoteMap := make(map[string]*entity.RemoteInfo)

		if jsonData != nil && len(jsonData) != 0 {
			err = json.Unmarshal(jsonData, &remoteMap)
			if err != nil {
				fmt.Println("success")
				return nil
			}
		}

		addr := fmt.Sprintf("%s:%d", host, port)

		if remoteMap[addr] == nil {
			fmt.Println("no such addr")
			return nil
		}

		delete(remoteMap, fmt.Sprintf("%s:%d", host, port))

		err = file.Truncate(0)
		if err != nil {
			log.Panic(err)
		}
		_, err = file.Seek(0, 0)
		if err != nil {
			log.Panic(err)
		}

		write := bufio.NewWriter(file)

		jsonData, _ = json.Marshal(remoteMap)

		write.Write(jsonData)
		write.Flush()
		fmt.Println("success")
		return nil
	},
}

func disConnectInit() {
	remoteCmd.AddCommand(disConnectCmd)
	disConnectCmd.Flags().StringVarP(&host, "host", "i", "", "remote device host")
	disConnectCmd.Flags().IntVarP(&port, "port", "p", 9123, "share port")
	disConnectCmd.MarkFlagRequired("host")
}
