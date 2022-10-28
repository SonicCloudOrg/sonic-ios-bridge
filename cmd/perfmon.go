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
	"fmt"
	"github.com/SonicCloudOrg/sonic-ios-bridge/src/entity"
	"github.com/SonicCloudOrg/sonic-ios-bridge/src/util"
	giDevice "github.com/electricbubble/gidevice"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
	"strconv"
)

var pefmonCmd = &cobra.Command{
	Use:   "perfmon",
	Short: "Get perfmon from your device.",
	Long:  "Get perfmon from your device.",
	RunE: func(cmd *cobra.Command, args []string) error {
		device := util.GetDeviceByUdId(udid)
		if device == nil {
			os.Exit(0)
		}
		if appName != "" && pid != -1 {
			fmt.Println("pid and appName cannot be used at the same time")
			os.Exit(0)
		}
		if appName != "" {
			var err error
			pid, err = util.GetApplicationPID(device, appName)
			if err != nil {
				os.Exit(0)
			}
		}
		var opts = &giDevice.PerfmonOption{
			PID:             strconv.Itoa(pid),
			OpenChanMEM:     getMEM,
			OpenChanNetWork: getNetWork,
			OpenChanCPU:     getCPU,
			OpenChanFPS:     getFPS,
			OpenChanGPU:     getGPU,
		}
		output, cancelFunc, err := device.GetPerfmon(opts)

		if err != nil {
			fmt.Println(err)
			os.Exit(0)
		}
		done := make(chan os.Signal, 1)
		signal.Notify(done, os.Interrupt, os.Kill)
		// add timer?

		for {
			select {
			case <-done:
				if cancelFunc != nil {
					cancelFunc()
				}
				return nil
			default:
				if data, ok := <-output; ok {
					d := util.ResultData(entity.CreatePerformanceData(data))
					fmt.Println(util.Format(d, isFormat, isJson))
				}
			}
		}
		return nil
	},
}

var (
	getCPU     bool
	getGPU     bool
	getMEM     bool
	getFPS     bool
	getNetWork bool
	pid        int
	appName    string
)

func init() {
	rootCmd.AddCommand(pefmonCmd)
	pefmonCmd.Flags().StringVarP(&udid, "udid", "u", "", "device's serialNumber ( default first device )")
	pefmonCmd.Flags().IntVarP(&pid, "pid", "p", -1, "get PID data")
	pefmonCmd.Flags().StringVarP(&appName, "app-name", "a", "", "get app data ( Valid for memory and CPU only )")
	pefmonCmd.Flags().BoolVar(&getCPU, "cpu", false, "get cpu data")
	pefmonCmd.Flags().BoolVar(&getMEM, "mem", false, "get memory data")
	pefmonCmd.Flags().BoolVar(&getGPU, "gpu", false, "get gpu data")
	pefmonCmd.Flags().BoolVar(&getFPS, "fps", false, "get fps data")
	pefmonCmd.Flags().BoolVar(&getNetWork, "network", false, "get networking data")
	pefmonCmd.Flags().BoolVarP(&isFormat, "format", "f", false, "convert to JSON string and format")
	pefmonCmd.Flags().BoolVarP(&isJson, "json", "j", false, "convert to JSON string")
}
