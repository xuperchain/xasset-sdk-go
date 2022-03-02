package cmd

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"

	"github.com/xuperchain/xasset-sdk-go/tools/xasset-cli/common"

	"github.com/spf13/cobra"
	"github.com/xuperchain/xasset-sdk-go/auth"
)

type HashCmd struct {
	BaseCmd
}

func GetHashCmd() *HashCmd {
	cmdIns := new(HashCmd)

	cmdIns.Cmd = &cobra.Command{
		Use:           "hash",
		Short:         "Xasset hash client.",
		Example:       common.CmdLineName + " hash sm3 [arguments]",
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	cmdIns.Cmd.AddCommand(GetSM3HashCmd().GetCmd())

	return cmdIns
}

// sm3 hash command
type SM3HashCmd struct {
	BaseCmd
	// 要绑定的变量类型只能使用内置基础类型
	FilePath string
}

func GetSM3HashCmd() *SM3HashCmd {
	cmdIns := new(SM3HashCmd)

	cmdIns.Cmd = &cobra.Command{
		Use:           "sm3",
		Short:         "SM3 sign.",
		Example:       common.CmdLineName + " hash sm3 -f [file path]",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmdIns.Hash()
		},
	}

	// 设置命令行参数并绑定变量
	cmdIns.Cmd.Flags().StringVarP(&cmdIns.FilePath, "file", "f", "", "path of file to be hashed")
	return cmdIns
}

// print new account
func (t *SM3HashCmd) Hash() error {
	content, err := ioutil.ReadFile(t.FilePath)
	if err != nil {
		fmt.Println(common.FailedRespMsg)
		return nil
	}
	hash := auth.HashBySM3(content)
	fmt.Println(hex.EncodeToString(hash))
	return nil
}
