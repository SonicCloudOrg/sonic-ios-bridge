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
	giDevice "github.com/SonicCloudOrg/sonic-gidevice"
	"github.com/SonicCloudOrg/sonic-ios-bridge/src/util"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
)

var pefmonCmd = &cobra.Command{
	Use:   "perfmon",
	Short: "Get perfmon from your device.",
	Long:  "Get perfmon from your device.",
	RunE: func(cmd *cobra.Command, args []string) error {
		device := util.GetDeviceByUdId(udid)
		if device == nil {
			fmt.Println("device is not found")
			os.Exit(0)
		}
		util.CheckMount(device)

		var data <-chan []byte
		var err error

		if processCpu {
			addCpuAttr()
		}

		if processMem {
			addMemAttr()
		}

		if pid != -1 {
			data, err = device.PerfStart(
				giDevice.WithPerfSystemCPU(sysCPU),
				giDevice.WithPerfSystemMem(sysMEM),
				giDevice.WithPerfSystemDisk(sysDisk),
				giDevice.WithPerfSystemNetwork(sysNetwork),
				giDevice.WithPerfNetwork(processNetwork),
				giDevice.WithPerfFPS(getFPS),
				giDevice.WithPerfGPU(getGPU),
				giDevice.WithPerfProcessAttributes(processAttributes...),
				giDevice.WithPerfPID(pid),
				giDevice.WithPerfOutputInterval(refreshTime),
			)
		} else if bundleId != "" {
			data, err = device.PerfStart(
				giDevice.WithPerfSystemCPU(sysCPU),
				giDevice.WithPerfSystemMem(sysMEM),
				giDevice.WithPerfSystemDisk(sysDisk),
				giDevice.WithPerfSystemNetwork(sysNetwork),
				giDevice.WithPerfNetwork(processNetwork),
				giDevice.WithPerfFPS(getFPS),
				giDevice.WithPerfGPU(getGPU),
				giDevice.WithPerfBundleID(bundleId),
				giDevice.WithPerfProcessAttributes(processAttributes...),
				giDevice.WithPerfOutputInterval(refreshTime),
			)
		} else {
			data, err = device.PerfStart(
				giDevice.WithPerfSystemCPU(sysCPU),
				giDevice.WithPerfSystemMem(sysMEM),
				giDevice.WithPerfSystemDisk(sysDisk),
				giDevice.WithPerfSystemNetwork(sysNetwork),
				giDevice.WithPerfFPS(getFPS),
				giDevice.WithPerfGPU(getGPU),
				giDevice.WithPerfOutputInterval(refreshTime),
			)
		}

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
				device.PerfStop()
				fmt.Println("force end perfmon")
				os.Exit(0)
			case d := <-data:
				fmt.Println(string(d))
			}
		}
		return nil
	},
}

var (
	sysCPU            bool
	getGPU            bool
	sysMEM            bool
	getFPS            bool
	sysDisk           bool
	sysNetwork        bool
	pid               int
	bundleId          string
	processNetwork    bool
	processCpu        bool
	processMem        bool
	processAttributes []string
	refreshTime       int
)

func addCpuAttr() {
	processAttributes = append(processAttributes, "cpuUsage")
}

func addMemAttr() {
	processAttributes = append(processAttributes, "memVirtualSize", "physFootprint", "memResidentSize", "memAnon")
}

func init() {
	rootCmd.AddCommand(pefmonCmd)
	pefmonCmd.Flags().StringVarP(&udid, "udid", "u", "", "device's serialNumber ( default first device )")
	pefmonCmd.Flags().IntVarP(&pid, "pid", "p", -1, "get PID data")
	pefmonCmd.Flags().StringVarP(&bundleId, "bundleId", "b", "", "target bundleId")
	pefmonCmd.Flags().BoolVar(&sysCPU, "sys-cpu", false, "get system cpu data")
	pefmonCmd.Flags().BoolVar(&sysMEM, "sys-mem", false, "get system memory data")
	pefmonCmd.Flags().BoolVar(&sysDisk, "sys-disk", false, "get system disk data")
	pefmonCmd.Flags().BoolVar(&sysNetwork, "sys-network", false, "get system networking data")
	pefmonCmd.Flags().BoolVar(&getGPU, "gpu", false, "get gpu data")
	pefmonCmd.Flags().BoolVar(&getFPS, "fps", false, "get fps data")
	pefmonCmd.Flags().BoolVar(&processNetwork, "proc-network", false, "get process network data")
	pefmonCmd.Flags().BoolVar(&processCpu, "proc-cpu", false, "get process cpu data")
	pefmonCmd.Flags().BoolVar(&processMem, "proc-mem", false, "get process mem data")
	pefmonCmd.Flags().IntVarP(&refreshTime, "refresh", "r", 1000, "data refresh time(millisecond)")
	pefmonCmd.Flags().BoolVarP(&isFormat, "format", "f", false, "convert to JSON string and format")
	pefmonCmd.Flags().BoolVarP(&isJson, "json", "j", false, "convert to JSON string")
}
