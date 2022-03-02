package main

import (
	"fmt"

	"github.com/xuperchain/xasset-sdk-go/tools/xasset-cli/cmd"
	"github.com/xuperchain/xasset-sdk-go/tools/xasset-cli/common"

	"github.com/spf13/cobra"
)

func main() {
	rootCmd, err := NewCliCommand()
	if err != nil {
		fmt.Print(common.FailedRespMsg)
		return
	}

	if err = rootCmd.Execute(); err != nil {
		fmt.Print(common.FailedRespMsg)
		return
	}
}

func NewCliCommand() (*cobra.Command, error) {
	rootCmd := &cobra.Command{
		Use:           common.CmdLineName + " <command> [arguments]",
		Short:         common.CmdLineName + " is a xasset terminal client.",
		Long:          common.CmdLineName + " is a xasset terminal client.",
		SilenceUsage:  true,
		SilenceErrors: true,
		Example:       common.CmdLineName + " account <sub_cmd> [arguments]",
	}

	rootCmd.AddCommand(cmd.GetAccountCmd().GetCmd())
	rootCmd.AddCommand(cmd.GetSignCmd().GetCmd())
	rootCmd.AddCommand(cmd.GetHashCmd().GetCmd())

	return rootCmd, nil
}
