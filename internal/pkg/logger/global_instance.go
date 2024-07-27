package logger

import "sync"

var gl V1
var lMutex sync.Mutex

func SetGlobalInstance(l V1) {
	lMutex.Lock()
	defer lMutex.Unlock()
	gl = l
}
