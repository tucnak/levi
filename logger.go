package levi

import (
	"fmt"
	"sync"
)

// Logger performs the logging once the request is finished.
type Logger interface {
	Report(*Lv) error
}

type StdLogger struct {
	sync.Mutex
}

func (std *StdLogger) Report(lv *Lv) error {
	std.Lock()
	for _, log := range lv.Logs {
		fmt.Println(log)
	}
	std.Unlock()
	return nil
}
