package cmd

import (
	"github.com/SonicCloudOrg/sonic-ios-bridge/cmd/afcUtil"
	"github.com/spf13/cobra"
)

var afcCmd = &cobra.Command{
	Use:   "afc",
	Short: "manipulate device files through afc commands",
	Long:  "manipulate device files through afc commands",
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.Help()
		return nil
	},
}

func init()  {
	rootCmd.AddCommand(afcCmd)
	afcUtil.InitAfc(afcCmd)
}
