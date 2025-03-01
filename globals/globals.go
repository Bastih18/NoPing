package globals

import "sync"

var Stop bool = false
var StopChan = make(chan struct{})
var Mutex = sync.Mutex{}

func StopProgram() {
	Mutex.Lock()
	Stop = true
	close(StopChan)
	Mutex.Unlock()
}
