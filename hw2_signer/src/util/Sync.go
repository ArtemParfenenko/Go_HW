package util

import "sync"

var mutexChannel = make(chan struct{}, 1)
var mutex = &sync.Mutex{}

type SyncExecutor func(action func())

func MutexExecutor(action func()) {
	mutex.Lock()
	defer mutex.Unlock()
	action()
}

func MutexChannelExecutor(action func()) {
	func(syncChannel chan struct{}) {
		syncChannel <- struct{}{}
		defer func() {
			<-syncChannel
		}()
		action()
	}(mutexChannel)
}
