// Copyright (C) 2014 Space Monkey, Inc.

package spacelog

import (
	"sync"
	"sync/atomic"
)

// Logger is the basic type that allows for logging. A logger has an associated
// name, given to it during construction, either through a logger collection,
// GetLogger, GetLoggerNamed, or another Logger's Scope method. A logger also
// has an associated level and handler, typically configured through the logger
// collection to which it belongs.
type Logger struct {
	level      LogLevel
	name       string
	collection *LoggerCollection

	handler_mtx sync.RWMutex
	handler     Handler
}

// Scope returns a new Logger with the same level and handler, using the
// receiver Logger's name as a prefix.
func (l *Logger) Scope(name string) *Logger {
	return l.collection.getLogger(l.name+"."+name, l.getLevel(),
		l.getHandler())
}

func (l *Logger) setLevel(level LogLevel) {
	atomic.StoreInt32((*int32)(&l.level), int32(level))
}

func (l *Logger) getLevel() LogLevel {
	return LogLevel(atomic.LoadInt32((*int32)(&l.level)))
}

func (l *Logger) setHandler(handler Handler) {
	l.handler_mtx.Lock()
	defer l.handler_mtx.Unlock()
	l.handler = handler
}

func (l *Logger) getHandler() Handler {
	l.handler_mtx.RLock()
	defer l.handler_mtx.RUnlock()
	return l.handler
}
