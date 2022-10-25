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
package afc

import (
	"fmt"
	giDevice "github.com/SonicCloudOrg/sonic-gidevice"
	"github.com/SonicCloudOrg/sonic-ios-bridge/src/util"
	"github.com/spf13/cobra"
	"os"
)

var afcRootCMD *cobra.Command

var udid, bundleId string

// InitAfc 用于初始化，在上层中调用这个方法，否则不会正常进行初始化
func InitAfc(afcCMD *cobra.Command) {
	afcRootCMD = afcCMD

	initMkDir()

	initTree()
	initCat()
	initLs()
	initStat()

	initPullCmd()
	initPush()

	initRM()
	initRMTree()
}

func getAFCServer() (afcServer giDevice.Afc) {
	device := util.GetDeviceByUdId(udid)
	if device == nil {
		os.Exit(0)
	}
	var err error
	if bundleId != "" {
		var houseArrestSrv giDevice.HouseArrest
		houseArrestSrv, err = device.HouseArrestService()
		if err != nil {
			fmt.Println(err)
			os.Exit(0)
		}
		afcServer, err = houseArrestSrv.Documents(bundleId)
	} else {
		afcServer, err = device.AfcService()
	}
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	return afcServer
}
