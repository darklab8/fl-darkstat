package utils

import (
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/darklab8/go-utils/utils/utils_logus"
)

func SleepAwaitCtrlC() {
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}

func SleepForever() {
	utils_logus.Log.Debug("awaiting smth forever in SleepForever")
	mu := sync.Mutex{}
	mu.Lock()
	mu.Lock()
}
