// Copyright 2021 The Alaya Network Authors
// This file is part of Alaya-Go.
//
// Alaya-Go is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// Alaya-Go is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with Alaya-Go. If not, see <http://www.gnu.org/licenses/>.

package ppos

import (
	"errors"

	"gopkg.in/urfave/cli.v1"

	"github.com/AlayaNetwork/Alaya-Go/common"
)

var (
	RestrictingCmd = cli.Command{
		Name:  "restricting",
		Usage: "use for restricting",
		Subcommands: []cli.Command{
			getRestrictingInfoCmd,
		},
	}
	getRestrictingInfoCmd = cli.Command{
		Name:   "getRestrictingInfo",
		Usage:  "4100,get restricting info,parameter:address",
		Before: netCheck,
		Action: getRestrictingInfo,
		Flags:  []cli.Flag{rpcUrlFlag, addressHRPFlag, addFlag, jsonFlag},
	}
)

func getRestrictingInfo(c *cli.Context) error {
	addstring := c.String(addFlag.Name)
	if addstring == "" {
		return errors.New("The locked position release to the account account is not set")
	}
	add, err := common.Bech32ToAddress(addstring)
	if err != nil {
		return err
	}
	return query(c, 4100, add)
}
