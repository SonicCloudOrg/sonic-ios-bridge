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
	giDevice "github.com/SonicCloudOrg/sonic-gidevice"
	"github.com/SonicCloudOrg/sonic-ios-bridge/src/entity"
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
		
		if (pid != -1 || bundleId != "") && !sysCPU && !sysMEM && !sysDisk && !sysNetwork && !getGPU && !getFPS && !processNetwork && !processMem && !processCpu {
			sysAllParamsSet()
			processNetwork = true
			processMem = true
			processCpu = true
		}

		if (pid == -1 && bundleId == "") && !sysCPU && !sysMEM && !sysDisk && !sysNetwork && !getGPU && !getFPS {
			sysAllParamsSet()
		}

		if processCpu {
			addCpuAttr()
		}

		if processMem {
			addMemAttr()
		}

		var perfOpts = []giDevice.PerfOption{
			giDevice.WithPerfSystemCPU(sysCPU),
			giDevice.WithPerfSystemMem(sysMEM),
			giDevice.WithPerfSystemDisk(sysDisk),
			giDevice.WithPerfSystemNetwork(sysNetwork),
			giDevice.WithPerfNetwork(processNetwork),
			giDevice.WithPerfFPS(getFPS),
			giDevice.WithPerfGPU(getGPU),
			giDevice.WithPerfOutputInterval(refreshTime),
		}
		if pid != -1 {
			perfOpts = append(perfOpts, giDevice.WithPerfPID(pid))
			perfOpts = append(perfOpts, giDevice.WithPerfProcessAttributes(processAttributes...))
		} else if bundleId != "" {
			perfOpts = append(perfOpts, giDevice.WithPerfBundleID(bundleId))
			perfOpts = append(perfOpts, giDevice.WithPerfProcessAttributes(processAttributes...))
		}

		data, err := device.PerfStart(perfOpts...)

		if err != nil {
			return util.NewErrorPrint(util.ErrSendCommand, "perfmon", err)
		}

		done := make(chan os.Signal, 1)
		signal.Notify(done, os.Interrupt, os.Kill)

		for {
			select {
			case <-done:
				device.PerfStop()
				fmt.Println("force end perfmon")
				os.Exit(0)
			case d := <-data:
				p := &entity.PerfData{
					PerfDataBytes: d,
				}
				fmt.Println(util.Format(p, isFormat, isJson))
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

func sysAllParamsSet() {
	getFPS = true
	getGPU = true
	sysCPU = true
	sysMEM = true
	sysDisk = true
	sysNetwork = true
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
