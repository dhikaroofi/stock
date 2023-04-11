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
	runtime.GOMAXPROCS(runtime.NumCPU())
	color.Green(banner, "1.0", runtime.GOOS, runtime.GOARCH, runtime.Version(), runtime.NumCPU())

	appExitSignal := make(chan bool)
	interruptSignal := make(chan os.Signal, 1)

	conf := config.LoadConfigFile("config.yaml")
	cmd.Init(conf, appExitSignal)
	signal.Notify(interruptSignal, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	<-interruptSignal
	appExitSignal <- true
	<-appExitSignal
	logger.SysInfo("system is shut down gracefully")
}
