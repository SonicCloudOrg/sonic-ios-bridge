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
	"github.com/spf13/cobra"
	"os"
	gPath "path"
)

var afcTreeCmd = &cobra.Command{
	Use:   "tree",
	Short: "tree structure view directory",
	Long:  "tree structure view directory",
	RunE: func(cmd *cobra.Command, args []string) error {
		afcServer := getAFCServer()
		showTree(afcServer, treeDir, 100)
		return nil
	},
}

var treeDir string

func initTree() {
	afcRootCMD.AddCommand(afcTreeCmd)
	afcTreeCmd.Flags().StringVarP(&udid, "udid", "u", "", "device's serialNumber ( default first device )")
	afcTreeCmd.Flags().StringVarP(&bundleId, "bundleId", "b", "", "app bundleId")
	afcTreeCmd.Flags().StringVarP(&treeDir, "folder", "f", "", "folder path to tree view")
	afcTreeCmd.MarkFlagRequired("folder")
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
