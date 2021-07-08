package interrupthandler

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

func HandleInterrupt(cb func()) {
	signals := make(chan os.Signal, 1)
	quit := make(chan bool, 1)

	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-signals
		log.Println("Program terminating. Got signal: ", sig)
		quit <- true
	}()
	<-quit
	cb()
	log.Println("Program terminated.")
}
