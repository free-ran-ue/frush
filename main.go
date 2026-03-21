package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"charm.land/lipgloss/v2"
	"github.com/chzyer/readline"
	"github.com/free-ran-ue/free-ran-ue/v2/model"
	"github.com/free-ran-ue/frush/cmd"
	"github.com/free-ran-ue/frush/constant"
	"github.com/free-ran-ue/frush/manager"
	"github.com/free-ran-ue/util"
	"golang.org/x/term"
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

func getConfig(gnbConfigPath, ueConfigPath string) (*model.GnbConfig, *model.UeConfig, error) {
	gnbConfig := model.GnbConfig{}
	if err := util.LoadFromYaml(gnbConfigPath, &gnbConfig); err != nil {
		return nil, nil, err
	}
	if err := util.ValidateGnb(&gnbConfig); err != nil {
		panic(err)
	}

	ueConfig := model.UeConfig{}
	if err := util.LoadFromYaml(ueConfigPath, &ueConfig); err != nil {
		return nil, nil, err
	}
	if err := util.ValidateUe(&ueConfig); err != nil {
		panic(err)
	}

	return &gnbConfig, &ueConfig, nil
}

func main() {
	printFrush()

	gnbConfig, ueConfig, err := getConfig(constant.TEMPLATE_GNB_YAML, constant.TEMPLATE_UE_YAML)
	if err != nil {
		panic(err)
	}

	manager.Manager = manager.NewManager(*gnbConfig, *ueConfig)
	manager.RootCtx, manager.RootCancel = context.WithCancel(context.Background())

	isTerminal := term.IsTerminal(int(os.Stdin.Fd()))
	prompt := constant.CMD_START
	if isTerminal {
		promptStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Bold(true)
		prompt = promptStyle.Render(constant.CMD_START)
	}

	completer := readline.NewPrefixCompleter(
		readline.PcItem(constant.CMD_STATUS),
		readline.PcItem(constant.CMD_GNB),
		readline.PcItem(constant.CMD_UE_REGISTER),
		readline.PcItem(constant.CMD_PING),
		readline.PcItem(constant.CMD_UE_DE_REGISTER),
		readline.PcItem(constant.CMD_DELETE_SUBSCRIBER),
		readline.PcItem(constant.CMD_ADD_SUBSCRIBER),
		readline.PcItem(constant.CMD_EXIT),
		readline.PcItem(constant.CMD_HELP),
	)

	config := &readline.Config{
		Prompt:       prompt,
		AutoComplete: completer,
	}

	rl, err := readline.NewEx(config)
	if err != nil {
		panic(err)
	}
	defer func() {
		manager.RootCancel()
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

		if cmds[0] == constant.CMD_EXIT {
			manager.RootCancel()
			errflag := false
			if manager.Manager.UeContext().GetStatus() == constant.CONTEXT_STATUS_UE_REGISTERED {
				if err := manager.Manager.UeContext().Stop(); err != nil {
					fmt.Println(err)
					errflag = true
				}
			}
			if manager.Manager.GnbContext().GetStatus() == constant.CONTEXT_STATUS_GNB_RUNNING {
				if err := manager.Manager.GnbContext().Stop(); err != nil {
					fmt.Println(err)
					errflag = true
				}
			}
			time.Sleep(constant.WAIT_TIME)
			if errflag {
				fmt.Println(constant.OUTPUT_FAILURE)
			} else {
				fmt.Println(constant.OUTPUT_SUCCESS)
			}
			return
		}

		if err := cmd.ExecuteWithArgs(cmds); err != nil {
			fmt.Println(err)
			fmt.Println("Type 'help' to see available commands.")
			fmt.Println(constant.OUTPUT_FAILURE)
		}
	}
}
