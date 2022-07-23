package cmd

import (
	"github.com/SonicCloudOrg/sonic-ios-bridge/cmd/afcUtil"
	"github.com/spf13/cobra"
	"strings"
)

var afcCmd = &cobra.Command{
	Use:   "afc",
	Short: "manipulate device files through afc commands",
	Long:  "manipulate device files through afc commands",
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}

var bundleID string

func init()  {
	rootCmd.AddCommand(afcCmd)
	afcCmd.Flags().StringVarP(&udid, "udid", "u", "", "device's serialNumber ( default first device )")
	afcCmd.Flags().StringVarP(&bundleID, "bundleId", "b", "", "app bundleId")
	afcCmd.SetUsageTemplate(strings.Replace(
		afcCmd.UsageTemplate(),
		"{{.UseLine}}{{end}}{{if .HasAvailableSubCommands}}\n  {{.CommandPath}} [command]{{end}}{{if gt (len .Aliases) 0}}",
		"sib afc [-h] [-b BUNDLE_ID] {ls,rm,cat,pull,push,stat,tree,rmtree,mkdir} arguments [arguments ...]",
		1))
	afcUtil.InitAfc(afcCmd,udid,bundleID)
}
