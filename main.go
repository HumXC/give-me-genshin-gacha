package main

import (
	"flag"
	"fmt"
	"give-me-genshin-gacha/cli"
	"os"
)

func main() {
	var commandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	cli := cli.NewCli()
	commandLine.Usage = func() {
		fmt.Fprintf(commandLine.Output(), "Usage:\n")
		commandLine.PrintDefaults()
		cli.Usage()
	}
	commandLine.Parse(os.Args[1:])

	// 有子命令，则运行 cli
	if len(commandLine.Args()) > 0 {
		cli.Run(commandLine.Args()[0:])
		return
	}
	// 否则启动 gui
	fmt.Println("GUI")
}
