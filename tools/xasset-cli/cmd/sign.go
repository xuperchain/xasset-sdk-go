package cmd

import (
	"fmt"

	"github.com/xuperchain/xasset-sdk-go/auth"
	"github.com/xuperchain/xasset-sdk-go/tools/xasset-cli/common"

	"github.com/spf13/cobra"
)

type SignCmd struct {
	BaseCmd
}

func GetSignCmd() *SignCmd {
	cmdIns := new(SignCmd)

	cmdIns.Cmd = &cobra.Command{
		Use:           "sign",
		Short:         "Xasset sign client.",
		Example:       common.CmdLineName + " sign ecsda [arguments]",
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	cmdIns.Cmd.AddCommand(GetEcsdaSignCmd().GetCmd())

	return cmdIns
}

// ecsda sign command
type EcsdaSignCmd struct {
	BaseCmd
	// 要绑定的变量类型只能使用内置基础类型
	PrivateKey string
	Msg        string
	Fmt        string
}

func GetEcsdaSignCmd() *EcsdaSignCmd {
	cmdIns := new(EcsdaSignCmd)

	cmdIns.Cmd = &cobra.Command{
		Use:           "ecsda",
		Short:         "Ecsda sign.",
		Example:       common.CmdLineName + " sign ecsda -k [private key] -m [msg] -f vis",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmdIns.Sign()
		},
	}

	// 设置命令行参数并绑定变量
	cmdIns.Cmd.Flags().StringVarP(&cmdIns.PrivateKey, "privtkey", "k", "", "account private key.")
	cmdIns.Cmd.Flags().StringVarP(&cmdIns.Msg, "msg", "m", "", "content to be signed")
	cmdIns.Cmd.Flags().StringVarP(&cmdIns.Fmt, "fmt", "f", "vis", "display format. std|vis")

	return cmdIns
}

// print new account
func (t *EcsdaSignCmd) Sign() error {
	sign, err := auth.XassetSignECDSA(t.PrivateKey, []byte(t.Msg))
	if err != nil {
		fmt.Print(common.FailedRespMsg)
		return nil
	}

	switch t.Fmt {
	case "std":
		fmt.Print(sign)
	case "vis":
		fmt.Printf("%s\n", sign)
	default:
		fmt.Print(common.FailedRespMsg)
	}

	return nil
}
