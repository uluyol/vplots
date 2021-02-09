package builtin

import "runtime"

func init() {
	runtime.LockOSThread()
}

var mainq chan func()

func AppMain(mainFunc func()) {
	mainq = make(chan func())

	mainDone := make(chan int)
	go func() {
		mainFunc()
		mainDone <- 1
	}()

	for {
		select {
		case f := <-mainq:
			f()
		case <-mainDone:
			return
		}
	}
}

func onMain(f func()) {
	c := make(chan int)
	mainq <- func() {
		f()
		c <- 1
	}
	<-c
}
