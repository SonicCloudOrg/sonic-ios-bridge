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
package afc

import (
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
			util.NewErrorPrint(util.ErrSendCommand, "house arrest", err)
			os.Exit(0)
		}
		afcServer, err = houseArrestSrv.Documents(bundleId)
	} else {
		afcServer, err = device.AfcService()
	}
	if err != nil {
		util.NewErrorPrint(util.ErrSendCommand, "afc", err)
		os.Exit(0)
	}
	return afcServer
}
