// Copyright 2018 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// See the License for the specific language governing permissions and
// limitations under the License.

package worker

import (
	"fmt"
	"strings"

	"github.com/juju/errors"
	"github.com/pingcap/tidb-enterprise-tools/dm/ctl/common"
	"github.com/pingcap/tidb-enterprise-tools/dm/pb"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
)

// NewSQLSkipCmd creates a SQLSkip command
func NewSQLSkipCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sql-skip <sub_task_name> <binlog_pos>",
		Short: "sql-skip skips specified binlog position",
		Run:   sqlSkipFunc,
	}
	return cmd
}

func sqlSkipFunc(cmd *cobra.Command, _ []string) {
	if len(cmd.Flags().Args()) != 2 {
		fmt.Println(cmd.Usage())
		return
	}
	subTaskName := cmd.Flags().Arg(0)
	if strings.TrimSpace(subTaskName) == "" {
		common.PrintLines("sub_task_name is empty")
		return
	}
	binlogPos := cmd.Flags().Arg(1)
	if err := common.CheckBinlogPos(binlogPos); err != nil {
		common.PrintLines("check binlog pos err %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	cli := common.WorkerClient()
	resp, err := cli.HandleSQLs(ctx, &pb.HandleSubTaskSQLsRequest{
		Name:      subTaskName,
		Op:        pb.SQLOp_SKIP,
		BinlogPos: binlogPos,
	})
	if err != nil {
		common.PrintLines("%s can not sql sql:\n%v", subTaskName, errors.ErrorStack(err))
		return
	}

	common.PrettyPrintResponse(resp)
}
