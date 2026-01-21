package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/chzyer/readline"
	"github.com/free-ran-ue/free-ran-ue/v2/model"
	fruUtil "github.com/free-ran-ue/free-ran-ue/v2/util"
	"github.com/free-ran-ue/frush/constant"
	"github.com/free-ran-ue/frush/manager"
	"github.com/free-ran-ue/frush/subscriber"
)

func printFrush() {
	fmt.Print(`=============== Welcome to use frush! ===============
======███████╗██████╗ ██╗   ██╗███████╗██╗  ██╗======
======██╔════╝██╔══██╗██║   ██║██╔════╝██║  ██║======
======█████╗  ██████╔╝██║   ██║███████╗███████║======
======██╔══╝  ██╔══██╗██║   ██║╚════██║██╔══██║======
======██║     ██║  ██║╚██████╔╝███████║██║  ██║======
======╚═╝     ╚═╝  ╚═╝ ╚═════╝ ╚══════╝╚═╝  ╚═╝======
=====================================================
`)
}

func usage() {
	fmt.Print(`
Commands:
	help    Show help
	exit    Exit

	add    Add a subscriber
	delete Delete a subscriber

	status  Show the status of gNB and UE
	gnb     Start gNB
	reg     Register UE
	dereg   De-register UE
`)
}

func getConfig(gnbConfigPath, ueConfigPath string) (*model.GnbConfig, *model.UeConfig, error) {
	gnbConfig := model.GnbConfig{}
	if err := fruUtil.LoadFromYaml(gnbConfigPath, &gnbConfig); err != nil {
		return nil, nil, err
	}
	if err := fruUtil.ValidateGnb(&gnbConfig); err != nil {
		panic(err)
	}

	ueConfig := model.UeConfig{}
	if err := fruUtil.LoadFromYaml(ueConfigPath, &ueConfig); err != nil {
		return nil, nil, err
	}
	if err := fruUtil.ValidateUe(&ueConfig); err != nil {
		panic(err)
	}

	return &gnbConfig, &ueConfig, nil
}

func printStatusTable(gnbName, ueImsi string, gnbStatus, ueStatus constant.ContextStatus) {
	nameHeader := "Name"
	maxNameLen := len(nameHeader)

	if len(gnbName) > maxNameLen {
		maxNameLen = len(gnbName)
	}
	if len(ueImsi) > maxNameLen {
		maxNameLen = len(ueImsi)
	}

	nameWidth := maxNameLen + 2
	stateWidth := 19

	fmt.Println("┌" + strings.Repeat("─", nameWidth) + "┬" + strings.Repeat("─", stateWidth) + "┐")
	fmt.Printf("│ %-*s│ %-*s│\n", nameWidth-1, nameHeader, stateWidth-1, "State")
	fmt.Println("├" + strings.Repeat("─", nameWidth) + "┼" + strings.Repeat("─", stateWidth) + "┤")

	fmt.Printf("│ %-*s│ %-*s│\n", nameWidth-1, gnbName, stateWidth-1, gnbStatus)
	fmt.Printf("│ %-*s│ %-*s│\n", nameWidth-1, ueImsi, stateWidth-1, ueStatus)

	fmt.Println("└" + strings.Repeat("─", nameWidth) + "┴" + strings.Repeat("─", stateWidth) + "┘")
}

func checkRoot() {
	if os.Geteuid() != 0 {
		fmt.Println("Please run as root: sudo ./frush")
		os.Exit(1)
	}
}

func main() {
	checkRoot()

	printFrush()

	gnbConfig, ueConfig, err := getConfig(constant.TEMPLATE_GNB_YAML, constant.TEMPLATE_UE_YAML)
	if err != nil {
		panic(err)
	}

	frushManager := manager.NewManager(*gnbConfig, *ueConfig)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	rl, err := readline.New(constant.CMD_START)
	if err != nil {
		panic(err)
	}
	defer func() {
		cancel()

		if frushManager.GnbContext().GetStatus() == constant.CONTEXT_STATUS_GNB_RUNNING {
			frushManager.GnbContext().Stop()
		}
		if frushManager.UeContext().GetStatus() == constant.CONTEXT_STATUS_UE_REGISTERED {
			frushManager.UeContext().Stop()
		}
		if err := rl.Close(); err != nil {
			panic(err)
		}
	}()

	for {
		line, err := rl.Readline()
		if err != nil {
			if err == readline.ErrInterrupt {
				fmt.Println(constant.SYSTEM_HINT_CTRL_C_EXIT)
				continue
			}
			panic(err)
		}

		cmds := strings.Fields(line)
		if len(cmds) == 0 {
			continue
		}

		switch cmds[0] {
		case constant.CMD_HELP:
			usage()
		case constant.CMD_EXIT:
			return
		case constant.CMD_ADD_SUBSCRIBER:
			if err := subscriber.AddSubscriber(constant.TEMPLATE_CONSOLE_ACCOUNT_JSON, constant.TEMPLATE_SUBSCRIBER_JSON); err != nil {
				fmt.Println(err)
			} else {
				fmt.Println(constant.OUTPUT_SUCCESS)
			}
		case constant.CMD_DELETE_SUBSCRIBER:
			if err := subscriber.DeleteSubscriber(constant.TEMPLATE_CONSOLE_ACCOUNT_JSON, constant.TEMPLATE_SUBSCRIBER_JSON); err != nil {
				fmt.Println(err)
			} else {
				fmt.Println(constant.OUTPUT_SUCCESS)
			}
		case constant.CMD_STATUS:
			printStatusTable(frushManager.GnbContext().GetName(), frushManager.UeContext().GetImsi(), frushManager.GnbContext().GetStatus(), frushManager.UeContext().GetStatus())
		case constant.CMD_GNB:
			if err := frushManager.GnbContext().Start(ctx); err != nil {
				fmt.Println(err)
			} else {
				time.Sleep(100 * time.Millisecond)
				fmt.Println(constant.OUTPUT_SUCCESS)
			}
		case constant.CMD_UE_REGISTER:
			if err := frushManager.UeContext().Start(ctx); err != nil {
				fmt.Println(err)
			} else {
				time.Sleep(200 * time.Millisecond)
				fmt.Println(constant.OUTPUT_SUCCESS)
			}
		case constant.CMD_UE_DE_REGISTER:
			frushManager.UeContext().Stop()
			time.Sleep(100 * time.Millisecond)
			fmt.Println(constant.OUTPUT_SUCCESS)
		default:
			fmt.Println(fmt.Sprintf(constant.SYSTEM_HINT_UNKNOWN_CMD, cmds[0]))
			usage()
		}
	}
}
