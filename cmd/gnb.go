package cmd

import (
	"fmt"
	"time"

	"github.com/free-ran-ue/frush/constant"
	"github.com/free-ran-ue/frush/manager"
	"github.com/spf13/cobra"
)

var gnbCmd = &cobra.Command{
	Use:   constant.CMD_GNB,
	Short: "Start gNB",
	Long:  "Start gNB",
	Run:   gnbFunc,
}

func init() {
	rootCmd.AddCommand(gnbCmd)
}

func gnbFunc(cmd *cobra.Command, args []string) {
	if err := manager.Manager.GnbContext().Start(manager.RootCtx); err != nil {
		fmt.Println(err)
		fmt.Println(constant.OUTPUT_FAILURE)
	} else {
		time.Sleep(constant.WAIT_TIME)
		fmt.Println(constant.OUTPUT_SUCCESS)
	}
}
