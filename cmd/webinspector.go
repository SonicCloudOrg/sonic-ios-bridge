package cmd

import (
	"fmt"
	"github.com/SonicCloudOrg/sonic-ios-bridge/src/webinspector"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
)

var webInspectorCmd = &cobra.Command{
	Use:   "webinspector",
	Short: "Enable iOS webinspector communication service",
	Long:  "Enable iOS webinspector communication service",
	Run: func(cmd *cobra.Command, args []string) {
		done := make(chan os.Signal, 1)
		signal.Notify(done)
		cannel := webinspector.InitWebInspectorServer(udid, port, isDebug)
		fmt.Println("service started successfully")
		go func() {
			select {
			case <-done:
				cannel()
			}
		}()

		r := gin.Default()
		r.GET("/", webinspector.PagesHandle)
		r.GET("/json", webinspector.PagesHandle)
		r.GET("/json/list", webinspector.PagesHandle)
		webinspector.SetIsAdapter(true)
		r.GET("/devtools/page/:id", webinspector.PageDebugHandle)
		r.Run(fmt.Sprintf("127.0.0.1:%d", port))
	},
}

var (
	port    int
	isDebug bool
)

func init() {
	rootCmd.AddCommand(webInspectorCmd)
	webInspectorCmd.Flags().StringVarP(&udid, "udid", "u", "", "device's serialNumber ( default first device )")
	webInspectorCmd.Flags().IntVarP(&port, "port", "p", 9222, "local proxy inspector communication port")
	webInspectorCmd.Flags().BoolVarP(&isDebug, "debug", "d", false, "whether to enable debug mode")
	//afc.InitAfc(afcCmd)
}
