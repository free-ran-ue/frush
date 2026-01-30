package cmd

import (
	"fmt"
	"strings"

	"github.com/free-ran-ue/frush/constant"
	"github.com/free-ran-ue/frush/manager"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   constant.CMD_STATUS,
	Short: "Show the status of gNB and UE",
	Long:  "Show the status of gNB and UE",
	Run:   statusFunc,
}

func init() {
	rootCmd.AddCommand(statusCmd)
}

func statusFunc(cmd *cobra.Command, args []string) {
	nameHeader := "Name"
	maxNameLen := len(nameHeader)

	if len(manager.Manager.GnbContext().GetName()) > maxNameLen {
		maxNameLen = len(manager.Manager.GnbContext().GetName())
	}
	if len(manager.Manager.UeContext().GetImsi()) > maxNameLen {
		maxNameLen = len(manager.Manager.UeContext().GetImsi())
	}

	nameWidth := maxNameLen + 2
	stateWidth := 19

	fmt.Println("┌" + strings.Repeat("─", nameWidth) + "┬" + strings.Repeat("─", stateWidth) + "┐")
	fmt.Printf("│ %-*s│ %-*s│\n", nameWidth-1, nameHeader, stateWidth-1, "State")
	fmt.Println("├" + strings.Repeat("─", nameWidth) + "┼" + strings.Repeat("─", stateWidth) + "┤")

	fmt.Printf("│ %-*s│ %-*s│\n", nameWidth-1, manager.Manager.GnbContext().GetName(), stateWidth-1, manager.Manager.GnbContext().GetStatus())
	fmt.Printf("│ %-*s│ %-*s│\n", nameWidth-1, manager.Manager.UeContext().GetImsi(), stateWidth-1, manager.Manager.UeContext().GetStatus())

	fmt.Println("└" + strings.Repeat("─", nameWidth) + "┴" + strings.Repeat("─", stateWidth) + "┘")
	fmt.Println(constant.OUTPUT_SUCCESS)
}
