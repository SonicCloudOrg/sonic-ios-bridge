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
	"github.com/SonicCloudOrg/sonic-ios-bridge/src/util"
	giDevice "github.com/electricbubble/gidevice"
	"github.com/spf13/cobra"
	"io"
	"os"
	gPath "path"
	"strings"
)

var afcCmd = &cobra.Command{
	Use:   "afc",
	Short: "manipulate device files through afc commands",
	Long:  "manipulate device files through afc commands",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			fmt.Println("arguments error: "+ strings.Join(args," "))
			fmt.Println()
			fmt.Println(cmd.UsageString())
			os.Exit(0)
		}
		device := util.GetDeviceByUdId(udid)
		if device == nil {
			os.Exit(0)
		}
		var afcServer giDevice.Afc
		var err error
		if bundleID != "" {
			var houseArrestSrv giDevice.HouseArrest
			houseArrestSrv, err = device.HouseArrestService()
			if err != nil {
				fmt.Println(err)
				os.Exit(0)
			}
			afcServer, err = houseArrestSrv.Documents(bundleID)
		} else {
			afcServer, err = device.AfcService()
		}
		if err != nil {
			fmt.Println(err)
			os.Exit(0)
		}
		var operator = args[0]
		switch operator {
		case "ls":
			lsShow(afcServer, args[1])
			break
		case "cat":
			catFile(afcServer, args[1])
			break
		case "rm":
			err := afcServer.Remove(args[1])
			if err != nil {
				fmt.Println(err)
				os.Exit(0)
			}
			fmt.Println("rm success")
			break
		case "stat":
			info, err := afcServer.Stat(args[1])
			if err != nil {
				os.Exit(0)
			}
			if info.IsDir() {
				fmt.Println("type:DIR")
			} else {
				fmt.Println("type:FILE")
			}
			fmt.Println("CTime:", info.CreationTime().Format("2006-01-02 15:04:05"))
			fmt.Println("MTime:", info.ModTime().Format("2006-01-02 15:04:05"))
			fmt.Println(fmt.Sprintf("Size:%d", info.Size()))
			break
		case "tree":
			showTree(afcServer, args[1], 100)
			break
		case "pull":
			if len(args) != 3 {
				fmt.Println("arguments error")
				os.Exit(0)
			}
			pullOperate(afcServer, args[1], args[2])
			fmt.Println(fmt.Sprintf("success,pull %s --> %s", args[1], args[2]))
			break
		case "push":
			if len(args) != 3 {
				fmt.Println("parameter error")
				os.Exit(0)
			}
			pushOperate(afcServer, args[1], args[2])
			fmt.Println(fmt.Sprintf("success,push %s --> %s", args[1], args[2]))
			break
		case "rmtree":
			removeTree(afcServer, args[1])
			fmt.Println("rmtree success")
			break
		case "mkdir":
			err := afcServer.Mkdir(args[1])
			if err != nil {
				fmt.Println(err)
				os.Exit(0)
			}
			fmt.Println("mkdir success")
			break
		default:
			fmt.Println("parameter error："+operator)
			fmt.Println()
			fmt.Println(cmd.UsageString())
		}
		return nil
	},
}

var bundleID string

func init() {
	rootCmd.AddCommand(afcCmd)
	afcCmd.Flags().StringVarP(&udid, "udid", "u", "", "device's serialNumber ( default first device )")
	afcCmd.Flags().StringVarP(&bundleID, "bundleId", "b", "", "app bundleId")
	afcCmd.SetUsageTemplate(strings.Replace(
		afcCmd.UsageTemplate(),
		"{{.UseLine}}{{end}}{{if .HasAvailableSubCommands}}\n  {{.CommandPath}} [command]{{end}}{{if gt (len .Aliases) 0}}",
		"sib afc [-h] [-b BUNDLE_ID] {ls,rm,cat,pull,push,stat,tree,rmtree,mkdir} arguments [arguments ...]",
		1))
}

var (
	levelFlag []bool // 路径级别标志
	fileCount,
	dirCount int
)

const (
	space  = "   "
	line   = "│  "
	last   = "└─ "
	middle = "├─ "
)

func showTree(afc giDevice.Afc, path string, subDepth int) {
	fmt.Println(gPath.Base(path))
	levelFlag = make([]bool, subDepth)
	walk(afc, path, 0)
}

func walk(afc giDevice.Afc, dir string, level int) {
	if len(levelFlag) <= level {
		fmt.Println("exceeded maximum depth")
		os.Exit(0)
	}
	levelFlag[level] = true
	if files, err := afc.ReadDir(dir); err == nil {
		for index, file := range files {
			if file == "." || file == ".." {
				continue
			}
			absFile := gPath.Join(dir, file)

			isLast := index == len(files)-1

			levelFlag[level] = !isLast
			afcInfo, err := afc.Stat(gPath.Join(dir, file))
			if err != nil {
				fmt.Println(err)
				os.Exit(0)
			}
			showLine(level, isLast, afcInfo)
			if afcInfo.IsDir() {
				walk(afc, absFile, level+1)
			}
		}
	} else {
		fmt.Println(err)
	}
}

