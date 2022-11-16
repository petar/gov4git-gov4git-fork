package cmd

import (
	"fmt"
	"os"

	"github.com/gov4git/gov4git/mod/balance"
	"github.com/gov4git/gov4git/mod/member"
	"github.com/spf13/cobra"
)

var (
	balanceCmd = &cobra.Command{
		Use:   "balance",
		Short: "Manage user balances",
		Long:  ``,
		Run:   func(cmd *cobra.Command, args []string) {},
	}

	balanceSetCmd = &cobra.Command{
		Use:   "set",
		Short: "Set user balance",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			balance.Set(
				ctx,
				setup.Community,
				member.User(balanceUser),
				balanceKey,
				balanceValue,
			)
		},
	}

	balanceGetCmd = &cobra.Command{
		Use:   "get",
		Short: "Get user balance",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			v := balance.Get(
				ctx,
				setup.Community,
				member.User(balanceUser),
				balanceKey,
			)
			fmt.Fprint(os.Stdout, v)
		},
	}

	balanceAddCmd = &cobra.Command{
		Use:   "add",
		Short: "Add to user balance",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			v := balance.Add(
				ctx,
				setup.Community,
				member.User(balanceUser),
				balanceKey,
				balanceValue,
			)
			fmt.Fprint(os.Stdout, v)
		},
	}

	balanceMulCmd = &cobra.Command{
		Use:   "mul",
		Short: "Multiply user balance",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			v := balance.Mul(
				ctx,
				setup.Community,
				member.User(balanceUser),
				balanceKey,
				balanceValue,
			)
			fmt.Fprint(os.Stdout, v)
		},
	}
)

var (
	balanceUser  string
	balanceKey   string
	balanceValue float64
)

func init() {
	// set
	balanceCmd.AddCommand(balanceSetCmd)
	balanceSetCmd.Flags().StringVar(&balanceUser, "user", "", "user alias")
	balanceSetCmd.MarkFlagRequired("user")
	balanceSetCmd.Flags().StringVar(&balanceKey, "key", "", "balance key (e.g. voting_credits)")
	balanceSetCmd.MarkFlagRequired("key")
	balanceSetCmd.Flags().Float64Var(&balanceValue, "value", 0.0, "balance value")
	balanceSetCmd.MarkFlagRequired("value")
	// get
	balanceCmd.AddCommand(balanceGetCmd)
	balanceGetCmd.Flags().StringVar(&balanceUser, "user", "", "user alias")
	balanceGetCmd.MarkFlagRequired("user")
	balanceGetCmd.Flags().StringVar(&balanceKey, "key", "", "balance key (e.g. voting_credits)")
	balanceGetCmd.MarkFlagRequired("key")
	// add
	balanceCmd.AddCommand(balanceAddCmd)
	balanceAddCmd.Flags().StringVar(&balanceUser, "user", "", "user alias")
	balanceAddCmd.MarkFlagRequired("user")
	balanceAddCmd.Flags().StringVar(&balanceKey, "key", "", "balance key (e.g. voting_credits)")
	balanceAddCmd.MarkFlagRequired("key")
	balanceAddCmd.Flags().Float64Var(&balanceValue, "value", 0.0, "value to add")
	balanceAddCmd.MarkFlagRequired("value")
	// mul
	balanceCmd.AddCommand(balanceMulCmd)
	balanceMulCmd.Flags().StringVar(&balanceUser, "user", "", "user alias")
	balanceMulCmd.MarkFlagRequired("user")
	balanceMulCmd.Flags().StringVar(&balanceKey, "key", "", "balance key (e.g. voting_credits)")
	balanceMulCmd.MarkFlagRequired("key")
	balanceMulCmd.Flags().Float64Var(&balanceValue, "value", 0.0, "value to multiply by")
	balanceMulCmd.MarkFlagRequired("value")
}
