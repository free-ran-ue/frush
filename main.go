package main

import (
	"context"
	"fmt"
	"os"
	"strings"

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

func printStatusTable(gnbName, ueName string, gnbStatus, ueStatus constant.ContextStatus) {
	fmt.Println("┌──────────┬─────────────────────┐")
	fmt.Println("│ Name     │ State               │")
	fmt.Println("├──────────┼─────────────────────┤")
	fmt.Printf("│ %-8s │ %-19s │\n", gnbName, gnbStatus)
	fmt.Printf("│ %-8s │ %-19s │\n", ueName, ueStatus)
	fmt.Println("└──────────┴─────────────────────┘")
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
		if frushManager.GnbContext().GetStatus() == constant.Context_Running {
			frushManager.GnbContext().Stop()
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
			printStatusTable(frushManager.GnbContext().GetName(), "", frushManager.GnbContext().GetStatus(), "")
		case constant.CMD_GNB:
			if err := frushManager.GnbContext().Start(ctx); err != nil {
				fmt.Println(err)
			} else {
				fmt.Println(constant.OUTPUT_SUCCESS)
			}
		default:
			fmt.Println(fmt.Sprintf(constant.SYSTEM_HINT_UNKNOWN_CMD, cmds[0]))
			usage()
		}
	}
}
