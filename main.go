package main

import (
	"fmt"
	"strings"

	"github.com/chzyer/readline"
	"github.com/free-ran-ue/frush/constant"
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
`)
}

func main() {
	printFrush()

	rl, err := readline.New(constant.CMD_START)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := rl.Close(); err != nil {
			panic(err)
		}
	}()

	// frushManager := manager.NewManager()

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
		default:
			fmt.Println(fmt.Sprintf(constant.SYSTEM_HINT_UNKNOWN_CMD, cmds[0]))
			usage()
		}
	}
}
