package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/xuperchain/xasset-sdk-go/auth"
	"github.com/xuperchain/xasset-sdk-go/tools/xasset-cli/common"

	"github.com/spf13/cobra"
)

type AccountCmd struct {
	BaseCmd
}

func GetAccountCmd() *AccountCmd {
	cmdIns := new(AccountCmd)

	cmdIns.Cmd = &cobra.Command{
		Use:           "account",
		Short:         "Xasset account operation.",
		Example:       common.CmdLineName + " account create [arguments]",
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	cmdIns.Cmd.AddCommand(GetCreateAccCmd().GetCmd())
	cmdIns.Cmd.AddCommand(GetRetrieveAccCmd().GetCmd())

	return cmdIns
}

// new account command
type CreateAccountCmd struct {
	BaseCmd
	// 要绑定的变量类型只能使用内置基础类型
	Strgth int
	Lang   int
	Fmt    string
}

func GetCreateAccCmd() *CreateAccountCmd {
	cmdIns := new(CreateAccountCmd)

	cmdIns.Cmd = &cobra.Command{
		Use:           "create",
		Short:         "Create new account.",
		Example:       common.CmdLineName + " account create -s 2 -l 2 -f vis",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmdIns.CreateAccount()
		},
	}

	// 设置命令行参数并绑定变量
	cmdIns.Cmd.Flags().IntVarP(&cmdIns.Strgth, "strgth", "s", 1, "mnemonic words strength. 1|2|3")
	cmdIns.Cmd.Flags().IntVarP(&cmdIns.Lang, "lang", "l", 1, "mnemonic words language. 1|2")
	cmdIns.Cmd.Flags().StringVarP(&cmdIns.Fmt, "fmt", "f", "vis", "display format. std|vis")

	return cmdIns
}

// print new account
func (t *CreateAccountCmd) CreateAccount() error {
	acc, err := auth.NewXchainEcdsaAccount(auth.MnemStrgth(t.Strgth), auth.MnemLang(t.Lang))
	if err != nil {
		fmt.Print(common.FailedRespMsg)
		return nil
	}

	switch t.Fmt {
	case "std":
		t.showJson(acc)
	case "vis":
		t.showVisual(acc)
	default:
		fmt.Print(common.FailedRespMsg)
	}

	return nil
}

func (t *CreateAccountCmd) showJson(acc *auth.Account) {
	js, err := json.Marshal(acc)
	if err != nil {
		fmt.Print(common.FailedRespMsg)
		return
	}

	fmt.Print(string(js))
}

func (t *CreateAccountCmd) showVisual(acc *auth.Account) {
	fmt.Printf("address:%s\n", acc.Address)
	fmt.Printf("private_key:%s\n", acc.PrivateKey)
	fmt.Printf("public_key:%s\n", acc.PublicKey)
	fmt.Printf("mnemonic:%s\n", acc.Mnemonic)
}

// retrieve account command
type RetrieveAccountCmd struct {
	BaseCmd
	// 要绑定的变量类型只能使用内置基础类型
	Mnemonic string
	Lang     int
	Fmt      string
}

func GetRetrieveAccCmd() *RetrieveAccountCmd {
	cmdIns := new(RetrieveAccountCmd)

	cmdIns.Cmd = &cobra.Command{
		Use:           "retrieve",
		Short:         "Retrieve account from mnemonic.",
		Example:       common.CmdLineName + " account retrieve -m 'your mnemonic' -l 1 -f vis",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmdIns.RetrieveAccount()
		},
	}

	// 设置命令行参数并绑定变量
	cmdIns.Cmd.Flags().StringVarP(&cmdIns.Mnemonic, "mnemonic", "m", "", "mnemonic words")
	cmdIns.Cmd.Flags().IntVarP(&cmdIns.Lang, "lang", "l", 1, "mnemonic words language. 1|2")
	cmdIns.Cmd.Flags().StringVarP(&cmdIns.Fmt, "fmt", "f", "vis", "display format. std|vis")

	return cmdIns
}

// print retrieve account
func (t *RetrieveAccountCmd) RetrieveAccount() error {
	acc, err := auth.RetrieveAccountByMnemonic(t.Mnemonic, t.Lang)
	if err != nil {
		fmt.Print(common.FailedRespMsg)
		return nil
	}

	switch t.Fmt {
	case "std":
		t.showJson(acc)
	case "vis":
		t.showVisual(acc)
	default:
		fmt.Print(common.FailedRespMsg)
	}

	return nil
}

func (t *RetrieveAccountCmd) showJson(acc *auth.Account) {
	js, err := json.Marshal(acc)
	if err != nil {
		fmt.Print(common.FailedRespMsg)
		return
	}

	fmt.Print(string(js))
}

func (t *RetrieveAccountCmd) showVisual(acc *auth.Account) {
	fmt.Printf("address:%s\n", acc.Address)
	fmt.Printf("private_key:%s\n", acc.PrivateKey)
	fmt.Printf("public_key:%s\n", acc.PublicKey)
	fmt.Printf("mnemonic:%s\n", acc.Mnemonic)
}
