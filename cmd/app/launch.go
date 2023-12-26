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
package app

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	giDevice "github.com/SonicCloudOrg/sonic-gidevice"
	"github.com/SonicCloudOrg/sonic-ios-bridge/src/util"
	envparse "github.com/hashicorp/go-envparse"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var bKillExisting bool
var envKVs []string

var launchCmd = &cobra.Command{
	Use:   "launch",
	Args:  cobra.ArbitraryArgs,
	Short: "Launch App",
	Long:  "Launch App",
	RunE: func(cmd *cobra.Command, args []string) error {
		device := util.GetDeviceByUdId(udid)
		if device == nil {
			os.Exit(0)
		}
		argv := []any{}
		for _, arg := range args {
			argv = append(argv, arg)
		}
		inputEnvVars, errParseEnv := envparse.Parse(bytes.NewReader([]byte(strings.Join(envKVs, "\n"))))
		if errParseEnv != nil {
			logrus.Warnf("Failed to parse env vars: %+v", errParseEnv)
		}
		envVars := map[string]any{}
		for k, v := range inputEnvVars {
			envVars[k] = v
		}
		_, errLaunch := device.AppLaunch(bundleId, giDevice.WithArguments(argv), giDevice.WithKillExisting(bKillExisting), giDevice.WithEnvironment(envVars))
		if errLaunch != nil {
			return util.NewErrorPrint(util.ErrSendCommand, "launch", errLaunch)
		}
		return nil
	},
}

func myHelpFunc(cmd *cobra.Command, args []string) {
	fmt.Printf(`%s
 
 Usage:
   %s -- [arguments [arguments ...]]
 
 Flags:
 %s`, cmd.Long, cmd.UseLine(), cmd.Flags().FlagUsages())
}

func initAppLaunch() {
	appRootCMD.AddCommand(launchCmd)
	launchCmd.Flags().StringVarP(&udid, "udid", "u", "", "device's serialNumber")
	launchCmd.Flags().StringVarP(&bundleId, "bundleId", "b", "", "target bundleId")
	launchCmd.MarkFlagRequired("bundleId")
	launchCmd.Flags().StringSliceVarP(&envKVs, "env", "e", []string{}, "environment variables; format: KEY=VALUE")
	launchCmd.Flags().BoolVar(&bKillExisting, "kill-existing", false, "kill the application if it is already running")
	launchCmd.SetHelpFunc(myHelpFunc)
	launchCmd.UseLine()
}
