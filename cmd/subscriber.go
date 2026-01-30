package cmd

import (
	"fmt"

	"github.com/free-ran-ue/frush/constant"
	"github.com/free-ran-ue/frush/subscriber"
	"github.com/spf13/cobra"
)

var addSubscriberCmd = &cobra.Command{
	Use:   constant.CMD_ADD_SUBSCRIBER,
	Short: "Add a subscriber",
	Long:  "Add a subscriber",
	Run:   addSubscriberFunc,
}

func init() {
	rootCmd.AddCommand(addSubscriberCmd)
}

func addSubscriberFunc(cmd *cobra.Command, args []string) {
	if err := subscriber.AddSubscriber(constant.TEMPLATE_CONSOLE_ACCOUNT_JSON, constant.TEMPLATE_SUBSCRIBER_JSON); err != nil {
		fmt.Println(err)
		fmt.Println(constant.OUTPUT_FAILURE)
	} else {
		fmt.Println(constant.OUTPUT_SUCCESS)
	}
}

var deleteSubscriberCmd = &cobra.Command{
	Use:   constant.CMD_DELETE_SUBSCRIBER,
	Short: "Delete a subscriber",
	Long:  "Delete a subscriber",
	Run:   deleteSubscriberFunc,
}

func init() {
	rootCmd.AddCommand(deleteSubscriberCmd)
}

func deleteSubscriberFunc(cmd *cobra.Command, args []string) {
	if err := subscriber.DeleteSubscriber(constant.TEMPLATE_CONSOLE_ACCOUNT_JSON, constant.TEMPLATE_SUBSCRIBER_JSON); err != nil {
		fmt.Println(err)
		fmt.Println(constant.OUTPUT_FAILURE)
	} else {
		fmt.Println(constant.OUTPUT_SUCCESS)
	}
}
