package log

import (
	c "golang.org/x/net/context"
	"sync"
)


const LoggerID = "logger_id"
var threadId = 0;
var mu sync.Mutex

func GetContext() c.Context {
	mu.Lock()
	ctx := c.Background()
	threadId++
	ctx2 := c.WithValue(ctx, LoggerID, threadId)
	mu.Unlock()
	return ctx2
}