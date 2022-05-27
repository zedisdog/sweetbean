package tools

import (
	"os"
	"os/signal"
	"syscall"
)

func Wait(closeFunc ...func()) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGKILL)
	<-c
	for _, cls := range closeFunc {
		cls()
	}
}
