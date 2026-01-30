package cmd

import (
	"fmt"
	"time"

	"github.com/free-ran-ue/frush/constant"
	"github.com/free-ran-ue/frush/manager"
	"github.com/free-ran-ue/util"
	"github.com/spf13/cobra"
)

var registerCmd = &cobra.Command{
	Use:   constant.CMD_UE_REGISTER,
	Short: "Register UE",
	Long:  "Register UE",
	Run:   registerFunc,
}

func init() {
	rootCmd.AddCommand(registerCmd)
}

func registerFunc(cmd *cobra.Command, args []string) {
	if err := manager.Manager.UeContext().Start(manager.RootCtx); err != nil {
		fmt.Println(err)
		fmt.Println(constant.OUTPUT_FAILURE)
	} else {
		time.Sleep(constant.WAIT_TIME)
		fmt.Println(constant.OUTPUT_SUCCESS)
	}
}

var deregisterCmd = &cobra.Command{
	Use:   constant.CMD_UE_DE_REGISTER,
	Short: "Deregister UE",
	Long:  "Deregister UE",
	Run:   deregisterFunc,
}

func init() {
	rootCmd.AddCommand(deregisterCmd)
}

func deregisterFunc(cmd *cobra.Command, args []string) {
	if err := manager.Manager.UeContext().Stop(); err != nil {
		fmt.Println(err)
		fmt.Println(constant.OUTPUT_FAILURE)
	} else {
		time.Sleep(constant.WAIT_TIME)
		fmt.Println(constant.OUTPUT_SUCCESS)
	}
}

var pingCmd = &cobra.Command{
	Use:   constant.CMD_PING,
	Short: "Ping the DN, if dn is not provided, ping 1.1.1.1 and 8.8.8.8",
	Long:  "Ping the DN, if dn is not provided, ping 1.1.1.1 and 8.8.8.8",
	Run:   pingFunc,
}

func init() {
	rootCmd.AddCommand(pingCmd)
}

func pingFunc(cmd *cobra.Command, args []string) {
	switch len(args) {
	case 0:
		fmt.Printf("Pinging %s...\n", constant.DN_ONE)
		if err := manager.Manager.UeContext().Ping(constant.DN_ONE); err != nil {
			fmt.Println(err)
			fmt.Println(constant.OUTPUT_FAILURE)
		} else {
			fmt.Println(constant.OUTPUT_SUCCESS)
		}
		fmt.Printf("Pinging %s...\n", constant.DN_EIGHT)
		if err := manager.Manager.UeContext().Ping(constant.DN_EIGHT); err != nil {
			fmt.Println(err)
			fmt.Println(constant.OUTPUT_FAILURE)
		} else {
			fmt.Println(constant.OUTPUT_SUCCESS)
		}
	case 1:
		if err := util.ValidateIp(args[1]); err != nil {
			fmt.Println(err)
			fmt.Println(constant.OUTPUT_FAILURE)
		} else {
			fmt.Printf("Pinging %s...\n", args[1])
			if err := manager.Manager.UeContext().Ping(args[1]); err != nil {
				fmt.Println(err)
				fmt.Println(constant.OUTPUT_FAILURE)
			} else {
				fmt.Println(constant.OUTPUT_SUCCESS)
			}
		}
	}
}
