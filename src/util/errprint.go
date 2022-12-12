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