func showLine(level int, isLast bool, info *giDevice.AfcFileInfo) {
	preFix := buildPrefix(level)
	outTemp, out := "%s%s%s", ""
	fName := info.Name()
	if info.IsDir() {
		fName = fmt.Sprintf("%s", fName)
		dirCount++
	} else {
		fileCount++
	}
	if isLast {
		out = fmt.Sprintf(outTemp, preFix, last, fName)
	} else {
		out = fmt.Sprintf(outTemp, preFix, middle, fName)
	}
	fmt.Println(out)
}

func buildPrefix(level int) string {
	result := ""
	for idx := 0; idx < level; idx++ {
		if levelFlag[idx] {
			result += line
		} else {
			result += space
		}
	}
	return result
}

func pushFile(afc giDevice.Afc, localPath string, devicePath string) {
	file, err := os.Open(localPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer func() {
		_ = file.Close()
	}()

	afcFile, err := afc.Open(devicePath, giDevice.AfcFileModeWr)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	defer func() {
		_ = afcFile.Close()
	}()
	if _, err = io.Copy(afcFile, file); err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

}

func pullFile(afc giDevice.Afc, devicePath string, localPath string) {
	afcFile, err := afc.Open(devicePath, giDevice.AfcFileModeRdOnly)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	defer func() {
		_ = afcFile.Close()
	}()
	file, err := os.Create(localPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	defer func() {
		_ = file.Close()
	}()
	if _, err = io.Copy(file, afcFile); err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
}

func pullOperate(afc giDevice.Afc, devicePath string, localPath string) {
	fileInfo, err := afc.Stat(devicePath)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	if fileInfo.IsDir() {
		localFile, err := os.ReadDir(localPath)
		if localFile == nil || err != nil {
			mkdirError := os.Mkdir(localPath, os.ModePerm)
			if mkdirError != nil {
				fmt.Println(mkdirError)
				os.Exit(0)
			}
		}
		fileNames, err := afc.ReadDir(devicePath)
		if err != nil {
			fmt.Println(err)
			os.Exit(0)
		}
		for _, file := range fileNames {
			if file == "." || file == ".." {
				continue
			}
			pullOperate(afc, gPath.Join(devicePath, file), gPath.Join(localPath, file))
		}
	} else {
		pullFile(afc, devicePath, localPath)
	}
}

func pushOperate(afc giDevice.Afc, localPath string, devicePath string) {
	localFile, err := os.Stat(localPath)
	if err != nil {
		os.Exit(0)
	}
	if localFile.IsDir() {
		aPathInfo, _ := afc.ReadDir(devicePath)
		if aPathInfo == nil {
			mkdirError := afc.Mkdir(devicePath)
			if mkdirError != nil {
				fmt.Println(mkdirError)
				os.Exit(0)
			}
		}
		childFiles, err := os.ReadDir(localPath)
		if err != nil {
			os.Exit(0)
		}
		for _, childFile := range childFiles {
			pushOperate(afc, gPath.Join(localPath, childFile.Name()), gPath.Join(devicePath, childFile.Name()))
		}
	} else {
		pushFile(afc, localPath, devicePath)
	}
}

func catFile(afc giDevice.Afc, filePath string) {
	fileInfo, err := afc.Stat(filePath)
	if err != nil {
		fmt.Println("file path is null")
		os.Exit(0)
	}
	p := make([]byte, fileInfo.Size())
	afcFile, err := afc.Open(filePath, giDevice.AfcFileModeRdOnly)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	defer func() {
		_ = afcFile.Close()
	}()
	n, err := afcFile.Read(p)
	if err == io.EOF {
		fmt.Println(err)
		os.Exit(0)
		//break
	}
	fmt.Print(string(p[:n]))
}

func lsShow(afc giDevice.Afc, filePath string) {
	fileNames, err := afc.ReadDir(filePath)
	if err != nil {
		os.Exit(0)
	}
	for _, fileName := range fileNames {
		if fileName == "." || fileName == ".." {
			continue
		}
		info, err := afc.Stat(gPath.Join(filePath, fileName))

		if err != nil {
			os.Exit(0)
		}
		if info.IsDir() {
			fmt.Println(fileName + "/")
		} else {
			fmt.Println(fmt.Sprintf("- %s %d", fileName, info.Size()))
		}
	}
}

func removeTree(afc giDevice.Afc, devicePath string) {
	fileInfo, err := afc.Stat(devicePath)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	if fileInfo.IsDir() {
		fileNames, err := afc.ReadDir(devicePath)
		if err != nil {
			fmt.Println(err)
			os.Exit(0)
		}
		for _, file := range fileNames {
			if file == "." || file == ".." {
				continue
			}
			var childPath string
			childPath = gPath.Join(devicePath, file)

			removeTree(afc, childPath)
		}

		err = afc.Remove(devicePath)
		if err != nil {
			fmt.Println(err)
			os.Exit(0)
		}
	} else {
		err := afc.Remove(devicePath)
		if err != nil {
			fmt.Println(err)
			os.Exit(0)
		}
	}
}
