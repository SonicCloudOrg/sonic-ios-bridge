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
package util

import "fmt"

const (
	ErrConnect     = "failed connecting to"
	ErrReadingMsg  = "failed reading msg"
	ErrSendCommand = "failed send the command"
	ErrMissingArgs = "missing arg(s)"
	ErrUnknown     = "unknown error"
	MountTips      = "you can use [sib mount] command to fix it and retry"
)

func NewErrorPrint(t string, msg string, err error) error {
	if len(msg) == 0 && err == nil {
		return fmt.Errorf("%s, %s", t, MountTips)
	}
	if len(msg) == 0 {
		return fmt.Errorf("%s, %s : %w", t, MountTips, err)
	}
	if err == nil {
		return fmt.Errorf("%s [%s], %s", t, msg, MountTips)
	}
	return fmt.Errorf("%s [%s], %s, err : %w", t, msg, MountTips, err)
}
