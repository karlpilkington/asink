/*
 Copyright (C) 2013 Aaron Lindsay <aaron@aclindsay.com>
*/

package asink

import (
	"os"
	"os/signal"
	"sync/atomic"
)

var exitWaiterCount int32
var exitCalled chan int
var exitWaiterChan chan int

func init() {
	exitWaiterCount = 0
	exitWaiterChan = make(chan int)
	exitCalled = make(chan int)
}

func SetupCleanExitOnSignals() {
	go setupCleanExitOnSignals()
}
func setupCleanExitOnSignals() {
	//wait to properly close the socket when we're exiting
	exitCode := 0
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	defer signal.Stop(sig)

	select {
	case <-sig:
	case exitCode = <-exitCalled:
	}

	for c := atomic.AddInt32(&exitWaiterCount, -1); c >= 0; c = atomic.AddInt32(&exitWaiterCount, -1) {
		exitWaiterChan <- exitCode
	}
}

func Exit(exitCode int) {
	exitCalled <- exitCode
}

func WaitOnExit() int {
	atomic.AddInt32(&exitWaiterCount, 1)
	return <-exitWaiterChan
}

func exitChanWaiter(waiterChan chan int) {
	waiterChan <- <-exitWaiterChan
}
func WaitOnExitChan(waiterChan chan int) {
	atomic.AddInt32(&exitWaiterCount, 1)
	go exitChanWaiter(waiterChan)
}
