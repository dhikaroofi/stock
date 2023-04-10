package main

import (
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/dhikaroofi/stock.git/internal/cmd"
	"github.com/dhikaroofi/stock.git/internal/config"
	"github.com/dhikaroofi/stock.git/pkg/logger"
	"github.com/fatih/color"
)

const banner = `
███████╗████████╗ ██████╗  ██████╗██╗  ██╗██████╗ ██╗████████╗
██╔════╝╚══██╔══╝██╔═══██╗██╔════╝██║ ██╔╝██╔══██╗██║╚══██╔══╝
███████╗   ██║   ██║   ██║██║     █████╔╝ ██████╔╝██║   ██║   
╚════██║   ██║   ██║   ██║██║     ██╔═██╗ ██╔══██╗██║   ██║   
███████║   ██║   ╚██████╔╝╚██████╗██║  ██╗██████╔╝██║   ██║   
╚══════╝   ╚═╝    ╚═════╝  ╚═════╝╚═╝  ╚═╝╚═════╝ ╚═╝   ╚═╝   
version %s | OS %s %s %s CPU %v
`

func main() {
	args := os.Args

	if len(args) > 1 {
		conf := config.LoadConfigFile("config.yaml")
		switch args[1] {
		case "part1":
			appExitSignal := make(chan bool)
			interruptSignal := make(chan os.Signal, 1)

			cmd.InitChallengePart1(conf, appExitSignal)
			signal.Notify(interruptSignal, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

			<-interruptSignal
			appExitSignal <- true
			<-appExitSignal
			logger.SysInfo("system is shut down gracefully")
			return
		case "part2":
			cmd.InitChallengePart2(conf)
		}
	} else {
		menuArguments()
	}
}

func menuArguments() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	color.Green(banner, "1.0", runtime.GOOS, runtime.GOARCH, runtime.Version(), runtime.NumCPU())
	color.Red("please use command below:")
	color.White("1.to run server for challenge part 1 you can type %s", color.HiRedString("'make run part1'"))
	color.White("2.to start challenge part 2 you can type %s", color.HiRedString("'make run challengePart2'"))
	color.Green("thank you")
}
