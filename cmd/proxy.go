package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
)

var port, target string

var proxyCmd = &cobra.Command{
	Use:   "proxy",
	Short: "Proxy port/unix path to local port.",
	Args:  cobra.MaximumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error{
		if len(args) == 0 {
			if len(port) != 0 && len(target) != 0 {
				fmt.Printf("proxy called %s %s", port, target)
				return nil
			} else {
				return errors.New("missing arg(s)")
			}
		} else if len(args) < 2 {
			return errors.New("missing arg(s)")
		} else {
			fmt.Printf("proxy called %s %s", args[0], args[1])
			return nil
		}
	},
}

func init() {
	rootCmd.AddCommand(proxyCmd)
	proxyCmd.Flags().StringVarP(&port, "port", "p", "", "local port")
	proxyCmd.Flags().StringVarP(&target, "target", "t", "", "target port/unix path")
}
