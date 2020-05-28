package levi

import (
	"fmt"
	"time"
)

// Logotype allows to differentiate between log messages.
type Logotype int

const (
	PRINT   Logotype = iota
	DEBUG            // …
	WARNING          // ?
	ERROR            // !
	PANIC            // *
)

func (k Logotype) String() string {
	switch k {
	case PRINT:
		return ""
	case DEBUG:
		return "…"
	case WARNING:
		return "?"
	case ERROR:
		return "!"
	case PANIC:
		return "*"
	default:
		panic("levi: unknown Log type")
	}
}

// Log is a single log entry.
type Log struct {
	How  Logotype
	When time.Duration
	What []byte
}

func (l Log) String() string {
	if l.How == PRINT {
		return string(l.What)
	}
	return fmt.Sprintf("<%s> [%s] %s", l.How, l.When, l.What)
}

func (lv *Lv) log(kind Logotype, stuff ...interface{}) {
	log := Log{kind, time.Now().Sub(lv.started), []byte(fmt.Sprintln(stuff...))}
	lv.Logs = append(lv.Logs, log)
}

func (lv *Lv) logf(kind Logotype, format string, stuff ...interface{}) {
	log := Log{kind, time.Now().Sub(lv.started), []byte(fmt.Sprintf(format, stuff...))}
	lv.Logs = append(lv.Logs, log)
}

func (lv *Lv) Debug(info ...interface{})              { lv.log(DEBUG, info...) }
func (lv *Lv) Debugf(fmt string, info ...interface{}) { lv.logf(DEBUG, fmt, info...) }
func (lv *Lv) Warn(msg ...interface{})                { lv.log(WARNING, msg...) }
func (lv *Lv) Warnf(fmt string, data ...interface{})  { lv.logf(WARNING, fmt, data...) }
func (lv *Lv) Error(err error)                        { lv.log(ERROR, err) }
func (lv *Lv) Errorf(fmt string, data ...interface{}) { lv.logf(ERROR, fmt, data...) }
func (lv *Lv) Panic(msg ...interface{})               { lv.log(PANIC, msg...) }
func (lv *Lv) Panicf(fmt string, data ...interface{}) { lv.logf(PANIC, fmt, data...) }
