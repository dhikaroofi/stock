package main

import (
	"fmt"
	"github.com/dhikaroofi/stock.git/internal/cmd"
	"github.com/dhikaroofi/stock.git/internal/config"
	"github.com/dhikaroofi/stock.git/pkg/logger"
	"os"
	"os/signal"
	"runtime"
	"syscall"
)

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())
	fmt.Printf("version %s | OS %s %s %s CPU %v\n", "1.0", runtime.GOOS, runtime.GOARCH, runtime.Version(), runtime.NumCPU())
	conf := config.LoadConfigFile("config.yaml")

	appExitSignal := make(chan bool)
	interruptSignal := make(chan os.Signal, 1)

	cmd.Init(conf, appExitSignal)
	signal.Notify(interruptSignal, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	for range interruptSignal {
		appExitSignal <- true
		<-appExitSignal

		close(appExitSignal)

		logger.SysInfo("system is shut down gracefully")

		return // Now we can safely exit the app
	}

}
